package health

import (
	"context"
	"errors"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReadHealth - web status
func (ctrl *Controller) ReadHealth(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl")

	healthCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	healthRes, err := ctrl.service.ReadHealth(healthCtx)
	if err != nil {
		err = ctx.Error(system.ErrServer)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// READ "X-Revision" HEADER
	xRevision, ok := ctx.Get("revision")
	if !ok {
		err = ctx.Error(errors.New("X-Revision header is empty"))
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// update revision ID
	healthRes.Revision = xRevision.(string)

	message := "health status"
	res := system.NewHttpResponse(true, message, gin.H{
		"health": healthRes,
	})
	ctx.JSON(200, res)
}
