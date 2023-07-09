package ports

import (
	"context"

	"auth/core/domain"
)

// HealthService is the interface that provides health methods.
type HealthService interface {
	// ReadHealth returns the health.
	ReadHealth(ctx context.Context) (*domain.Health, error)
}
