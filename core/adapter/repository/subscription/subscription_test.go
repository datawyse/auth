package subscription

import (
	"context"
	"testing"
	"time"

	"auth/internal/db/mongodb"
	tests2 "auth/tests"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewSubscriptionRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("create subscription", func(mt *mtest.T) {
		log := tests2.NewTestLog(t)
		config := tests2.NewTestConfig(t)
		mockedDb := &mongodb.MongoDb{Db: mt.DB, Client: mt.Client}

		repo, err := NewSubscriptionRepository(ctx, log, config, mockedDb)
		assert.Nil(t, err)
		assert.NotNil(t, repo)
	})
}
