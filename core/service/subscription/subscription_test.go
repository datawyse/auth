package subscription

import (
	"context"
	"testing"
	"time"

	"auth/core/ports/mocks"
	tests2 "auth/tests"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSubscriptionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	defer ctrl.Finish()

	t.Run("subscription instance", func(t *testing.T) {
		log := tests2.NewTestLog(t)
		config := tests2.NewTestConfig(t)
		repo := mocks.NewMockSubscriptionRepository(ctrl)
		uuidService := mocks.NewMockUUIDService(ctrl)

		subsService, err := NewSubscriptionService(ctx, log, config, repo, uuidService)

		assert.Nil(t, err)
		assert.NotNil(t, subsService)
	})
}
