package interceptors

import (
	"auth/core/ports"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryInterceptor(authServer ports.AuthServerService, log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Info("--> unary interceptor: ", zap.String("method: ", info.FullMethod))

		//// Extract the authorization token from the gRPC metadata
		//md, ok := metadata.FromIncomingContext(ctx)
		//if !ok {
		//	return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//tokens := md.Get("authorization")
		//
		//if len(tokens) == 0 {
		//	return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//bearerToken := strings.Split(tokens[0], " ")
		//if len(bearerToken) != 2 {
		//	return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//if strings.ToLower(bearerToken[0]) != "bearer" {
		//	return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//token := bearerToken[1]
		//if len(token) == 0 {
		//	return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//// Validate the authorization token with Keycloak
		//_, err := authServer.VerifyToken(ctx, token)
		//if err != nil {
		//	return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token")
		//}

		// Call the gRPC handler with the authorized context
		return handler(ctx, req)
	}
}
