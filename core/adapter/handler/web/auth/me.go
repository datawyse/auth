package auth

import (
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) Me(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl/me")

	userId := ctx.MustGet("userId").(string)
	userProfile, err := ctrl.service.User(userId)
	if err != nil {
		ctrl.log.Error("error getting user profile", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := "user profile details"
	res := system.NewHttpResponse(true, message, gin.H{
		"user": userProfile,
	})
	ctx.JSON(200, res)
}
