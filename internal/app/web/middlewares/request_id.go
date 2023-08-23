package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func RequestIdMiddleware(log *otelzap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.NewUUID()
		if err != nil {
			ctx.Next()
		}

		// Set X-Request-Id header
		log.Debug("request id", zap.String("requestId", id.String()))

		ctx.Set("requestId", id.String())
		ctx.Writer.Header().Set("X-Request-Id", id.String())
		ctx.Next()
	}
}
