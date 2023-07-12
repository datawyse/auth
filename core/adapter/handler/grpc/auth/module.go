package auth

import "go.uber.org/fx"

var HandlerModule = fx.Options(
	fx.Provide(NewAuthGRPCService),
	fx.Provide(NewSubscriptionGRPCService),
)
