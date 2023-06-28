package subscriptions

import (
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl Controller) GetSubscriptions(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl/subscriptions")

	userId := ctx.MustGet("userId").(string)

	user, err := ctrl.userService.User(userId)
	if err != nil {
		ctrl.log.Error("error getting user", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	subscription, err := ctrl.service.FindSubscriptionByID(user.Subscription.String())
	if err != nil {
		ctrl.log.Error("error getting subscription", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	res := system.NewHttpResponse(true, "user subscription", gin.H{"subscription": subscription})
	ctx.JSON(200, res)
}
