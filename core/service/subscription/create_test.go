package subscription

import (
	"testing"

	"auth/core/domain"
	"auth/core/ports/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateSubscription(t *testing.T) {

	t.Run("creating subscription", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		subscription, err := domain.NewSubscription()
		assert.Nil(t, err)
		assert.NotNil(t, subscription)

		subsService := mocks.NewMockSubscriptionService(ctrl)
		subsService.EXPECT().CreateSubscription(subscription).Return(subscription.Id, nil)

		subsId, err := subsService.CreateSubscription(subscription)
		assert.Nil(t, err)
		assert.NotNil(t, subsId)
		assert.Equal(t, subsId, subscription.Id)
	})
}
