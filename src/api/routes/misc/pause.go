package misc

import (
	"net/http"

	"github.com/HackUCF/quincy/api/services"
	"github.com/gin-gonic/gin"
)

// Pause handles quincy pause. This for things like lunch during a competition.
func Pause(c *gin.Context) {

	if services.IsPaused() {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already paused"})
		return
	}

	services.Pause()

	c.JSON(http.StatusOK, gin.H{"message": "quincy paused"})
}

// Unpause handles quincy unpause. This for things like lunch during a competition.
func Unpause(c *gin.Context) {

	if !services.IsPaused() {
		c.AbortWithStatusJSON(http.StatusTeapot, gin.H{"message": "quincy already unpaused"})
		return
	}

	services.Unpause()

	c.JSON(http.StatusOK, gin.H{"message": "quincy unpaused"})
}
