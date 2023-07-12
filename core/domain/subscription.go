package domain

import (
	"auth/core/domain/system"
)

type Subscription struct {
	system.BaseCollection `mapstructure:",squash" bson:"inline"`

	Organizations            int32 `mapstructure:"organizations" bson:"organizations" json:"organizations" validate:"required"`                                  // 1
	ActiveOrganizations      int32 `mapstructure:"activeOrganizations" bson:"activeOrganizations" json:"activeOrganizations" validate:"required"`                // 0
	OrganizationProjects     int32 `mapstructure:"organizationProjects" bson:"organizationProjects" json:"organizationProjects" validate:"required"`             // 3
	PaidProjects             int32 `mapstructure:"paidProjects" bson:"paidProjects" json:"paidProjects" validate:"required"`                                     // 0
	FreeProjects             int32 `mapstructure:"freeProjects" bson:"freeProjects" json:"freeProjects" validate:"required"`                                     // 3
	EnterpriseProjects       int32 `mapstructure:"enterpriseProjects" bson:"enterpriseProjects" json:"enterpriseProjects" validate:"required"`                   // 0
	ActiveFreeProjects       int32 `mapstructure:"activeFreeProjects" bson:"activeFreeProjects" json:"activeFreeProjects" validate:"required"`                   // 5
	PausedFreeProjects       int32 `mapstructure:"pausedFreeProjects" bson:"pausedFreeProjects" json:"pausedFreeProjects" validate:"required"`                   // 0
	ActivePaidProjects       int32 `mapstructure:"activePaidProjects" bson:"activePaidProjects" json:"activePaidProjects" validate:"required"`                   // 0                            // 0
	PausedPaidProjects       int32 `mapstructure:"pausedPaidProjects" bson:"pausedPaidProjects" json:"pausedPaidProjects" validate:"required"`                   // 0
	ActiveEnterpriseProjects int32 `mapstructure:"activeEnterpriseProjects" bson:"activeEnterpriseProjects" json:"activeEnterpriseProjects" validate:"required"` // 0
	PausedEnterpriseProjects int32 `mapstructure:"pausedEnterpriseProjects" bson:"pausedEnterpriseProjects" json:"pausedEnterpriseProjects" validate:"required"` // 0
}

func NewSubscription(subscriptionType string) (*Subscription, error) {
	subs := &Subscription{
		BaseCollection:           *system.NewBaseCollection(),
		Organizations:            3,
		ActiveOrganizations:      0,
		OrganizationProjects:     2,
		PaidProjects:             0,
		FreeProjects:             2,
		EnterpriseProjects:       0,
		ActiveFreeProjects:       0,
		PausedFreeProjects:       0,
		ActivePaidProjects:       0,
		PausedPaidProjects:       0,
		ActiveEnterpriseProjects: 0,
		PausedEnterpriseProjects: 0,
	}

	return subs, nil
}

func (s *Subscription) AddOrganization(subscriptionType string) error {
	if s.ActiveOrganizations < s.Organizations {
		s.ActiveOrganizations++
		return nil
	}

	return system.ErrInvalidInput
}
