package ports

import (
	"context"

	"auth/core/domain/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
)

type AuthServerService interface {
	// AccessToken get access token
	AccessToken(ctx context.Context) (string, error)

	Login(ctx context.Context, email, password string) (*gocloak.JWT, error)

	NewClient(ctx context.Context) *gocloak.GoCloak

	CreateUser(ctx context.Context, input *http.SignupInput) (*gocloak.User, error)

	GetUserById(ctx context.Context, id string) (*gocloak.User, error)

	// RetrospectToken get permissions with roles
	RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error)

	// VerifyToken verify token
	VerifyToken(ctx context.Context, accessToken string) (*jwt.Token, error)
}
