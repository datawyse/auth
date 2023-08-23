package auth

import (
	"context"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) Me(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl/me")

	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	userId := ctx.MustGet("userId").(string)
	userProfile, err := ctrl.service.User(authCtx, userId)
	if err != nil {
		ctrl.log.Error("error getting user profile", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := "user profile details"
	res := system.NewHttpResponse(true, message, gin.H{"user": userProfile}, 200)
	ctx.JSON(200, res)
}
