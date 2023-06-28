package tests

import (
	"sync"
	"testing"

	"auth/internal/app"
)

type TestApp struct {
	*app.ProjectService

	mux sync.Mutex

	// EventCalls defines a map to inspect which app events
	// (and how many times) were triggered.
	// The following events are not counted because they execute always:
	// - OnBeforeBootstrap
	// - OnAfterBootstrap
	// - OnBeforeServe
	EventCalls map[string]int
}

// ResetEventCalls resets the EventCalls counter.
func (t *TestApp) ResetEventCalls() {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.EventCalls = make(map[string]int)
}

func (t *TestApp) Cleanup() {
	t.ResetEventCalls()
}

func (t *TestApp) registerEventCall(name string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.EventCalls == nil {
		t.EventCalls = make(map[string]int)
	}

	t.EventCalls[name]++

	return nil
}

// NewTestApp creates and initializes a test application instance.
// It is the caller's responsibility to call `app.Cleanup()`
// when the app is no longer needed.
func NewTestApp(t *testing.T) (*TestApp, error) {
	log := NewTestLog(t)
	projectService := &app.ProjectService{}
	appConfig := NewTestConfig(t)

	// load data dir and db connections
	if err := projectService.Bootstrap(log, appConfig); err != nil {
		return nil, err
	}

	testApp := &TestApp{ProjectService: projectService}

	testApp.OnBeforeServe(log, appConfig).Add(func(e *app.ServeEvent) error {
		return testApp.registerEventCall("OnBeforeServe")
	})
	testApp.OnBeforeBootstrap(log, appConfig).Add(func(e *app.BootstrapEvent) error {
		return testApp.registerEventCall("OnBeforeBootstrap")
	})
	testApp.OnAfterBootstrap(log, appConfig).Add(func(e *app.BootstrapEvent) error {
		return testApp.registerEventCall("OnAfterBootstrap")
	})
	testApp.OnBeforeApiError(log, appConfig).Add(func(e *app.ApiErrorEvent) error {
		return testApp.registerEventCall("OnBeforeApiError")
	})
	testApp.OnAfterApiError(log, appConfig).Add(func(e *app.ApiErrorEvent) error {
		return testApp.registerEventCall("OnAfterApiError")
	})
	testApp.OnBeforeMongoDbConnect(log, appConfig).Add(func(e *app.BeforeMongoDbConnectionEvent) error {
		return testApp.registerEventCall("OnBeforeMongoDbConnect")
	})
	testApp.OnAfterMongoDbConnect(log, appConfig).Add(func(e *app.AfterMongoDbConnectionEvent) error {
		return testApp.registerEventCall("OnAfterMongoDbConnect")
	})

	return testApp, nil
}
