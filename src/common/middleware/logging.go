/*
Package middleware implements some common gin middleware that can be shared between the API and a potential golang frontend.
Includes panic recovery and request logging. Request logs contain both a raw duration (for machine parsing) and a
human-readable pretty_duration field.
*/
package middleware

import (
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/common/log"
	"github.com/gin-gonic/gin"
)

// prettyDuration prints how long a request takes in human readable format.
// this is not good for exporting logs, but helpful for reading output.
func prettyDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		seconds := d.Seconds()
		return fmt.Sprintf("%.3fs", seconds)

	case d >= time.Millisecond:
		milis := float64(d.Microseconds()) / 1000.0
		return fmt.Sprintf("%.3fms", milis)

	default:
		micros := float64(d.Nanoseconds()) / 1000.0
		return fmt.Sprintf("%.3fµs", micros)
	}
}

// Logging is a middleware that tracks request information.
// It shows how long requests take, as well as some other information.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// record the start time
		start := time.Now()

		// run the route
		c.Next()

		// record the end time
		duration := time.Since(start)

		// log it!
		log.Debug(
			"request",
			"duration", duration,
			"pretty_duration", prettyDuration(duration),
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"agent", c.Request.UserAgent(),
			"client", c.ClientIP(),
			"path", c.Request.URL.Path,
		)
	}
}
