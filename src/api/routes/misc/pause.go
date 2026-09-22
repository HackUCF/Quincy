package misc

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/services"
	"github.com/gin-gonic/gin"
)

// Pause handles quincy pause. This for things like lunch during a competition.
//
//	@Summary		Pause the competition
//	@Description	Halts scoring for planned interruptions such as a lunch break. While paused, agents requesting a check receive a no-op assignment instead, so no checks run and no scores are recorded. Returns 418 if the competition is already paused. The pause state is held in memory and does not survive an API server restart.
//	@Tags			misc
//	@Produce		json
//	@Success		200	{object}	object
//	@Failure		418	{object}	object
//	@Router			/pause [post]
func Pause(c *gin.Context) {

	if services.IsPaused() {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already paused"})
		return
	}

	services.Pause()

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

	if !services.IsPaused() {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already unpaused"})
		return
	}

	services.Unpause()

	c.JSON(http.StatusOK, gin.H{"message": "quincy unpaused"})
}
