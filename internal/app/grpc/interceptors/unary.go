package interceptors

import (
	"context"
	"log"

	"auth/core/ports"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func KeycloakAuthorizationInterceptor(authServer ports.AuthServerService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		// Extract the authorization token from the gRPC metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		}
		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		}
		token := tokens[0]

		// Validate the authorization token with Keycloak
		_, err := authServer.RetrospectToken(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token")
		}

		// Call the gRPC handler with the authorized context
		return handler(ctx, req)
	}
}
