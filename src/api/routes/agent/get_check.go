package agent

import (
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/services"
	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/gin-gonic/gin"
)

// GetCheck is a route that returns the next check to run.
// This is meant for interaction from the agents.
//
//	@Summary		Get next check to run
//	@Description	Returns the next fully-rendered service check for the agent to execute. Rotates round-robin across all services and teams.
//	@Tags			agent
//	@Produce		json
//	@Success		200	{object}	types.Service
//	@Failure		400	{object}	object
//	@Failure		500	{object}	object
//	@Router			/agent/new-check [get]
func GetCheck(c *gin.Context) {

	// grab globals from gin context
	cfg := config.Get(c)
	db, err := conn.GetE(c)

	// check for db errors, failed if the db sink is enabled
	// if this fails `db` is safely null and will be ignored by GetNext.
	if err != nil && cfg.Sinks.DBEnabled() {
		resp := gin.H{
			"message": "failed to get database connection from request context",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	check, err := services.GetNext(c.Request.Context(), cfg, db)
	if err != nil {
		resp := gin.H{
			"message": "failed to get check",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, check)
}
