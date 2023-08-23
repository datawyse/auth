package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx/fxevent"
	"log"
	"net/http"
	"sync"
	"time"

	"auth/core"
	"auth/internal"
	"auth/internal/app"
	"auth/internal/app/grpc"
	"auth/internal/app/grpc/server"
	"auth/internal/app/web"
	"auth/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

var resource *sdkresource.Resource
var initResourcesOnce sync.Once

func initTracerExporter(config *internal.AppConfig) *otlptrace.Exporter {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if config.AppMode != "production" {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.CollectorURL),
		),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	return exporter
}

func initMetricExporter(config *internal.AppConfig) sdkmetric.Exporter {
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.CollectorURL),
	}

	exporter, err := otlpmetricgrpc.New(context.Background(), options...)
	if err != nil {
		log.Fatal(err.Error())
	}

	return exporter
}

func initResource(config *internal.AppConfig) *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				attribute.String("library.language", "go"),
				attribute.String("app.mode", config.AppMode),
				attribute.String("app.version", config.AppVersion),
				attribute.String("service.name", config.ServiceName),
			),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initMeterProvider(config *internal.AppConfig) *sdkmetric.MeterProvider {
	otelMetricExporter := initMetricExporter(config)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(initResource(config)),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otelMetricExporter)),
	)
	otel.SetMeterProvider(mp)
	return mp
}

func initTracerProvider(config *internal.AppConfig) *trace.TracerProvider {
	otelResource := initResource(config)
	otelTracerExporter := initTracerExporter(config)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(otelTracerExporter),
		trace.WithResource(otelResource),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func StartGRPCServer(lf fx.Lifecycle, log *zap.Logger, grpcServer *server.Server, config *internal.AppConfig) {
	metric := initMeterProvider(config)
	tracer := initTracerProvider(config)

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
				log.Fatal(err.Error())
			}

			log.Info("Starting gRPC server")
			go func() {
				log.Info("Starting gRPC server", zap.String("addr", fmt.Sprint(":", config.GRPCPort)))
				if err := grpcServer.Start(); err != nil {
					log.Fatal(err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping gRPC server")

			// stop tracer
			if err := tracer.Shutdown(ctx); err != nil {
				log.Error("Failed to stop tracer", zap.Error(err))
				panic(err)
			}

			// stop metric
			if err := metric.Shutdown(ctx); err != nil {
				log.Error("Failed to stop metric", zap.Error(err))
				panic(err)
			}

			return grpcServer.Stop()
		},
	})
}

func StartHTTPServer(lf fx.Lifecycle, log *zap.Logger, router *gin.Engine, authApp app.App, config *internal.AppConfig) {
	metric := initMeterProvider(config)
	tracer := initTracerProvider(config)

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
				log.Fatal(err.Error())
			}

			log.Info("Starting http server")

			if err := authApp.Bootstrap(log, config); err != nil {
				log.Fatal(err.Error())
			}

			if err := authApp.RefreshSettings(log, config); err != nil {
				log.Fatal(err.Error())
			}

			go func() {
				log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
				// service connections
				if err := srv.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						log.Error("Server closed unexpectedly", zap.Error(err))
						log.Fatal(err.Error())
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server", zap.String("addr", srv.Addr))

			// shutdown the server with a timeout
			defer func(authApp app.App, logger *zap.Logger, config *internal.AppConfig) {
				// stop tracer
				if err := tracer.Shutdown(ctx); err != nil {
					logger.Error("Failed to stop tracer", zap.Error(err))
					panic(err)
				}

				// stop metric
				if err := metric.Shutdown(ctx); err != nil {
					logger.Error("Failed to stop metric", zap.Error(err))
					panic(err)
				}

				err := authApp.ResetBootstrapState(logger, config)
				if err != nil {
					logger.Error("Failed to reset bootstrap state", zap.Error(err))
					panic(err)
				}
			}(authApp, log, config)
			return srv.Shutdown(ctx)
		},
	})
}

func Serve() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web server (default to 127.0.0.1:3000)",
		Run: func(command *cobra.Command, args []string) {
			// get config file path from command flags
			configPath, err := command.Flags().GetString("config")
			if err != nil {
				panic(err)
			}

			emptyContext := context.Background()
			// create an empty context as base and reuse it later throughout the application
			ctx := func() context.Context {
				return emptyContext
			}

			authServiceApp := fx.New(
				fx.WithLogger(func(log *otelzap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: log.Logger}
				}),

				//  provide configFile path as a global dependency for the application
				fx.Provide(
					func() string {
						return configPath
					},
				),

				fx.Provide(
					// inject empty context on application startup
					fx.Annotate(ctx, fx.As(new(context.Context))),
				),

				fx.Provide(
					internal.AppConfigModule,
				),

				fx.Provide(func(config *internal.AppConfig) (*zap.Logger, error) {
					// inject the logger as a global dependency for the application
					if config.AppMode == "production" {
						return zap.NewProduction()
					}
					return zap.NewDevelopment()
				}),

				fx.Provide(
					// inject the logger as a global dependency for the application
					otelzap.New,
				),

				// inject the application database
				db.Module,

				// inject the application services
				web.Module,
				grpc.Module,

				// inject the application controllers
				core.Module,

				// inject the application controllers and register them with the router as routes and handlers for the application services
				fx.Invoke(
					StartHTTPServer,
					StartGRPCServer,
					func(logger *zap.Logger) {
						logger.Debug("Logger module invoked")
					},
					func(log *zap.Logger) {
						log.Info("Starting application with config file: " + configPath)
					},
				),
			)

			startCtx, cancel := context.WithTimeout(ctx(), 5*time.Second)
			defer cancel()
			if err := authServiceApp.Start(startCtx); err != nil {
				log.Fatal(fmt.Errorf("app.Start: %w", err))
			}

			<-authServiceApp.Wait()
		},
	}

	command.PersistentFlags().StringP("config", "c", "config.yml", "config file path")
	return command
}
