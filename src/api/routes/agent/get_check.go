package agent

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/services"
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
//	@Router			/agent/new-check [get]
func GetCheck(c *gin.Context) {
	db := conn.Get(c)

	check, err := services.GetNext(c.Request.Context(), db)
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
