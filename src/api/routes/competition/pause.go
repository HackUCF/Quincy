package competition

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/gin-gonic/gin"
)

// ChangePauseState builds a handler that moves one pause kind into one state.
// The state and pause kind are bound at route registration time, so each of the
// four pause/unpause paths is this same handler with different arguments.
//
//	@Summary		Pause or unpause the competition
//	@Description	Moves one of the two independent pauses into the requested state. A scoring pause leaves checks running but stops results counting toward team totals; a check pause stops checks being handed to agents entirely. The target state and pause kind are fixed by the path, so there is no request body. Requesting a state the competition is already in changes nothing and returns 418.
//	@Tags			comp
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Failure		500	{object}	object
//	@Failure		501	{object}	object
//	@Router			/comp/pause-checks [post]
//	@Router			/comp/unpause-checks [post]
//	@Router			/comp/pause-scoring [post]
//	@Router			/comp/unpause-scoring [post]
func ChangePauseState(desiredState types.PauseState, pauseType types.PauseType) func(c *gin.Context) {

	return func(c *gin.Context) {
		db := conn.Get(c)

		changed, err := pauses.ChangePauseState(c.Request.Context(), db, desiredState, pauseType)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message":       "failed to change pause state",
					"error":         err.Error(),
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
