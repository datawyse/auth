package server

import (
	"auth/core/domain/system"
	"context"

	"auth/internal"
	"auth/internal/app"
	"auth/internal/app/web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinRouter(engine *gin.Engine, log *zap.Logger) *gin.RouterGroup {
	log.Info("setting routes to v1")

	v1 := engine.Group("api/v1")
	return v1
}

// NewWebEngine returns a new gin router
func NewWebEngine(ctx context.Context, config *internal.AppConfig, log *zap.Logger, authApp app.App) (*gin.Engine, error) {
	// set app release mode
	if config.AppMode != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.RedirectFixedPath = true
	engine.RedirectTrailingSlash = false
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	// update engine context
	engine.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", gin.MIMEJSON)
		ctx.Next()
	})

	// trigger the custom BeforeServe hook for the created web router
	// allowing users to further adjust its options or register new routes
	serveEvent := &app.ServeEvent{
		App:    authApp,
		Router: engine,
	}
	if err := authApp.OnBeforeServe(log, config).Trigger(serveEvent); err != nil {
		return nil, err
	}

	// 404 router
	engine.NoRoute(func(ctx *gin.Context) {
		_ = ctx.Error(system.ErrNotFound)
		return
	})

	return engine, nil
}

var EngineModule = fx.Provide(NewWebEngine)
var ServiceAppModule = fx.Provide(app.AppModule)
var GinRouterModule = fx.Provide(NewGinRouter)
var MiddlewareModule = fx.Options(fx.Invoke(middlewares.InitMiddleware))
