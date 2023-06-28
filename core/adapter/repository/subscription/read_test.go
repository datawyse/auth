package subscription

import (
	"context"
	"testing"
	"time"

	"auth/core/domain"
	"auth/internal/db/mongodb"
	tests2 "auth/tests"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRepository_FindSubscriptionByID(t *testing.T) {
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
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "auth.subscription", mtest.FirstBatch, bson.D{
			{"_id", subs.Id},
			{"organizations", subs.Organizations},
			{"free_projects", subs.FreeProjects},
			{"organization_projects", subs.OrganizationProjects},
		}))

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

		freeProjectsCount := subs.FreeProjects
		organizationProjectsCount := subs.OrganizationProjects
		subs, err = repo.FindSubscriptionByID(ctx, subsId)

		assert.Nil(t, err)
		assert.NotNil(t, subsId)
		assert.Equal(t, subsId, subs.Id)
		assert.Equal(t, freeProjectsCount, subs.FreeProjects)
		assert.Equal(t, organizationProjectsCount, subs.OrganizationProjects)
	})
}
