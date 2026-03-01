package graphs

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/graphs"
	"github.com/gin-gonic/gin"
)

// GetScoreboard returns html and js containing a simple scoreaboard.
func GetScoreboard(c *gin.Context) {
	data, err := graphs.GetScoreboardData()
	if err != nil {
		resp := gin.H{
			"message": "failed to get scoreboard data",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "scoreboard.html", data)
	if err != nil {
		resp := gin.H{
			"message": "failed to render scoreboard",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
}
