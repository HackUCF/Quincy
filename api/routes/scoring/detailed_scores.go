package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/db/scoring"
)

// GetDetailedScores returns the final stats per team per box per service.
// This is as detailed as final scores can be.
func GetDetailedScores(c *gin.Context) {

	// pull up to date scores from the db
	scores, err := scoring.GetDetailedScores()
	if err != nil {
		resp := gin.H{
			"message": "failed to get final scores",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, scores)
}
