package user

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo Repository) UpdateUser(ctx context.Context, input *domain.User) (*domain.User, error) {
	repo.log.Info("updating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(repo.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(repo.config.ServiceName).Start(ctx, "repository.user.update_user")
	defer span.End()

	// update user
	filter := bson.D{{"_id", input.Id}}
	var updateValues bson.D
	if (input.LastSignInAt != time.Time{}) {
		updateValues = append(updateValues, bson.E{Key: "lastSignInAt", Value: input.LastSignInAt})
	}
	updateValues = append(updateValues, bson.E{Key: "updatedAt", Value: time.Now()})

	if len(updateValues) == 0 {
		return input, nil
	}

	update := bson.D{{"$set", updateValues}}
	result, err := repo.db.Db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		repo.log.Info("user not found")
		return nil, system.ErrInvalidInput
	}

	return input, nil
}
