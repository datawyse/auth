package profile

import (
	"context"
	"fmt"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) GetProfile(ctx *gin.Context) {
	ctrl.log.Info("updating profile up")

	profileCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.RequestTimeout)*time.Second)
	defer cancel()

	id := ctx.Param("id")
	if id == "" {
		ctrl.log.Error("error binding json", zap.Error(system.ErrInvalidInput))

		err := ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
			return
		}
	}

	user, err := ctrl.userService.User(profileCtx, id)
	if err != nil {
		ctrl.log.Error("error signing up", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := fmt.Sprintf("user details")
	ctx.JSON(200, system.NewHttpResponse(true, message, gin.H{"user": user}, 200))
}
