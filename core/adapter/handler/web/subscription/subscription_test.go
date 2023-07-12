package subscription

import (
	"context"
	"testing"
	"time"

	"auth/core/ports/mocks"
	"auth/tests"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSubscriptionsController(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	defer ctrl.Finish()

	t.Run("should create NewSubscriptionsController instance", func(t *testing.T) {
		api := tests.SetUpRouterGroup()
		log := tests.NewTestLog(t)
		config := tests.NewTestConfig(t)
		validate := validator.New()
		service := mocks.NewMockSubscriptionService(ctrl)
		userService := mocks.NewMockUserService(ctrl)

		subscription, err := NewSubscriptionsController(ctx, api, log, config, validate, service, userService)
		assert.Nil(t, err)
		assert.NotNil(t, subscription)
	})
}
