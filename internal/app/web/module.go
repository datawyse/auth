package web

import (
	"auth/internal/app/web/server"

	"go.uber.org/fx"
)

var Module = fx.Options(server.ServiceAppModule, server.EngineModule, server.GinRouterModule, server.MiddlewareModule)
