package adapter

import (
	"auth/core/adapter/handler"
	"auth/core/adapter/repository"

	"go.uber.org/fx"
)

var Module = fx.Options(
	handler.Module,
	repository.Module,
)
