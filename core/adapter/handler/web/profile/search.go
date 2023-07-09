package profile

import (
	"context"
	"fmt"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) SearchProfile(ctx *gin.Context) {
	ctrl.log.Info("updating profile up")

	profileCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.RequestTimeout)*time.Second)
	defer cancel()

	// extract query params for username and email
	username := ctx.Query("username")
	email := ctx.Query("email")

	// check if both are empty
	if username == "" && email == "" {
		ctrl.log.Error("error binding json", zap.Error(system.ErrInvalidInput))

		err := ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
			return
		}
	}

	// check which exist first and then call the appropriate service
	if username != "" {
		user, err := ctrl.userService.UserByUsername(profileCtx, username)
		if err != nil {
			ctrl.log.Error("error signing up", zap.Error(err))

			err = ctx.Error(err)
			if err != nil {
				ctrl.log.Error("error aborting with error", zap.Error(err))
			}
			return
		}

		message := fmt.Sprintf("user details")
		ctx.JSON(200, system.NewHttpResponse(true, message, gin.H{"user": user}))
	} else {
		user, err := ctrl.userService.UserByEmail(profileCtx, email)
		if err != nil {
			ctrl.log.Error("error signing up", zap.Error(err))

			err = ctx.Error(err)
			if err != nil {
				ctrl.log.Error("error aborting with error", zap.Error(err))
			}
			return
		}

		message := fmt.Sprintf("user details")
		ctx.JSON(200, system.NewHttpResponse(true, message, gin.H{"user": user}))
	}
}
