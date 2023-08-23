package auth

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (ctrl *Controller) Login(ctx *gin.Context) {
	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(authCtx)
	tracerProvider := span.TracerProvider()
	authCtx, span = tracerProvider.Tracer(ctrl.config.ServiceName).Start(authCtx, "adapter.auth.login")
	defer span.End()

	ctrl.log.Ctx(authCtx).Info("login")
	span.SetAttributes(attribute.String("adapter.name", "auth.Login"))

	var loginInput http.LoginInput
	if err := ctx.BindJSON(&loginInput); err != nil {
		ctrl.log.Ctx(authCtx).Error("error binding json", zap.Error(err))
		_ = ctx.Error(system.ErrInvalidInput)
		return
	}

	if err := ctrl.validate.Validate(authCtx, loginInput); err != nil {
		ctrl.log.Ctx(authCtx).Error("error validating input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	token, err := ctrl.service.Login(authCtx, loginInput.Email, loginInput.Password)
	if err != nil {
		ctrl.log.Ctx(authCtx).Error("error logging in", zap.Error(err))

		_ = ctx.Error(err)
		return
	}

	message := "user successfully logged in"
	res := system.NewHttpResponse(true, message, gin.H{"token": token}, 200)
	ctx.JSON(200, res)
}
