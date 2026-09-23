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
//	@Description	Halts scoring for planned interruptions such as a lunch break. Agents keep receiving checks and keep running them while paused; results are still archived and still reported as the current status, but they are not added to a team's pass and total counters, so no team loses uptime. Returns 418 if the competition is already paused. The pause state is stored in the database and survives an API server restart, so this endpoint requires the PostgreSQL sink and returns 501 without it.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Failure		500	{object}	object
//	@Failure		501	{object}	object
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
//	@Description	Resumes scoring after a pause. Check results begin counting toward team pass and total counters again from the moment of the call. Returns 418 if the competition is not currently paused. Requires the PostgreSQL sink; returns 501 without it.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Failure		500	{object}	object
//	@Failure		501	{object}	object
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
//	@Description	Reports whether the competition is currently paused and when it entered that state. Returns an object with two fields: `is_paused` (bool) and `since` (RFC 3339 timestamp). The timestamp is the moment of the last recorded pause or unpause, read from the database. Requires the PostgreSQL sink; returns 501 without it.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		500	{object}	object
//	@Failure		501	{object}	object
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
