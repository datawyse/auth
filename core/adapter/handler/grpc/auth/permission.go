package auth

import (
	"github.com/team-management-io/proto/golang/auth"
)

type gRPCPermissionService struct {
	auth.UnimplementedPermissionServiceServer
}

func NewPermissionGRPCService() auth.PermissionServiceServer {
	return &gRPCPermissionService{}
}

// GetUserPermissions returns the permissions of a user
func (s *gRPCPermissionService) GetUserPermissions(req *auth.UserIdRequestAuth, perm auth.PermissionService_GetUserPermissionsServer) error {
	panic("not implemented") // TODO: Implement
}
