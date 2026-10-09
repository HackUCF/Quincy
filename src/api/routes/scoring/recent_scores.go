package scoring

import (
	"net/http"
	"strconv"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/scoring"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/gin-gonic/gin"
)

// maxRecentChecks caps how much history one request can pull per service.
const maxRecentChecks = 100

// defaultRecentChecks is used when the caller does not ask for a specific count.
const defaultRecentChecks = 10

// GetRecentNChecks returns the N most recent checks for every service of one team.
//
//	@Summary		Get the N most recent checks for a team
//	@Description	Returns the most recent checks for every service on every box of a single team. Sorted by box, then service, then newest check first.
//	@Tags			scores
//	@Produce		json
//	@Param			team	path		int	true	"Team number"
//	@Param			n		query		int	false	"Checks to return per service (default 10, max 100)"
//	@Success		200		{array}		types.Score
//	@Failure		400		{object}	object
//	@Failure		501		{object}	object
//	@Router			/scores/recent/{team} [get]
func GetRecentNChecks(c *gin.Context) {
	db := conn.Get(c)
	cfg := config.Get(c)

	// team number comes from the path and must be a real team
	teamRaw, err := strconv.ParseUint(c.Param("team"), 10, 32)
	if err != nil {
		resp := gin.H{
			"message": "invalid team number",
			"error":   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// verify in range
	teamNum := types.TeamNum(teamRaw)
	if teamNum < 1 || teamNum > cfg.NumTeams {

		resp := gin.H{
			"message": "team number out of range",
			"error":   "team must be between 1 and " + strconv.FormatUint(uint64(cfg.NumTeams), 10),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// optionally allow the user to specify a number of checks
	numChecks := defaultRecentChecks

	if raw := c.Query("n"); raw != "" {

		// convert to int
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			resp := gin.H{
				"message": "invalid check count",
				"error":   err.Error(),
			}
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		// ensure it's sane
		if parsed < 1 || parsed > maxRecentChecks {
			resp := gin.H{
				"message": "check count out of range",
				"error":   "n must be between 1 and " + strconv.Itoa(maxRecentChecks),
			}
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		numChecks = parsed
	}

	// grab from by
	scores, err := scoring.GetRecentNScores(c.Request.Context(), db, cfg, teamNum, numChecks)
	if err != nil {

		resp := gin.H{
			"message": "could not get recent checks",
			"error":   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, scores)
}
