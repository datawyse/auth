package app

import (
	"context"
	"fmt"

	"auth/internal"

	"go.uber.org/zap"
)

type App interface {
	// Bootstrap takes care for initializing the application
	// (open db connections, load settings, etc.)
	Bootstrap(*zap.Logger, *internal.AppConfig) error

	// ResetBootstrapState takes care for releasing initialized app resources
	// (ex. closing db connections).
	ResetBootstrapState(*zap.Logger, *internal.AppConfig) error

	// Settings returns the loaded app settings.
	Settings(*zap.Logger, *internal.AppConfig) *Settings

	// RefreshSettings re-initializes and reloads the stored application settings.
	RefreshSettings(*zap.Logger, *internal.AppConfig) error

	// ---------------------------------------------------------------
	// App event hooks
	// ---------------------------------------------------------------

	// OnBeforeBootstrap hook is triggered before initializing the base
	// application resources (ex. before db open and initial settings load).
	OnBeforeBootstrap(*zap.Logger, *internal.AppConfig) *Hook[*BootstrapEvent]

	// OnAfterBootstrap hook is triggered after initializing the base
	// application resources (ex. after db open and initial settings load).
	OnAfterBootstrap(*zap.Logger, *internal.AppConfig) *Hook[*BootstrapEvent]

	// OnBeforeServe hook is triggered before serving the internal router (echo),
	// allowing you to adjust its options and attach new routes.
	OnBeforeServe(*zap.Logger, *internal.AppConfig) *Hook[*ServeEvent]

	// OnBeforeApiError hook is triggered right before sending an error API
	// response to the client, allowing you to further modify the error data
	// or to return a completely different API response (using [hook.StopPropagation]).
	OnBeforeApiError(*zap.Logger, *internal.AppConfig) *Hook[*ApiErrorEvent]

	// OnAfterApiError hook is triggered right after sending an error API
	// response to the client.
	// It could be used to log the final API error in external services.
	OnAfterApiError(*zap.Logger, *internal.AppConfig) *Hook[*ApiErrorEvent]

	//	OnBeforeMongoDbConnect events
	OnBeforeMongoDbConnect(*zap.Logger, *internal.AppConfig) *Hook[*BeforeMongoDbConnectionEvent]

	//	OnAfterMongoDbConnect events
	OnAfterMongoDbConnect(*zap.Logger, *internal.AppConfig) *Hook[*AfterMongoDbConnectionEvent]
}

// ProjectService is the main application interface. It is used to bootstrap the application. It also provides hooks for the application events.
type ProjectService struct {
	settings *Settings

	// app event hooks
	onBeforeBootstrap      *Hook[*BootstrapEvent]
	onAfterBootstrap       *Hook[*BootstrapEvent]
	onBeforeServe          *Hook[*ServeEvent]
	onBeforeApiError       *Hook[*ApiErrorEvent]
	onAfterApiError        *Hook[*ApiErrorEvent]
	onBeforeMongoDbConnect *Hook[*BeforeMongoDbConnectionEvent]
	onAfterMongoDbConnect  *Hook[*AfterMongoDbConnectionEvent]
}

func (app *ProjectService) OnBeforeServe(log *zap.Logger, config *internal.AppConfig) *Hook[*ServeEvent] {
	log.Debug(fmt.Sprintf("@@ OnBeforeServe HOOK"))
	return app.onBeforeServe
}

func (app *ProjectService) OnBeforeBootstrap(log *zap.Logger, config *internal.AppConfig) *Hook[*BootstrapEvent] {
	log.Debug(fmt.Sprintf("@@ OnBeforeBootstrap HOOK"))
	return app.onBeforeBootstrap
}

func (app *ProjectService) OnAfterBootstrap(log *zap.Logger, config *internal.AppConfig) *Hook[*BootstrapEvent] {
	log.Debug(fmt.Sprintf("@@ OnAfterBootstrap HOOK"))
	return app.onBeforeBootstrap
}

func (app *ProjectService) OnBeforeApiError(log *zap.Logger, config *internal.AppConfig) *Hook[*ApiErrorEvent] {
	log.Debug(fmt.Sprintf("@@ OnBeforeApiError HOOK"))
	return app.onBeforeApiError
}

func (app *ProjectService) OnAfterApiError(log *zap.Logger, config *internal.AppConfig) *Hook[*ApiErrorEvent] {
	log.Debug(fmt.Sprintf("@@ OnAfterApiError HOOK"))
	return app.onAfterApiError
}

// OnBeforeMongoDbConnect returns the loaded app settings.
func (app *ProjectService) OnBeforeMongoDbConnect(log *zap.Logger, config *internal.AppConfig) *Hook[*BeforeMongoDbConnectionEvent] {
	log.Debug(fmt.Sprintf("@@ OnBeforeMongoDbConnect HOOK"))
	return app.onBeforeMongoDbConnect
}

// OnAfterMongoDbConnect returns the loaded app settings.
func (app *ProjectService) OnAfterMongoDbConnect(log *zap.Logger, config *internal.AppConfig) *Hook[*AfterMongoDbConnectionEvent] {
	log.Debug(fmt.Sprintf("@@ OnAfterMongoDbConnect HOOK"))
	return app.onAfterMongoDbConnect
}

func (app *ProjectService) RefreshSettings(log *zap.Logger, config *internal.AppConfig) error {
	log.Debug(fmt.Sprintf("refreshing app settings"))

	if app.settings == nil {
		app.settings = NewSettings(config)
	}

	return nil
}

// Settings returns the loaded app settings.
func (app *ProjectService) Settings(log *zap.Logger, config *internal.AppConfig) *Settings {
	log.Debug(fmt.Sprintf("getting app settings"))
	return app.settings
}

func (app *ProjectService) Bootstrap(log *zap.Logger, config *internal.AppConfig) error {
	log.Debug(fmt.Sprintf("bootstrapping app"))

	event := &BootstrapEvent{App: app}
	if err := app.OnBeforeBootstrap(log, config).Trigger(event); err != nil {
		return err
	}

	// clear resources of previous core state (if any)
	if err := app.ResetBootstrapState(log, config); err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("running bootstrap"))

	// we don't check for an error because the db migrations may have not been executed yet
	if err := app.RefreshSettings(log, config); err != nil {
		return fmt.Errorf("error while refreshing app settings: %w", err)
	}

	if err := app.OnAfterBootstrap(log, config).Trigger(event); err != nil {
		log.Debug(fmt.Sprintf(err.Error()))
	}

	return nil
}

// ResetBootstrapState takes care for releasing initialized app resources
// (ex. closing db connections).
func (app *ProjectService) ResetBootstrapState(log *zap.Logger, config *internal.AppConfig) error {
	log.Debug(fmt.Sprintf("resetting app"))

	app.settings = nil

	return nil
}

func NewDataFlow(ctx context.Context, config *internal.AppConfig) App {
	app := &ProjectService{
		settings: NewSettings(config),

		onBeforeServe:          &Hook[*ServeEvent]{},
		onBeforeBootstrap:      &Hook[*BootstrapEvent]{},
		onAfterBootstrap:       &Hook[*BootstrapEvent]{},
		onBeforeApiError:       &Hook[*ApiErrorEvent]{},
		onAfterApiError:        &Hook[*ApiErrorEvent]{},
		onBeforeMongoDbConnect: &Hook[*BeforeMongoDbConnectionEvent]{},
		onAfterMongoDbConnect:  &Hook[*AfterMongoDbConnectionEvent]{},
	}

	return app
}

var AppModule = NewDataFlow
