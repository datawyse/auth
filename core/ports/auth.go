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
	Login(username string, password string) (*domain.AuthToken, error)

	// RefreshToken refreshes a token.
	RefreshToken(refreshToken string) (*domain.AuthToken, error)

	// User returns the current user.
	User(userId string) (*domain.UserInfo, error)

	UpdateUser(ctx context.Context, input *http.UpdateUserInput) (*domain.User, error)
}
