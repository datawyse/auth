package app

import (
	"sync"

	"auth/internal"
)

// Settings defines common app configuration options.
type Settings struct {
	mux sync.RWMutex

	Config *internal.AppConfig `json:"config"`
}

// NewSettings creates and returns a new default Settings instance.
func NewSettings(config *internal.AppConfig) *Settings {
	return &Settings{
		Config: config,
	}
}
