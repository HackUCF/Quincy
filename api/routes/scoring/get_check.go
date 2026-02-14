package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/services"
)

// GetCheck is a route that returns the next check to run.
// This is meant for interaction from the agents.
func GetCheck(c *gin.Context) {
	check, err := services.GetNext()
	if err != nil {
		resp := gin.H{
			"message": "failed to get check",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, check)
}
