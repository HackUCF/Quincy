package misc

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/gin-gonic/gin"
)

// Pause handles quincy pause. This for things like lunch during a competition.
//
//	@Summary		Pause the competition
//	@Description	Halts scoring for planned interruptions such as a lunch break. While paused, agents requesting a check receive a no-op assignment instead, so no checks run and no scores are recorded. Returns 418 if the competition is already paused. The pause state is held in memory: an API server restart discards it and falls back to the start_paused setting in the config.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Router			/pause [post]
func Pause(c *gin.Context) {
	db := conn.Get(c)

	changed, err := pauses.ChangePauseState(c.Request.Context(), db, types.Paused)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "failed to pause competition", "error": err},
		)
		return
	}

	if !changed {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already paused"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quincy paused"})
}

// Unpause handles quincy unpause. This for things like lunch during a competition.
//
//	@Summary		Resume the competition
//	@Description	Resumes scoring after a pause. Agents begin receiving real checks again from wherever the queue left off. Returns 418 if the competition is not currently paused.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Router			/unpause [post]
func Unpause(c *gin.Context) {

	db := conn.Get(c)

	changed, err := pauses.ChangePauseState(c.Request.Context(), db, types.Unpaused)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "failed to unpause competition", "error": err},
		)
		return
	}

	if !changed {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already unpaused"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quincy unpaused"})
}

// PauseStatus handles reporting whether quincy is currently paused.
//
//	@Summary		Get the competition pause status
//	@Description	Reports whether the competition is currently paused and when it entered that state. Returns an object with two fields: `is_paused` (bool) and `since` (RFC 3339 timestamp). The timestamp is the moment of the last pause or unpause, or the time the API server started if the state has not changed since boot.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Router			/pause-status [get]
func PauseStatus(c *gin.Context) {

	db := conn.Get(c)

	record, err := pauses.GetPauseRecord(c.Request.Context(), db)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "failed to get pause status from db", "error": err},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_paused": record.State, "since": record.Since})
}
