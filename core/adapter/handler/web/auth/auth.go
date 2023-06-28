package auth

import (
	"context"

	"auth/core/ports"
	"auth/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Controller - health controller
type Controller struct {
	log *zap.Logger
	*internal.AppConfig
	ctx      context.Context
	validate *validator.Validate
	service  ports.AuthService
}

// NewAuthController - new health controller
func NewAuthController(log *zap.Logger, api *gin.RouterGroup, ctx context.Context, config *internal.AppConfig, authService ports.AuthService) *Controller {
	controller := &Controller{
		log:       log,
		AppConfig: config,
		ctx:       ctx,
		validate:  validator.New(),
		service:   authService,
	}

	// Declare routing for specific routes.
	authRoute := api.Group("/user")

	authRoute.GET("/me", controller.Me)
	authRoute.POST("/login", controller.Login)
	authRoute.POST("/signup", controller.Signup)
	authRoute.POST("/refresh-token", controller.RefreshToken)

	authRoute.PUT("/profile/:id", controller.UpdateProfile)

	return controller
}

var HandlerModule = fx.Options(fx.Invoke(NewAuthController))
