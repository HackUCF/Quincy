package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/db/scoring"
)

// GetTeamScores returns the final stats per team.
func GetTeamScores(c *gin.Context) {

	// pull up to date scores from the db
	scores, err := scoring.GetTeamScores()
	if err != nil {
		resp := gin.H{
			"message": "failed to get final scores per team",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, scores)
}
