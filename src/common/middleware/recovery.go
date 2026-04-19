package middleware

import (
	"runtime/debug"

	"github.com/HackUCF/quincy/common/log"
	"github.com/gin-gonic/gin"
)

// Recovery is the first middleware in the chain.
// It recovers from panics that happen in routes.
// There is a default one but it has ugly logging.
// Without this middleware, every panic in a route would crash the entire server.
func Recovery(includeStackTrace bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// recover if a panic has happened
		// defer to run after function exit
		defer func() {
			if err := recover(); err != nil {

				// optional stack trace
				var stacktrace string
				if includeStackTrace {
					stacktrace = string(debug.Stack())
				}

				args := []any{
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"ip", c.ClientIP(),
				}

				if includeStackTrace {
					args = append(args, "stacktrace", stacktrace)
				}

				log.Error("panic recovered", args...)

				c.AbortWithStatus(500)
			}
		}()

		// call every other handler
		// this includes the route and other middleware
		c.Next()
	}
}
