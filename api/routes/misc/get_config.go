package misc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/config"
)

// GetConfig is a route that returns the api configuration file.
// This can be useful for frontends to determine the list of valid boxes, services, and user lists.
func GetConfig(c *gin.Context) {
	cfg, err := config.Get()
	if err != nil {
		resp := gin.H{
			"message": "failed to get config",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, cfg)
}
