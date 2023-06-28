package system

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	return uuid.New()
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// ToUUID to uuid.UUID
func ToUUID(id string) (uuid.UUID, error) {
	if !IsValidUUID(id) {
		return uuid.Nil, ErrInvalidUUID
	}

	return uuid.Parse(id)
}
