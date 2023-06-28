package tests

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func NewTestRedisDb(t *testing.T) *redis.Client {
	db, _ := redismock.NewClientMock()

	return db
}
