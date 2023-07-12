package subscription

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"auth/core/ports/mocks"
	"auth/tests"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestController_GetSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		api := tests.SetUpRouterGroup()
		log := tests.NewTestLog(t)
		config := tests.NewTestConfig(t)
		validate := validator.New()
		service := mocks.NewMockSubscriptionService(ctrl)
		userService := mocks.NewMockUserService(ctrl)

		subscription, err := NewSubscriptionsController(ctx, api, log, config, validate, service, userService)
		assert.Nil(t, err)
		assert.NotNil(t, subscription)

		router := tests.SetUpRouter()
		router.GET("/api/v1/auth/subscriptions", subscription.GetSubscriptions)
		req, _ := http.NewRequest("GET", "/api/v1/auth/subscriptions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
