package health

import (
	"context"
	"time"

	"auth/core/domain"
)

// ReadHealth returns the health.
func (s *Service) ReadHealth(ctx context.Context) (*domain.Health, error) {
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")

	timezone := time.Now().Format("UTC")
	datetime := currentDate + " " + currentTime + " " + timezone

	isMongoHealthy, err := s.mongo.IsHealthy()
	if err != nil {
		return nil, err
	}

	var mongoStatus = "ERROR"
	if isMongoHealthy {
		mongoStatus = "OK"
	}

	isRedisHealthy, err := s.redis.IsHealthy()
	if err != nil {
		return nil, err
	}

	var redisStatus = "ERROR"
	if isRedisHealthy {
		redisStatus = "OK"
	}

	return &domain.Health{
		Status:   "OK",
		Message:  "System is healthy",
		Version:  "v1.0.0",
		Date:     currentDate,
		Datetime: datetime,
		Time:     currentTime,
		Timezone: timezone,
		Revision: "1",
		Redis:    mongoStatus,
		Mongo:    redisStatus,
	}, nil
}
