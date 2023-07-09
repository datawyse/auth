package profile

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
	log *zap.Logger
	*internal.AppConfig
	ctx         context.Context
	validate    *validator.Validate
	service     ports.AuthService
	userService ports.UserService
}

func NewProfileController(log *zap.Logger, api *gin.RouterGroup, config *internal.AppConfig, ctx context.Context, authService ports.AuthService, userService ports.UserService) (*Controller, error) {
	controller := &Controller{
		log:         log,
		AppConfig:   config,
		ctx:         ctx,
		validate:    validator.New(),
		service:     authService,
		userService: userService,
	}

	// Declare routing for specific routes.
	profileRoute := api.Group("/profile")

	profileRoute.GET("/:id", controller.GetProfile)
	profileRoute.PUT("/:id", controller.UpdateProfile)
	profileRoute.GET("/search", controller.SearchProfile)
	profileRoute.GET("/import", controller.ImportProfiles)

	return controller, nil
}

var HandlerModule = fx.Options(fx.Invoke(NewProfileController))
