/*
Package scoring contains all the routes involving scorechecks.
This includes agent routes to get new checks and post completed ones,
as well as frontend routes to get final scores or the current status.
*/
package scoring

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/api/db/scoring"
	"github.com/gin-gonic/gin"
)

// GetRecentChecks gets a list of the the most recent checks for every service.
// It is sorted by team number, then box ID, then service ID
func GetRecentChecks(c *gin.Context) {
	db := conn.Get(c)

	status, err := scoring.GetCurrentServiceStatus(db)
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
