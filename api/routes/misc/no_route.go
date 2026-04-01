/*
Package misc contains routes that don't belong in other categories.
This includes error routes and API config routes.
*/
package misc

import (
	"net/http"
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
// To be fully accurate, this needs to be the last route registered on the server.
func NoRoute(router *gin.Engine) func(c *gin.Context) {

	// pull the routes from the router (at route definition time)
	routes := router.Routes()

	return func(c *gin.Context) {

		// construct a list of similar routes
		var similar []similarRoute

		// remove trailing slash if present
		reqPath := strings.TrimRight(c.Request.URL.Path, "/")

		// loop through all routes in the api
		for _, route := range routes {

			// do the same thing
			routePath := strings.TrimRight(route.Path, "/")

			// identify similar routes
			if strings.HasPrefix(routePath, reqPath) || strings.HasPrefix(reqPath, routePath) {
				similar = append(similar, similarRoute{
					Path:   routePath,
					Method: route.Method,
				})
			}
		}

		// return 404 with similar routes
		c.JSON(http.StatusNotFound, gin.H{
			"message":        "route not found or method not allowed",
			"path":           c.Request.URL.Path,
			"method":         c.Request.Method,
			"similar_routes": similar,
		})
	}
}
