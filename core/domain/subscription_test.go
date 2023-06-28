package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestSubscription(t *testing.T) {
	t.Run("should return subscription if valid", func(t *testing.T) {
		subscription, err := NewSubscription()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if subscription == nil {
			t.Error("expected subscription, got nil")
		}

		if subscription.Id == uuid.Nil {
			t.Error("expected subscription id, got empty string")
		}
	})
}
