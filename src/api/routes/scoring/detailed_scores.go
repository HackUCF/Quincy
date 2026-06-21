package scoring

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/scoring"
	"github.com/gin-gonic/gin"
)

// GetDetailedScores returns the final stats per team per box per service.
// This is as detailed as final scores can be.
//
//	@Summary		Get fully detailed scores
//	@Description	Returns pass/fail stats keyed by team, then box, then service. Shape: {"team": {"box": {"service": ScoreResult}}}.
//	@Tags			scores
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		400	{object}	object
//	@Router			/scores/detailed [get]
func GetDetailedScores(c *gin.Context) {
	db := conn.Get(c)

	// pull up to date scores from the db
	scores, err := scoring.GetDetailedScores(c.Request.Context(), db)
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
