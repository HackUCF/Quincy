package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/db/scoring"
)

// GetServiceScores returns the final stats per box per service.
func GetServiceScores(c *gin.Context) {

	// pull up to date scores from the db
	scores, err := scoring.GetServiceScores()
	if err != nil {
		resp := gin.H{
			"message": "failed to get final scores per box per service",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, scores)
}
