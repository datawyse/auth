package app

import (
	"github.com/gin-gonic/gin"
)

// -------------------------------------------------------------------
// Serve events data
// -------------------------------------------------------------------

// BootstrapEvent is the event data for the bootstrap event of the application
type BootstrapEvent struct {
	App App
}

// ServeEvent is the event data for the serve event of the application
type ServeEvent struct {
	App    App
	Router *gin.Engine
}

// BeforeMongoDbConnectionEvent is the event data for the before mongodb connection event of the application
type BeforeMongoDbConnectionEvent struct {
	App         App
	DatabaseUrl string
	Database    string
}

// AfterMongoDbConnectionEvent is the event data for the after mongodb connection event of the application
type AfterMongoDbConnectionEvent struct {
	App App
}

// ApiErrorEvent is the event data for the api error event of the application
type ApiErrorEvent struct {
	HttpContext gin.Context
	Error       error
}
