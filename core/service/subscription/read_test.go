package subscription

import (
	"testing"

	"auth/core/domain"
	"auth/core/ports/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_FindSubscriptionByID(t *testing.T) {
	t.Run("find By id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		subscription, err := domain.NewSubscription()
		assert.Nil(t, err)
		assert.NotNil(t, subscription)

		subsService := mocks.NewMockSubscriptionService(ctrl)
		subsService.EXPECT().FindSubscriptionByID(subscription.Id.String()).Return(subscription, nil)

		id := subscription.Id
		subscription, err = subsService.FindSubscriptionByID(subscription.Id.String())
		assert.Nil(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, subscription.Id, id)
	})
}
