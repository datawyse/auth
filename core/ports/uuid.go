package ports

import "github.com/google/uuid"

type UUIDService interface {
	FromString(id string) (uuid.UUID, error)
	NewUUID() uuid.UUID
	IsValidUUID(id string) (bool, error)
}
