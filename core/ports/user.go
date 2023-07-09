package ports

import (
	"context"

	"auth/core/domain"

	"github.com/google/uuid"
)

type UserService interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, input *domain.User) (string, error)

	// User returns the user with the given id.
	User(ctx context.Context, id string) (*domain.UserInfo, error)

	// UpdateUser updates a user
	UpdateUser(input *domain.User) (*domain.User, error)

	// UserByEmail returns the user with the given email.
	UserByEmail(ctx context.Context, email string) (*domain.UserInfo, error)

	UserByUsername(ctx context.Context, username string) (*domain.UserInfo, error)

	// Users returns all users.
	Users(ctx context.Context) ([]*domain.UserInfo, error)
}

type UserRepository interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, input *domain.User) (uuid.UUID, error)

	// User returns the user with the given id.
	User(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, input *domain.User) (*domain.User, error)

	// Users returns all users.
	Users(ctx context.Context) ([]*domain.User, error)
}
