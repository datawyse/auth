package middlewares

import (
	"context"

	"auth/core/service/app-validator"
	"auth/internal"

	"github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitMiddleware(ctx context.Context, app *gin.Engine, log *zap.Logger, appConfig *internal.AppConfig, appValidator *app_validator.AppValidator) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	app.Use(cors.New(config))

	app.Use(helmet.Default())

	app.Use(gin.Recovery())

	app.Use(RevisionMiddleware(log))

	app.Use(RequestIdMiddleware(log))

	app.Use(Authorization(log, appConfig))

	app.Use(ErrorHandler(log, appValidator))
}
