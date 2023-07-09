package http

import "time"

type LoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

type SignupInput struct {
	FirstName  string              `json:"firstName"  validate:"omitempty,min=3,max=32"`
	LastName   string              `json:"lastName"   validate:"omitempty,min=3,max=32"`
	Email      string              `json:"email"      validate:"required,email"          binding:"required"`
	Username   string              `json:"username"   validate:"required,min=3,max=32"   binding:"required"`
	Password   string              `json:"password"   validate:"required,min=6"          binding:"required"`
	Attributes map[string][]string `json:"attributes"`
}

type UpdateUserInput struct {
	Id        string `json:"id" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	// Language      string              `json:"language"`
	// AccountType   string              `mapstructure:"accountType" json:"accountType,omitempty"`
	LastSignInAt time.Time `mapstructure:"lastSignInAt" json:"lastSignInAt,omitempty"`
	// Organizations []string            `json:"organizations"`
	// Roles         []string            `json:"roles"`
	Attributes map[string][]string `json:"attributes"`
}
