package grpc

import (
	"fmt"
	"log"
	"net"

	"auth/internal"
	"auth/internal/app/grpc/interceptors"

	"github.com/team-management-io/proto/golang/auth"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	port                    string
	server                  *grpc.Server
	authServiceServer       auth.AuthServiceServer
	roleServiceServer       auth.RoleServiceServer
	permissionServiceServer auth.PermissionServiceServer
}

func NewGRPCServer(config *internal.AppConfig, authServiceServer auth.AuthServiceServer, roleServiceServer auth.RoleServiceServer, permissionServiceServer auth.PermissionServiceServer) *Server {
	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(interceptors.UnaryInterceptor))
	opts = append(opts, grpc.StreamInterceptor(interceptors.StreamInterceptor))

	grpcServer := grpc.NewServer(opts...)

	return &Server{
		port:                    config.GRPCPort,
		server:                  grpcServer,
		authServiceServer:       authServiceServer,
		roleServiceServer:       roleServiceServer,
		permissionServiceServer: permissionServiceServer,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	auth.RegisterAuthServiceServer(s.server, s.authServiceServer)
	auth.RegisterRoleServiceServer(s.server, s.roleServiceServer)
	auth.RegisterPermissionServiceServer(s.server, s.permissionServiceServer)

	reflection.Register(s.server)

	return s.server.Serve(lis)
}

func (s *Server) Stop() error {
	s.server.Stop()
	return nil
}

var ServerModule = fx.Provide(NewGRPCServer)
