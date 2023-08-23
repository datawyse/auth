package auth

import (
	"auth/core/ports"
	"auth/internal"
	"context"
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

func (svc *gRPCSubscriptionService) GetSubscription(ctx context.Context, request *auth.SubscriptionRequest) (*auth.SubscriptionResponse, error) {
	svc.log.Info("getting subscription")

	subs, err := svc.subscriptionService.FindSubscriptionByID(ctx, request.SubscriptionId)
	if err != nil {
		return nil, err
	}

	subscription := &auth.Subscription{
		Id:                       subs.Id.String(),
		Organizations:            subs.Organizations,
		ActiveOrganizations:      subs.ActiveOrganizations,
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
	}

	return &auth.SubscriptionResponse{
		Subscription: subscription,
	}, nil
}

func (svc *gRPCSubscriptionService) GetUserSubscription(request *auth.GetUserSubscriptionRequest, server auth.SubscriptionService_GetUserSubscriptionServer) error {
	//TODO implement me
	return status.Errorf(codes.Unimplemented, "method GetUserSubscription not implemented")
}
