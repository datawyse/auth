package middlewares

import (
	"auth/core/ports"
	"context"
	"errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"

	"auth/core/domain/system"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(log *otelzap.Logger, appValidator ports.AppValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// check if errors are validator.ValidationErrors
		for _, ginErr := range ctx.Errors {
			var validationErrors validator.ValidationErrors
			if errors.As(ginErr.Err, &validationErrors) {
				// extract the field and error message for each error
				validationErrors := appValidator.ValidationErrors(context.Background(), validationErrors)
				apiErrorResponse := system.NewHttpResponse(false, system.ErrValidationError.Message, validationErrors, system.ErrValidationError.Code)
				ctx.JSON(http.StatusBadRequest, apiErrorResponse)
				return
			}

			// check if error is of type ApiError
			var apiError *system.ApiError
			if errors.As(ginErr.Err, &apiError) {
				apiErrorResponse := system.NewHttpResponse(false, apiError.Message, nil, apiError.Code)
				ctx.JSON(apiError.Code, apiErrorResponse)
				return
			}

			apiErrorResponse := system.NewHttpResponse(false, ginErr.Error(), nil, 400)
			ctx.JSON(http.StatusBadRequest, apiErrorResponse)
			return
		}
	}
}
