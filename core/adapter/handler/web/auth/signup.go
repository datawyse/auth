package auth

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"time"

	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (ctrl *Controller) Signup(ctx *gin.Context) {
	ctrl.log.Info("signup")

	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	tracer := otel.GetTracerProvider().Tracer(ctrl.config.ServiceName)
	authCtx, span := tracer.Start(authCtx, "auth.adapter.signup")
	defer span.End()

	// get the request id from the header
	requestId, ok := ctx.Get("requestId")
	if !ok {
		// raise error and abort
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": "missing request id",
		})
		return
	}

	//span := trace.SpanFromContext(ctx.Request.Context())
	//span.SetAttributes(attribute.String("controller.name", "auth.signup"))

	var signupInput http.SignupInput
	if err := ctx.ShouldBindJSON(&signupInput); err != nil {
		ctrl.log.Error("error binding json", zap.Error(err))

		err = ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Ctx(authCtx).Error("error aborting with error", zap.Error(err))
		}
		return
	}

	if err := ctrl.validate.Validate(authCtx, signupInput); err != nil {
		ctrl.log.Ctx(authCtx).Error("error validating input", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Ctx(authCtx).Error("error aborting with error", zap.Error(err))
		}
		return
	}

	userId, err := ctrl.service.Signup(authCtx, &signupInput)
	if err != nil {
		ctrl.log.Ctx(authCtx).Error("error signing up", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Ctx(authCtx).Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// Pass log data using log events and include the requestId
	span.AddEvent("ControllerLogicCompleted", trace.WithAttributes(
		attribute.String("message", "Controller logic completed successfully"),
		attribute.String("requestId", requestId.(string)),
	))

	message := fmt.Sprintf("user successfully created")
	res := system.NewHttpResponse(true, message, gin.H{"userId": userId}, 201)
	ctx.JSON(201, res)
}
