package user

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo Repository) UpdateUser(ctx context.Context, input *domain.User) (*domain.User, error) {
	repo.log.Info("updating user")

	ctx, cancel := context.WithTimeout(repo.ctx, time.Duration(repo.config.DatabaseTimeout)*time.Second)
	defer cancel()

	// update user
	filter := bson.D{{"_id", input.Id}}
	var updateValues bson.D
	if len(input.Roles) > 0 {
		updateValues = append(updateValues, bson.E{Key: "roles", Value: input.Roles})
	}
	if len(input.Organizations) > 0 {
		updateValues = append(updateValues, bson.E{Key: "organizations", Value: input.Organizations})
	}
	if input.Language != "" {
		updateValues = append(updateValues, bson.E{Key: "language", Value: input.Language})
	}

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
