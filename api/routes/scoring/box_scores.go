package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/db/scoring"
)

// GetBoxScores returns the final stats per box.
func GetBoxScores(c *gin.Context) {

	// pull up to date scores from the db
	scores, err := scoring.GetBoxScores()
	if err != nil {
		resp := gin.H{
			"message": "failed to get final scores per box",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, scores)
}
