package auth

import (
	"auth/core/domain/http"
	"auth/core/ports"
	"auth/internal"
	"auth/internal/utils"
	"context"
	"github.com/datawyse/proto/golang/auth"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type gRPCAuthService struct {
	log                 *zap.Logger
	config              *internal.AppConfig
	service             ports.AuthService
	userService         ports.UserService
	subscriptionService ports.SubscriptionService
	auth.UnimplementedAuthServiceServer
}

func NewAuthGRPCService(service ports.AuthService, userService ports.UserService, log *zap.Logger, config *internal.AppConfig, subscriptionService ports.SubscriptionService) (auth.AuthServiceServer, error) {
	return &gRPCAuthService{
		log:                 log,
		config:              config,
		service:             service,
		userService:         userService,
		subscriptionService: subscriptionService,
	}, nil
}

func (svc *gRPCAuthService) CreateUser(ctx context.Context, request *auth.CreateUserRequest) (*auth.CreationResponseAuth, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "grpc.auth.get_user")
	defer span.End()

	//TODO implement me
	panic("implement me")
}

func (svc *gRPCAuthService) CreateUsers(server auth.AuthService_CreateUsersServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "grpc.auth.get_user")
	defer span.End()

	//TODO implement me
	panic("implement me")
}

func (svc *gRPCAuthService) GetUser(ctx context.Context, request *auth.GetUserRequest) (*auth.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "grpc.auth.get_user")
	defer span.End()

	userId := request.GetUserId()
	user, err := svc.service.User(ctx, userId)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]*auth.Attributes)
	for name, attr := range *user.Attributes {
		attribute := &auth.Attributes{
			Attributes: attr,
		}
		attributes[name] = attribute
	}

	createdAt := timestamppb.New(user.CreatedAt)
	updatedAt := timestamppb.New(user.UpdatedAt)

	// convert user to proto
	userProto := &auth.User{
		Id:         user.Id,
		Username:   *user.Username,
		Email:      *user.Email,
		FirstName:  *user.FirstName,
		LastName:   *user.LastName,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Attributes: attributes,
	}

	return &auth.UserResponse{
		User: userProto,
	}, nil
}

func (svc *gRPCAuthService) UpdateUser(ctx context.Context, request *auth.UpdateUserRequest) (*auth.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "grpc.auth.get_user")
	defer span.End()

	user := request.GetUser()
	svc.log.Debug("input ", zap.Any("input", user))

	// find user by id
	userInfo, err := svc.userService.User(ctx, user.Id)
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	// get attributes from user
	// get list of organization ids
	var organizationIds []string
	for name, attr := range *userInfo.Attributes {
		if name == "organizations" {
			organizationIds = attr
		}
	}

	attributes := make(map[string][]string)
	for name, attr := range user.Attributes {
		attributes[name] = attr.Attributes
	}

	// check if new organization has been created or not
	var isOrgCreated bool = false
	for name, attr := range attributes {
		if name == "organizations" {
			for _, org := range attr {
				if !utils.Contains(organizationIds, org) {
					isOrgCreated = true
				}
			}
		}
	}

	userInput := &http.UpdateUserInput{
		Id:         user.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Username:   user.Username,
		Attributes: attributes,
	}

	_, err = svc.service.UpdateUser(ctx, userInput)
	if err != nil {
		svc.log.Error("error signing up", zap.Error(err))

		// create grpc error
		return nil, err
	}

	// if new organization has been created, update subscription
	if isOrgCreated {

		// find subscription id
		var subscriptionId string
		for name, attr := range user.Attributes {
			if name == "subscription" {
				subscriptionId = attr.Attributes[0]
			}
		}

		// get user subscription
		subs, err := svc.subscriptionService.FindSubscriptionByID(ctx, subscriptionId)
		if err != nil {
			svc.log.Error("error getting user subscription", zap.Error(err))
			return nil, err
		}

		// update subscriptions
		if err := subs.AddOrganization("free"); err != nil {
			svc.log.Error("error adding organization to subscription", zap.Error(err))
			return nil, err
		}

		if _, err := svc.subscriptionService.UpdateSubscription(ctx, subs); err != nil {
			svc.log.Error("error updating subscription", zap.Error(err))
			return nil, err
		}
	}

	return &auth.UserResponse{User: user}, nil
}

