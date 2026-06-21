package scoring

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/scoring"
	"github.com/gin-gonic/gin"
)

// GetBoxScores returns the final stats per box.
//
//	@Summary		Get cumulative scores per box
//	@Description	Returns pass/fail check stats aggregated across all teams and services, keyed by box name.
//	@Tags			scores
//	@Produce		json
//	@Success		200	{object}	map[string]types.ScoreResult
//	@Failure		400	{object}	object
//	@Router			/scores/box [get]
func GetBoxScores(c *gin.Context) {
	db := conn.Get(c)

	// pull up to date scores from the db
	scores, err := scoring.GetBoxScores(c.Request.Context(), db)
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
