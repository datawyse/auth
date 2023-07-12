package auth

import (
	"context"

	"auth/core/ports"
	"auth/internal"

	"github.com/datawyse/proto/golang/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCSubscriptionService struct {
	log                 *zap.Logger
	config              *internal.AppConfig
	subscriptionService ports.SubscriptionService
	auth.UnimplementedSubscriptionServiceServer
}

func NewSubscriptionGRPCService(log *zap.Logger, config *internal.AppConfig, subscriptionService ports.SubscriptionService) (auth.SubscriptionServiceServer, error) {
	return &gRPCSubscriptionService{
		log:                 log,
		config:              config,
		subscriptionService: subscriptionService,
	}, nil
}

func (svc *gRPCSubscriptionService) GetSubscription(ctx context.Context, request *auth.SubscriptionRequest) (*auth.Subscription, error) {
	svc.log.Info("getting subscription")

	subs, err := svc.subscriptionService.FindSubscriptionByID(ctx, request.SubscriptionId)
	if err != nil {
		return nil, err
	}

	return &auth.Subscription{
		Id:                       subs.Id.String(),
		Organizations:            subs.Organizations,
		ActiveOrganization:       subs.ActiveOrganizations,
		OrganizationProjects:     subs.OrganizationProjects,
		PaidProjects:             subs.PaidProjects,
		FreeProjects:             subs.FreeProjects,
		EnterpriseProjects:       subs.EnterpriseProjects,
		ActiveFreeProjects:       subs.ActiveFreeProjects,
		PausedFreeProjects:       subs.PausedFreeProjects,
		ActivePaidProjects:       subs.ActivePaidProjects,
		PausedPaidProjects:       subs.PausedPaidProjects,
		ActiveEnterpriseProjects: subs.ActiveEnterpriseProjects,
		PausedEnterpriseProjects: subs.PausedEnterpriseProjects,
		CreatedAt:                timestamppb.New(subs.CreatedAt),
		UpdatedAt:                timestamppb.New(subs.UpdatedAt),
	}, nil
}

func (svc *gRPCSubscriptionService) GetUserSubscription(*auth.UserIdRequestAuth, auth.SubscriptionService_GetUserSubscriptionServer) error {
	return status.Errorf(codes.Unimplemented, "method GetUserSubscription not implemented")
}
