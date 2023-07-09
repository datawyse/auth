package auth

import (
	"context"
	"fmt"
	"time"

	"auth/core/domain/http"
	"auth/core/ports"
	"auth/internal"

	"github.com/datawyse/proto/golang/auth"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCAuthService struct {
	service     ports.AuthService
	userService ports.UserService
	log         *zap.Logger
	config      *internal.AppConfig
	auth.UnimplementedAuthServiceServer
}

func NewAuthGRPCService(service ports.AuthService, userService ports.UserService, log *zap.Logger, config *internal.AppConfig) auth.AuthServiceServer {
	return &gRPCAuthService{
		log:         log,
		config:      config,
		service:     service,
		userService: userService,
	}
}

func (svc *gRPCAuthService) GetUser(ctx context.Context, req *auth.UserIdRequestAuth) (*auth.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	userId := req.GetUserId()
	user, err := svc.service.User(ctx, userId)
	if err != nil {
		return nil, err
	}

	// var roles []string
	// for _, role := range user.Roles {
	// 	roles = append(roles, role.String())
	// }

	// var organizations []string
	// for _, organization := range user.Organizations {
	// 	organizations = append(organizations, organization.String())
	// }

	createdAt := timestamppb.New(user.CreatedAt)
	updatedAt := timestamppb.New(user.UpdatedAt)

	// convert user to proto
	userProto := &auth.User{
		Id:        user.Id,
		Username:  *user.Username,
		Email:     *user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		// Organizations: organizations,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return userProto, nil
}

func (svc *gRPCAuthService) CreateUser(ctx context.Context, req *auth.User) (*auth.CreationResponseAuth, error) {
	svc.log.Info("creating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	var signupInput http.SignupInput
	signupInput.Email = req.GetEmail()
	signupInput.FirstName = req.GetFirstName()
	signupInput.LastName = req.GetLastName()
	signupInput.Password = req.GetPassword()
	signupInput.Username = req.GetUsername()

	result, err := svc.service.Signup(ctx, &signupInput)
	if err != nil {
		svc.log.Error("error signing up", zap.Error(err))

		// create grpc error
		return nil, err
	}

	return &auth.CreationResponseAuth{Id: result}, nil
}

func (svc *gRPCAuthService) UpdateUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	svc.log.Info("updating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.GRPCServiceTimeout)*time.Second)
	defer cancel()

	svc.log.Debug("input ", zap.Any("input", user))

	userInput := &http.UpdateUserInput{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		// Language:      user.Language,
		// AccountType:   user.AccountType.String(),
		// Roles:         user.Roles,
		// Organizations: user.Organizations,
	}

	result, err := svc.service.UpdateUser(ctx, userInput)
	if err != nil {
		svc.log.Error("error signing up", zap.Error(err))

		// create grpc error
		return nil, err
	}

	fmt.Println("result: ", result)
	return user, nil
}

func (svc *gRPCAuthService) CreateUsers(server auth.AuthService_CreateUsersServer) error {
	panic("implement me")
}
