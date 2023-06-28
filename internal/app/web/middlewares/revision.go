package middlewares

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RevisionMiddleware(log *zap.Logger) gin.HandlerFunc {
	// Revision file contents will be only loaded once per process
	data, err := ioutil.ReadFile("REVISION")

	// If we cant read file, just skip to the next request handler
	// This is pretty much a NOOP middleware :)
	if err != nil {
		// Make sure to log error so it could be spotted
		log.Error("failed to read revision file", zap.Error(err))

		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Clean up the value since it could contain line breaks
	revision := strings.TrimSpace(string(data))

	// Set out header value for each response
	return func(c *gin.Context) {
		// only GET requests are cached
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Set X-Revision header
		log.Debug("revision", zap.String("revision", revision))

		c.Writer.Header().Set("X-Revision", revision)
		c.Next()
	}
}
