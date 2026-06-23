package scoring

import (
	"net/http"

	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/scoring"
	"github.com/gin-gonic/gin"
)

// GetTeamScores returns the final stats per team.
//
//	@Summary		Get cumulative scores per team
//	@Description	Returns pass/fail check stats aggregated across all boxes and services, keyed by team number.
//	@Tags			scores
//	@Produce		json
//	@Success		200	{object}	map[string]types.ScoreResult
//	@Failure		400	{object}	object
//	@Failure		501	{object}	object
//	@Router			/scores/team [get]
func GetTeamScores(c *gin.Context) {
	db := conn.Get(c)

	// pull up to date scores from the db
	scores, err := scoring.GetTeamScores(c.Request.Context(), db)
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
