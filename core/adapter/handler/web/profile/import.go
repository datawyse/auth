package profile

import (
	"context"
	"fmt"
	"os"
	"time"

	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) ImportProfiles(ctx *gin.Context) {
	ctrl.log.Info("updating profile up")

	_, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.RequestTimeout)*time.Second)
	defer cancel()

	// read csv file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctrl.log.Error("error reading csv file", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// check if file is csv
	if file.Header.Get("Content-Type") != "text/csv" {
		ctrl.log.Error("error reading csv file", zap.Error(system.ErrInvalidInput))

		err = ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// save the file
	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctrl.log.Error("error saving csv file", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	// read csv files
	csvFile, err := os.Open(file.Filename)
	if err != nil {
		ctrl.log.Error("error opening csv file", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			ctrl.log.Error("error closing csv file", zap.Error(err))
		}
	}(csvFile)

	// create users
	// create profiles
	// create roles
	// create permissions

	message := fmt.Sprintf("import successful")
	ctx.JSON(200, system.NewHttpResponse(true, message, gin.H{"users": "users"}))
}
