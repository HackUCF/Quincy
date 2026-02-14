/*
Package misc contains routes that don't belong in other categories.
This includes error routes and API config routes.
*/
package misc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NoRoute responds with a 404 when no other route is matched.
func NoRoute(c *gin.Context) {
	resp := gin.H{
		"message": "route not found or method not allowed",
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	}
	c.JSON(http.StatusNotFound, resp)
}
