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

	Id            uuid.UUID   `mapstructure:"id" bson:"_id" json:"id"`
	Verified      bool        `mapstructure:"verified" bson:"verified" json:"verified"`
	Active        bool        `mapstructure:"active" bson:"active" json:"active"`
	Avatar        string      `mapstructure:"avatar" bson:"avatar,omitempty" json:"avatar"`
	Language      string      `mapstructure:"language" bson:"language,omitempty" json:"language"`
	AccountType   string      `mapstructure:"accountType" bson:"accountType,omitempty" json:"accountType"`
	Roles         []uuid.UUID `mapstructure:"roles" bson:"roles,omitempty" json:"roles"`
	Organizations []uuid.UUID `mapstructure:"organizations" bson:"organizations,omitempty" json:"organizations"`
	Subscription  uuid.UUID   `mapstructure:"subscription" bson:"subscription" json:"subscription"`
	LastSignInAt  time.Time   `mapstructure:"lastSignInAt" bson:"lastSignInAt" json:"lastSignInAt,omitempty"`
}

type KeycloakUser = gocloak.User
type SystemUser = User
type UserInfo struct {
	Id string `json:"id"`
	KeycloakUser
	SystemUser
}

func NewUser(authId uuid.UUID) (*User, error) {
	user := &User{
		BaseCollection: *system.NewBaseCollection(),
		Id:             authId,
		Verified:       false,
		Active:         true,
		LastSignInAt:   time.Now(),
	}

	return user, nil
}

func NewUserInfo(keycloakUser *gocloak.User, systemUser *User) *UserInfo {
	// remove unwanted fields from keycloak user
	keycloakUser.Attributes = nil
	keycloakUser.CreatedTimestamp = nil
	keycloakUser.EmailVerified = nil
	keycloakUser.Enabled = nil
	keycloakUser.FederationLink = nil
	keycloakUser.DisableableCredentialTypes = nil
	keycloakUser.Totp = nil

	return &UserInfo{
		Id:           systemUser.Id.String(),
		SystemUser:   *systemUser,
		KeycloakUser: *keycloakUser,
	}
}

// GetId returns the model id.
func (user *User) GetId() uuid.UUID {
	return user.Id
}

// SetId sets the model id to the provided string value.
func (user *User) SetId(id uuid.UUID) {
	user.Id = id
}

// SetSubscriptionId sets the model subscription to the provided string value.
func (user *User) SetSubscriptionId(subscription uuid.UUID) {
	user.Subscription = subscription
}
