package ports

import (
	"context"
	"github.com/google/uuid"
)

type UUIDService interface {
	FromString(ctx context.Context, id string) (uuid.UUID, error)
	NewUUID(ctx context.Context) uuid.UUID
	IsValidUUID(ctx context.Context, id string) (bool, error)
}
