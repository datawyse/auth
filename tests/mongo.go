package tests

import (
	"testing"

	"auth/internal/db/mongodb"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func NewTestMongoDb(t *testing.T) *mongodb.MongoDb {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mockedDb := &mongodb.MongoDb{
		Db:     mt.DB,
		Client: mt.Client,
	}

	return mockedDb
}