//func (svc *gRPCAuthService) CreateUser(ctx context.Context, req *auth.User) (*auth.CreationResponseAuth, error) {
//	svc.log.Info("creating user")
//
//	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
//	defer cancel()
//
//	var signupInput http.SignupInput
//	signupInput.Email = req.GetEmail()
//	signupInput.FirstName = req.GetFirstName()
//	signupInput.LastName = req.GetLastName()
//	// signupInput.Password = req.GetPassword()
//	signupInput.Username = req.GetUsername()
//
//	result, err := svc.service.Signup(ctx, &signupInput)
//	if err != nil {
//		svc.log.Error("error signing up", zap.Error(err))
//
//		// create grpc error
//		return nil, err
//	}
//
//	return &auth.CreationResponseAuth{Id: result}, nil
//}

//func (svc *gRPCAuthService) UpdateUser(ctx context.Context, user *auth.User) (*auth.User, error) {
//	svc.log.Info("updating user")
//
//	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
//	defer cancel()
//
//	svc.log.Debug("input ", zap.Any("input", user))
//
//	// find user by id
//	userInfo, err := svc.userService.User(ctx, user.Id)
//	if err != nil {
//		svc.log.Error("error getting user", zap.Error(err))
//		return nil, err
//	}
//
//	// get attributes from user
//	// get list of organization ids
//	var organizationIds []string
//	for name, attr := range *userInfo.Attributes {
//		if name == "organizations" {
//			organizationIds = attr
//		}
//	}
//
//	attributes := make(map[string][]string)
//	for name, attr := range user.Attributes {
//		attributes[name] = attr.Attributes
//	}
//
//	// check if new organization has been created or not
//	var isOrgCreated bool = false
//	for name, attr := range attributes {
//		if name == "organizations" {
//			for _, org := range attr {
//				if !utils.Contains(organizationIds, org) {
//					isOrgCreated = true
//				}
//			}
//		}
//	}
//
//	userInput := &http.UpdateUserInput{
//		Id:         user.Id,
//		FirstName:  user.FirstName,
//		LastName:   user.LastName,
//		Email:      user.Email,
//		Username:   user.Username,
//		Attributes: attributes,
//	}
//
//	_, err = svc.service.UpdateUser(ctx, userInput)
//	if err != nil {
//		svc.log.Error("error signing up", zap.Error(err))
//
//		// create grpc error
//		return nil, err
//	}
//
//	// if new organization has been created, update subscription
//	if isOrgCreated {
//
//		// find subscription id
//		var subscriptionId string
//		for name, attr := range user.Attributes {
//			if name == "subscription" {
//				subscriptionId = attr.Attributes[0]
//			}
//		}
//
//		// get user subscription
//		subs, err := svc.subscriptionService.FindSubscriptionByID(ctx, subscriptionId)
//		if err != nil {
//			svc.log.Error("error getting user subscription", zap.Error(err))
//			return nil, err
//		}
//
//		// update subscriptions
//		if err := subs.AddOrganization("free"); err != nil {
//			svc.log.Error("error adding organization to subscription", zap.Error(err))
//			return nil, err
//		}
//
//		if _, err := svc.subscriptionService.UpdateSubscription(ctx, subs); err != nil {
//			svc.log.Error("error updating subscription", zap.Error(err))
//			return nil, err
//		}
//	}
//
//	return user, nil
//}
//
//func (svc *gRPCAuthService) CreateUsers(server auth.AuthService_CreateUsersServer) error {
//	panic("implement me")
//}
