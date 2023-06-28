package core

import (
	"auth/core/adapter"
	"auth/core/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	service.Module,
	adapter.Module,
)
