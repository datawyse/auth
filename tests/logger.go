package tests

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func NewTestLog(t *testing.T) *zap.Logger {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	return observedLogger
}
