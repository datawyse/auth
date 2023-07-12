package grpc

import (
	"fmt"
	"log"
	"net"

	"auth/core/ports"
	"auth/internal"
	"auth/internal/app/grpc/interceptors"

	"github.com/datawyse/proto/golang/auth"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	port                      string
	server                    *grpc.Server
	authServer                ports.AuthServerService
	authServiceServer         auth.AuthServiceServer
	subscriptionServiceServer auth.SubscriptionServiceServer
}

func NewGRPCServer(config *internal.AppConfig, authServiceServer auth.AuthServiceServer, subscriptionServiceServer auth.SubscriptionServiceServer, authServer ports.AuthServerService) *Server {
	var opts []grpc.ServerOption

	opts = append(opts, grpc.StreamInterceptor(interceptors.StreamInterceptor))
	// opts = append(opts, grpc.UnaryInterceptor(interceptors.KeycloakAuthorizationInterceptor(authServer)))

	grpcServer := grpc.NewServer(opts...)

	return &Server{
		port:                      config.GRPCPort,
		server:                    grpcServer,
		authServiceServer:         authServiceServer,
		subscriptionServiceServer: subscriptionServiceServer,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	auth.RegisterAuthServiceServer(s.server, s.authServiceServer)
	auth.RegisterSubscriptionServiceServer(s.server, s.subscriptionServiceServer)

	reflection.Register(s.server)

	return s.server.Serve(lis)
}

func (s *Server) Stop() error {
	s.server.Stop()
	return nil
}

var ServerModule = fx.Provide(NewGRPCServer)
