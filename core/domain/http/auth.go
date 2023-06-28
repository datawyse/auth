package http

import "time"

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

type SignupInput struct {
	FirstName string `json:"firstName" binding:"required" validate:"required,min=3,max=32"`
	LastName  string `json:"lastName" binding:"required" validate:"required,min=3,max=32"`
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Username  string `json:"username" binding:"required" validate:"required,min=3,max=32"`
	Password  string `json:"password" binding:"required" validate:"required,min=6"`
}

type UpdateUserInput struct {
	Id            string    `json:"id" binding:"required" validate:"required"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Language      string    `json:"language"`
	AccountType   string    `mapstructure:"accountType" bson:"accountType,omitempty" json:"accountType,omitempty"`
	LastSignInAt  time.Time `mapstructure:"lastSignInAt" bson:"lastSignInAt" json:"lastSignInAt,omitempty"`
	Organizations []string  `json:"organizations"`
	Roles         []string  `json:"roles"`
}
