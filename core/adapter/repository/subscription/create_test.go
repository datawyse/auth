package subscription

import (
	"context"
	"testing"
	"time"

	"auth/core/domain"
	"auth/internal/db/mongodb"
	tests2 "auth/tests"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRepository_CreateSubscription(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("creating", func(mt *mtest.T) {
		log := tests2.NewTestLog(t)
		config := tests2.NewTestConfig(t)
		mockedDb := &mongodb.MongoDb{Db: mt.DB, Client: mt.Client}
		subsCollection := mt.DB.Collection("subscriptions")

		subs, err := domain.NewSubscription()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &Repository{
			db:         mockedDb,
			log:        log,
			ctx:        ctx,
			config:     config,
			Collection: subsCollection,
		}

		subsId, err := repo.CreateSubscription(ctx, subs)

		assert.Nil(t, err)
		assert.NotNil(t, subsId)
		assert.Equal(t, subs.Id, subsId)
	})
}
