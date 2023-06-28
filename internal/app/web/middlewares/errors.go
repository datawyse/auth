package middlewares

import (
	"net/http"

	"auth/core/domain/system"
	"auth/core/service/app-validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.Logger, appValidator *app_validator.AppValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// check if errors are validator.ValidationErrors
		for _, ginErr := range ctx.Errors {
			log.Debug("api error")

			if err, ok := ginErr.Err.(validator.ValidationErrors); ok {
				// extract the field and error message for each error
				validationErrors := appValidator.ValidationErrors(err)
				apiError := system.NewHttpResponse(false, "input validation failed", validationErrors)
				ctx.JSON(http.StatusBadRequest, apiError)
				return // exit on first error

			}

			// check if error is of type ApiError
			if err, ok := ginErr.Err.(*system.ApiError); ok {
				apiError := system.NewHttpResponse(false, err.Message, nil)
				ctx.JSON(err.Code, apiError)
				return // exit on first error
			}

			apiError := system.NewHttpResponse(false, ginErr.Error(), nil)
			ctx.JSON(http.StatusBadRequest, apiError)
			return // exit on first error
		}
	}
}
