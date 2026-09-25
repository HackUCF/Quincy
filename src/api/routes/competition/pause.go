package competition

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/gin-gonic/gin"
)

func ChangePauseState(desiredState types.PauseState, pauseType types.PauseType) func(c *gin.Context) {

	return func(c *gin.Context) {
		db := conn.Get(c)

		changed, err := pauses.ChangePauseState(c.Request.Context(), db, desiredState, pauseType)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message":       "failed to change pause state",
					"error":         err,
					"desired_state": desiredState,
					"pause_type":    pauseType,
				},
			)
			return
		}

		if !changed {
			c.AbortWithStatusJSON(
				http.StatusTeapot,
				gin.H{
					"message":       "pause already in desired state",
					"desired_state": desiredState,
					"pause_type":    pauseType,
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "quincy pause state changed"},
		)
	}
}
