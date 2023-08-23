package middlewares

import (
	"auth/core/ports"
	"auth/internal"
	"context"
	"github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitMiddleware(ctx context.Context, app *gin.Engine, log *otelzap.Logger, config *internal.AppConfig, appValidator ports.AppValidator) {
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowCredentials = true
	ginConfig.AllowOrigins = []string{"http://localhost:5555", "http://localhost:9080", "http://localhost:3000"}
	app.Use(cors.New(ginConfig))

	app.Use(helmet.Default())
	app.Use(RevisionMiddleware(log))
	app.Use(RequestIdMiddleware(log))
	app.Use(Authorization(log, config))

	app.Use(otelgin.Middleware(config.ServiceName))
	app.Use(GetIDMiddleware(log))

	app.Use(ErrorHandler(log, appValidator))
	app.Use(gin.Recovery())

	//app.Use(ginzap.RecoveryWithZap(log.Logger, true))
	//app.Use(ginzap.Ginzap(log.Logger, time.RFC3339, true))
}
