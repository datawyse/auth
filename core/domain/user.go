package domain

import (
	"time"

	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
)

type AccountType int

func (a AccountType) String() string {
	return [...]string{"", "admin", "organization", "guest"}[a]
}

func (a AccountType) EnumIndex() int {
	return int(a)
}

const (
	ACCOUNT_UNSPECIFIED AccountType = iota
	ACCOUNT_ADMIN
	ACCOUNT_ORGANIZATION
	ACCOUNT_GUEST
)

// User - user entity
type User struct {
	// store the base collection as an embedded struct
	system.BaseCollection `mapstructure:",squash" bson:"inline"`

	Id           uuid.UUID `mapstructure:"id" bson:"_id" json:"id"`
	LastSignInAt time.Time `mapstructure:"lastSignInAt" bson:"lastSignInAt" json:"lastSignInAt,omitempty"`
}

type KeycloakUser = gocloak.User
type SystemUser = User
type UserInfo struct {
	Id string `json:"id"`
	KeycloakUser
	SystemUser
}

// TypeName implements openapi.Typer interface for Fruit.
func (info *UserInfo) TypeName() string {
	return "UserInfo"
}

// TypeName implements openapi.Typer interface for Fruit.
func (user *User) TypeName() string {
	return "User"
}

func NewUser(authId uuid.UUID) (*User, error) {
	user := &User{
		BaseCollection: *system.NewBaseCollection(),
		Id:             authId,
		LastSignInAt:   time.Now(),
	}

	return user, nil
}

func NewUserInfo(keycloakUser *gocloak.User, systemUser *User) *UserInfo {
	// remove unwanted fields from keycloak user
	keycloakUser.CreatedTimestamp = nil
	keycloakUser.FederationLink = nil
	keycloakUser.DisableableCredentialTypes = nil
	keycloakUser.Totp = nil

	return &UserInfo{
		Id:           systemUser.Id.String(),
		SystemUser:   *systemUser,
		KeycloakUser: *keycloakUser,
	}
}

func (info *UserInfo) ToUser() (*User, error) {
	id, err := system.ToUUID(info.Id)
	if err != nil {
		return nil, system.ErrInvalidInput
	}

	user := &User{
		BaseCollection: info.BaseCollection,
		Id:             id,
		LastSignInAt:   time.Now(),
	}

	return user, nil
}

func (info *UserInfo) toGocloakUser() *gocloak.User {
	return &info.KeycloakUser
}

// GetId returns the model id.
func (user *User) GetId() uuid.UUID {
	return user.Id
}

// SetId sets the model id to the provided string value.
func (user *User) SetId(id uuid.UUID) {
	user.Id = id
}
