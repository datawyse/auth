package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"auth/core"
	"auth/internal"
	"auth/internal/app"
	"auth/internal/app/grpc"
	"auth/internal/app/web"
	"auth/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func StartGRPCServer(lf fx.Lifecycle, log *zap.Logger, grpcServer *grpc.Server, config *internal.AppConfig) {
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting gRPC server")
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
			return grpcServer.Stop()
		},
	})
}

func StartHTTPServer(lf fx.Lifecycle, log *zap.Logger, router *gin.Engine, authApp app.App, config *internal.AppConfig) {
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting http server")
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
				fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: log}
				}),

				//  provide configFile path as a global dependency for the application
				fx.Provide(
					func() string {
						return configPath
					},
				),

				fx.Provide(
					// inject the logger as a global dependency for the application
					// check for authServiceApp mode and inject the appropriate logger
					zap.NewDevelopment,

					// inject empty context on application startup
					fx.Annotate(ctx, fx.As(new(context.Context))),

					// load the application config
					internal.AppConfigModule,
				),

				// inject the application database
				db.Module,

				// inject the gin router
				web.WebAppModule,
				grpc.ServerModule,

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
