/*
Package misc contains routes that don't belong in other categories.
This includes error routes and API config routes.
*/
package misc

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type similarRoute struct {
	Path   string `json:"path"`
	Method string `json:"methods"`
}

// NoRoute responds with a 404 when no other route is matched.
// Provides the user a list of "similar" endpoints available.
// This is helpful while manually testing.
func NoRoute(router *gin.Engine) func(c *gin.Context) {
	routes := router.Routes()

	return func(c *gin.Context) {
		// get the path and parent for the request
		reqPath := c.Request.URL.Path
		reqParent := filepath.Dir(reqPath)

		// accumulate similar routes
		similar := make([]similarRoute, 0)

		for _, route := range routes {
			routePath := strings.TrimRight(route.Path, "/")

			// check if the requested route has the same parent as an existing one
			if strings.HasPrefix(routePath, reqParent) {
				// add it to a list
				newRoute := similarRoute{
					Path:   routePath,
					Method: route.Method,
				}
				similar = append(similar, newRoute)
			}
		}

		// add all routes if nothing available
		if len(similar) == 0 {
			for _, route := range routes {
				newRoute := similarRoute{
					Path:   route.Path,
					Method: route.Method,
				}
				similar = append(similar, newRoute)
			}
		}

		resp := gin.H{
			"message":        "route not found or method not allowed",
			"path":           c.Request.URL.Path,
			"method":         c.Request.Method,
			"similar_routes": similar, // send to user
		}
		c.JSON(http.StatusNotFound, resp)
	}
}
