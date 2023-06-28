package subscriptions

import (
	"context"

	"auth/core/ports"
	"auth/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Controller struct {
	log         *zap.Logger
	config      *internal.AppConfig
	ctx         context.Context
	validate    *validator.Validate
	service     ports.SubscriptionService
	userService ports.UserService
}

func NewSubscriptionsController(ctx context.Context, api *gin.RouterGroup, log *zap.Logger, config *internal.AppConfig, validate *validator.Validate, service ports.SubscriptionService, userService ports.UserService) (*Controller, error) {
	controller := &Controller{
		log:         log,
		config:      config,
		ctx:         ctx,
		validate:    validate,
		service:     service,
		userService: userService,
	}

	// Declare routing for specific routes.
	subscriptionsRoute := api.Group("/subscriptions")
	subscriptionsRoute.GET("", controller.GetSubscriptions)

	return controller, nil
}

var HandlerModule = fx.Options(
	fx.Invoke(NewSubscriptionsController),
)
