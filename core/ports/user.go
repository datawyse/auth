package ports

import (
	"context"

	"auth/core/domain"

	"github.com/google/uuid"
)

type UserService interface {
	// CreateUser creates a new user.
	CreateUser(input *domain.User) (string, error)

	// User returns the user with the given id.
	User(id string) (*domain.UserInfo, error)

	// UpdateUser updates a user
	UpdateUser(input *domain.User) (*domain.User, error)
}

type UserRepository interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, input *domain.User) (uuid.UUID, error)

	// User returns the user with the given id.
	User(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, input *domain.User) (*domain.User, error)
}
