package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (repo Repository) UpdateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	repo.log.Info("updating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(repo.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(repo.config.ServiceName).Start(ctx, "repository.subscription.update_subscription")
	defer span.End()

	// update subscription in db
	updatedSubscription := bson.D{
		{"$set", bson.D{
			{"organizations", input.Organizations},
			{"activeOrganizations", input.ActiveOrganizations},
			{"organizationProjects", input.OrganizationProjects},
			{"paidProjects", input.PaidProjects},
			{"freeProjects", input.FreeProjects},
			{"enterpriseProjects", input.EnterpriseProjects},
			{"activeFreeProjects", input.ActiveFreeProjects},
			{"pausedFreeProjects", input.ActiveFreeProjects},
			{"activePaidProjects", input.ActivePaidProjects},
			{"pausedPaidProjects", input.ActivePaidProjects},
			{"activeEnterpriseProjects", input.ActiveEnterpriseProjects},
			{"pausedEnterpriseProjects", input.ActiveEnterpriseProjects},
			{"updatedAt", input.UpdatedAt},
		}},
	}
	_, err := repo.db.Db.Collection("subscriptions").UpdateOne(ctx, bson.D{{"_id", input.Id}}, updatedSubscription)
	if err != nil {
		repo.log.Error("error updating subscription", zap.Error(err))
		return uuid.Nil, err
	}

	return input.Id, nil
}
