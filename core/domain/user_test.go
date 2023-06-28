package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserModel(t *testing.T) {
	t.Run("Test User Model", func(t *testing.T) {
		user, err := NewUser(uuid.New())
		if err != nil {
			t.Errorf("error creating user model: %v", err)
		}

		if user.Id == uuid.Nil {
			t.Errorf("error creating user model: %v", err)
		}

		if user.CreatedAt.IsZero() {
			t.Errorf("error creating user model: %v", err)
		}
		if user.UpdatedAt.IsZero() {
			t.Errorf("error creating user model: %v", err)
		}

		if user.CreatedAt != user.UpdatedAt {
			t.Errorf("error creating user model: %v", err)
		}
	})
}
