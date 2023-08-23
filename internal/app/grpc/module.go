package grpc

import (
	"auth/internal/app/grpc/clients"
	"auth/internal/app/grpc/server"

	"go.uber.org/fx"
)

var Module = fx.Options(
	server.Module, clients.Module,
)
