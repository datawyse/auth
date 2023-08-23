package interceptors

import (
	"auth/core/ports"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StreamInterceptor(authServer ports.AuthServerService, log *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Info("--> stream interceptor: ", zap.String("method: ", info.FullMethod))

		//// Extract the authorization token from the gRPC metadata
		//md, ok := metadata.FromIncomingContext(stream.Context())
		//if !ok {
		//	return status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//tokens := md.Get("authorization")
		//
		//bearerToken := strings.Split(tokens[0], " ")
		//if len(bearerToken) != 2 {
		//	return status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//if strings.ToLower(bearerToken[0]) != "bearer" {
		//	return status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//token := bearerToken[1]
		//if len(token) == 0 {
		//	return status.Errorf(codes.Unauthenticated, "authorization token not found")
		//}
		//
		//// Validate the authorization token with Keycloak
		//_, err := authServer.VerifyToken(stream.Context(), token)
		//if err != nil {
		//	return status.Errorf(codes.Unauthenticated, "invalid authorization token")
		//}

		// Call the gRPC handler with the authorized context
		return handler(srv, stream)
	}
}
