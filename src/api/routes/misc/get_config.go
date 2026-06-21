package misc

import (
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/gin-gonic/gin"
)

// GetConfig is a route that returns the api configuration file.
// This can be useful for frontends to determine the list of valid boxes, services, and user lists.
//
//	@Summary		Get API configuration
//	@Description	Returns the loaded API config, including box definitions, user lists, and team count. Useful for building dynamic PCR forms.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	config.APIConfigSpec
//	@Router			/config [get]
func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Get(c))
}
