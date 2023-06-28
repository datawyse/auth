package interceptors

import (
	"log"

	"google.golang.org/grpc"
)

func StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}
