package permissions

import (
	"context"

	"auth/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Controller struct {
	log      *zap.Logger
	config   *internal.AppConfig
	ctx      context.Context
	validate *validator.Validate
}

func NewPermissionsController(api *gin.RouterGroup, log *zap.Logger, config *internal.AppConfig, ctx context.Context) *Controller {
	controller := &Controller{
		log:      log,
		config:   config,
		ctx:      ctx,
		validate: validator.New(),
	}

	// Declare routing for specific routes.
	permissionsRoute := api.Group("/permissions")

	permissionsRoute.GET("", controller.GetPermissions)

	return controller
}

var HandlerModule = fx.Options(
	fx.Invoke(NewPermissionsController),
)
