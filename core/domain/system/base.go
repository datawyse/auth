package system

import (
	"time"

	"github.com/google/uuid"
)

type BaseCollection struct {
	isNewFlag bool

	Id        uuid.UUID `mapstructure:"id" bson:"_id" json:"id"`
	IsDeleted bool      `mapstructure:"isDeleted" bson:"isDeleted" json:"isDeleted,omitempty"`
	Remarks   string    `mapstructure:"remarks" bson:"remarks" json:"remarks,omitempty"`
	CreatedAt time.Time `mapstructure:"createdAt" bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time `mapstructure:"updatedAt" bson:"updatedAt" json:"updatedAt,omitempty"`
}

func NewBaseCollection() *BaseCollection {
	now := time.Now()

	return &BaseCollection{
		isNewFlag: true,
		Id:        uuid.New(),
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
		Remarks:   "system",
	}
}

// HasId returns whether the model has a nonzero id.
func (baseCollection *BaseCollection) HasId() bool {
	return baseCollection.Id != uuid.Nil
}

// GetId returns the model id.
func (baseCollection *BaseCollection) GetId() uuid.UUID {
	return baseCollection.Id
}

// SetId sets the model id to the provided string value.
func (baseCollection *BaseCollection) SetId(id uuid.UUID) {
	baseCollection.Id = id
}

// MarkAsNew sets the model isNewFlag enforcing [m.IsNew()] to be true.
func (baseCollection *BaseCollection) MarkAsNew() {
	baseCollection.isNewFlag = true
}

// UnmarkAsNew resets the model isNewFlag.
func (baseCollection *BaseCollection) UnmarkAsNew() {
	baseCollection.isNewFlag = false
}

// IsNew indicates what type of db query (insert or update)
// should be used with the model instance.
func (baseCollection *BaseCollection) IsNew() bool {
	return baseCollection.isNewFlag || !baseCollection.HasId()
}

// GetCreated returns the model Created datetime.
func (baseCollection *BaseCollection) GetCreated() time.Time {
	return baseCollection.CreatedAt
}

// GetUpdated returns the model Updated datetime.
func (baseCollection *BaseCollection) GetUpdated() time.Time {
	return baseCollection.UpdatedAt
}

// RefreshId generates and sets a new model id.
// The generated id is a cryptographically random 15 characters length string.
func (baseCollection *BaseCollection) RefreshId() {
	if baseCollection.Id == uuid.Nil { // no previous id
		baseCollection.MarkAsNew()
	}
	baseCollection.Id = uuid.New()
}

// RefreshCreated updates the model Created field.go with the current datetime.
func (baseCollection *BaseCollection) RefreshCreated() {
	baseCollection.CreatedAt = time.Now()
}

// RefreshUpdated updates the model Updated field.go with the current datetime.
func (baseCollection *BaseCollection) RefreshUpdated() {
	baseCollection.UpdatedAt = time.Now()
}

// SetCreatedAt - created at
func (baseCollection *BaseCollection) SetCreatedAt() {
	baseCollection.CreatedAt = time.Now()
}

// SetUpdatedAt - created at
func (baseCollection *BaseCollection) SetUpdatedAt() {
	baseCollection.UpdatedAt = time.Now()
}
