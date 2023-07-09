package ports

import (
	"context"

	"auth/core/domain"
	"auth/core/domain/http"
)

type AuthService interface {
	// Signup creates a new user.
	Signup(ctx context.Context, input *http.SignupInput) (string, error)

	// Login authenticates a user.
	Login(ctx context.Context, email, password string) (*domain.AuthToken, error)

	// RefreshToken refreshes a token.
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthToken, error)

	// User returns the current user.
	User(ctx context.Context, userId string) (*domain.UserInfo, error)

	UpdateUser(ctx context.Context, input *http.UpdateUserInput) (*domain.UserInfo, error)
}
