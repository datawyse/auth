package auth

import (
	"context"

	"github.com/datawyse/proto/golang/auth"
)

type gRPCRoleService struct {
	auth.UnimplementedRoleServiceServer
}

func NewRoleGRPCService() auth.RoleServiceServer {
	return &gRPCRoleService{}
}

// GetUserRoles returns the roles of a user
func (s *gRPCRoleService) GetUserRoles(ctx context.Context, req *auth.UserIdRequestAuth) (*auth.Role, error) {
	panic("not implemented") // TODO: Implement
}
