package tests

import (
	"testing"

	"auth/internal"

	"go.uber.org/zap"
)

func NewTestConfig(t *testing.T) *internal.AppConfig {
	logger := NewTestLog(t)
	logger.Log(zap.DebugLevel, "creating config")
	config, err := internal.NewConfig(logger, "")
	if err != nil {
		panic("error loading test config")
	}

	return config
}
