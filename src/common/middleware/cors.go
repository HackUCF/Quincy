package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InsecureCORS is a middleware that responds to option headers allowing all origins, methods, and headers.
// This should be replaces with a secure set up for production.
func InsecureCORS() gin.HandlerFunc {
	// just make it work
	return cors.New(cors.Config{
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
		AllowWebSockets: true,
	})
}
