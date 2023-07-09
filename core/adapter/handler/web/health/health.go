package health

import (
	"auth/core/ports"
	"auth/core/service/app-validator"
	"auth/internal"
	"auth/internal/db/mongodb"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Controller - health controller
type Controller struct {
	log          *zap.Logger
	config       *internal.AppConfig
	db           *mongodb.MongoDb
	appValidator *app_validator.AppValidator
	service      ports.HealthService
}

// NewHealthCtrl - new health controller
func NewHealthCtrl(log *zap.Logger, config *internal.AppConfig, api *gin.RouterGroup, db *mongodb.MongoDb, healthService ports.HealthService, validator *app_validator.AppValidator) *Controller {
	controller := &Controller{
		log:          log,
		config:       config,
		db:           db,
		appValidator: validator,
		service:      healthService,
	}

	// Declare routing for specific routes.
	route := api.Group("/health")
	route.GET("", controller.ReadHealth)

	return controller
}

var HandlerModule = fx.Options(fx.Invoke(NewHealthCtrl))
