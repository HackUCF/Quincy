/*
Package scoring contains all the routes involving scorechecks.
This includes agent routes to get new checks and post completed ones,
as well as frontend routes to get final scores or the current status.
*/
package scoring

import (
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/scoring"
	"github.com/gin-gonic/gin"
)

// GetRecentChecks gets a list of the the most recent checks for every service.
// It is sorted by team number, then box ID, then service ID
//
//	@Summary		Get most recent check result per service per team
//	@Description	Returns the latest pass/fail score for every (team, box, service) combination. Sorted by team, box, service.
//	@Tags			scores
//	@Produce		json
//	@Success		200	{array}		types.Score
//	@Failure		400	{object}	object
//	@Failure		501	{object}	object
//	@Router			/scores/current [get]
func GetRecentChecks(c *gin.Context) {
	db := conn.Get(c)
	cfg := config.Get(c)

	status, err := scoring.GetCurrentServiceStatus(c.Request.Context(), db, cfg)
	if err != nil {
		resp := gin.H{
			"message": "could not get current service status",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, status)
}
