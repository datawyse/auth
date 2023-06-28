package handler

import (
	authgrpc "auth/core/adapter/handler/grpc/auth"
	"auth/core/adapter/handler/web/auth"
	"auth/core/adapter/handler/web/permissions"
	"auth/core/adapter/handler/web/subscriptions"

	"go.uber.org/fx"
)

var Module = fx.Options(
	auth.HandlerModule,
	permissions.HandlerModule,
	subscriptions.HandlerModule,
	authgrpc.HandlerModule,
)
