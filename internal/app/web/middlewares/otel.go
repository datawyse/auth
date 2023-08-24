package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func GetIDMiddleware(log *otelzap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if span is valid
		if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
			TraceID := trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
			SpanID := trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()

			c.Writer.Header().Set("X-Trace-ID", TraceID)
			c.Writer.Header().Set("X-Span-ID", SpanID)

			log.Ctx(c.Request.Context()).Info("trace id: ", zap.String("trace_id", TraceID))
			log.Ctx(c.Request.Context()).Info("span id: ", zap.String("span_id", SpanID))
		}
	}
}
