package health

import (
	"auth/core/ports"
	"auth/internal"
	"auth/internal/db/mongodb"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Controller - health controller
type Controller struct {
	log     *otelzap.Logger
	config  *internal.AppConfig
	db      *mongodb.MongoDb
	service ports.HealthService
}

// NewHealthCtrl - new health controller
func NewHealthCtrl(log *otelzap.Logger, config *internal.AppConfig, api *gin.RouterGroup, db *mongodb.MongoDb, healthService ports.HealthService) *Controller {
	controller := &Controller{
		log:     log,
		config:  config,
		db:      db,
		service: healthService,
	}

	// Declare routing for specific routes.
	route := api.Group("/health")
	route.GET("", controller.ReadHealth)

	return controller
}

var HandlerModule = fx.Options(fx.Invoke(NewHealthCtrl))
