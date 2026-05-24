package graphs

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/graphs"
	"github.com/gin-gonic/gin"
)

// GetScoreboard returns html and js containing a simple scoreaboard.
func GetScoreboard(c *gin.Context) {
	db := conn.Get(c)

	data, err := graphs.GetScoreboardData(c.Request.Context(), db)
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
