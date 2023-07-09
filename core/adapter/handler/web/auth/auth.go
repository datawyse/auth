package auth

import (
	"auth/core/ports"
	"auth/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Controller - health controller
type Controller struct {
	log      *zap.Logger
	config   *internal.AppConfig
	validate *validator.Validate
	service  ports.AuthService
}

// NewAuthController - new health controller
func NewAuthController(log *zap.Logger, api *gin.RouterGroup, config *internal.AppConfig, authService ports.AuthService) *Controller {
	controller := &Controller{
		log:      log,
		config:   config,
		validate: validator.New(),
		service:  authService,
	}

	// Declare routing for specific routes.
	authRoute := api.Group("/user")

	authRoute.GET("/me", controller.Me)
	authRoute.POST("/login", controller.Login)
	authRoute.POST("/signup", controller.Signup)
	authRoute.POST("/refresh-token", controller.RefreshToken)

	return controller
}

var HandlerModule = fx.Options(fx.Invoke(NewAuthController))
