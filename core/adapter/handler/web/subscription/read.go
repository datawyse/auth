package subscription

import (
	"context"
	"fmt"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl Controller) GetSubscriptions(ctx *gin.Context) {
	ctrl.log.Info("get subscriptions")

	subsCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	userId := ctx.MustGet("userId").(string)

	user, err := ctrl.userService.User(subsCtx, userId)
	if err != nil {
		ctrl.log.Error("error getting user", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// read user subscription from attributes
	var subscriptionId string
	for name, value := range *user.Attributes {
		if name == "subscription" {
			// extract from list of subscriptions
			subscriptionId = value[0]
			break
		}
	}

	if subscriptionId == "" {
		err = fmt.Errorf("no subscription found for user")
		ctrl.log.Error("error getting subscription", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	subscription, err := ctrl.service.FindSubscriptionByID(subsCtx, subscriptionId)
	if err != nil {
		ctrl.log.Error("error getting subscription", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	res := system.NewHttpResponse(true, "user subscription", gin.H{"subscription": subscription}, 200)
	ctx.JSON(200, res)
}
