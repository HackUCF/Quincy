package misc

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/gin-gonic/gin"
)

// GetConfig is a route that returns the api configuration file.
// This can be useful for frontends to determine the list of valid boxes, services, and user lists.
func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Get())
}
