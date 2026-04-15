package scoring

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/scoring"
	"github.com/gin-gonic/gin"
)

// GetServiceScores returns the final stats per box per service.
func GetServiceScores(c *gin.Context) {
	db := conn.Get(c)

	// pull up to date scores from the db
	scores, err := scoring.GetServiceScores(db)
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
