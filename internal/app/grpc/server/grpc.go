package server

import (
	"auth/internal/app/grpc/interceptors"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"log"
	"net"

	"auth/core/ports"
	"auth/internal"
	"github.com/datawyse/proto/golang/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

func NewGRPCServer(config *internal.AppConfig, log *zap.Logger, authServiceServer auth.AuthServiceServer, subscriptionServiceServer auth.SubscriptionServiceServer, authServer ports.AuthServerService) (*Server, error) {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor(), interceptors.UnaryInterceptor(authServer, log)))
	opts = append(opts, grpc.ChainStreamInterceptor(otelgrpc.StreamServerInterceptor(), interceptors.StreamInterceptor(authServer, log)))
	grpcServer := grpc.NewServer(opts...)

	return &Server{
		port:                      config.GRPCPort,
		server:                    grpcServer,
		authServiceServer:         authServiceServer,
		subscriptionServiceServer: subscriptionServiceServer,
	}, nil
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

var Module = fx.Provide(NewGRPCServer)
