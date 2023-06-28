package domain

import (
	"auth/core/domain/system"
)

type Subscription struct {
	system.BaseCollection `mapstructure:",squash" bson:"inline"`

	Organizations        int8 `mapstructure:"organization" bson:"organization" json:"organization" validate:"required"`                            // 1
	ActiveOrganizations  int8 `mapstructure:"active_organization" bson:"active_organization" json:"active_organization" validate:"required"`       // 0
	OrganizationProjects int8 `mapstructure:"organization_projects" bson:"organization_projects" json:"organization_projects" validate:"required"` // 0
	PaidProjects         int8 `mapstructure:"paid_projects" bson:"paid_projects" json:"paid_projects" validate:"required"`                         // 0
	FreeProjects         int8 `mapstructure:"free_projects" bson:"free_projects" json:"free_projects" validate:"required"`                         // 3
	ActiveFreeProjects   int8 `mapstructure:"active_free_projects" bson:"active_free_projects" json:"active_free_projects" validate:"required"`    // 5
	PausedFreeProjects   int8 `mapstructure:"paid_free_projects" bson:"paid_free_projects" json:"paid_free_projects" validate:"required"`          // 0
	ProProjects          int8 `mapstructure:"pro_projects" bson:"pro_projects" json:"pro_projects" validate:"required"`                            // 0
	EnterpriseProjects   int8 `mapstructure:"enterprise_projects" bson:"enterprise_projects" json:"enterprise_projects" validate:"required"`       // 0
}

func NewSubscription() (*Subscription, error) {
	subs := &Subscription{
		BaseCollection:       *system.NewBaseCollection(),
		Organizations:        1,
		ActiveOrganizations:  0,
		OrganizationProjects: 1,
		PaidProjects:         0,
		FreeProjects:         3,
		ActiveFreeProjects:   0,
		PausedFreeProjects:   0,
		ProProjects:          0,
		EnterpriseProjects:   0,
	}

	return subs, nil
}
