package auth

import (
	"auth/core/ports"
	"auth/internal"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
)

// Controller - health controller
type Controller struct {
	log      *otelzap.Logger
	config   *internal.AppConfig
	validate ports.AppValidator
	service  ports.AuthService
}

// NewAuthController - new health controller
func NewAuthController(log *otelzap.Logger, api *gin.RouterGroup, config *internal.AppConfig, authService ports.AuthService, validator ports.AppValidator) *Controller {
	controller := &Controller{
		log:      log,
		config:   config,
		validate: validator,
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
