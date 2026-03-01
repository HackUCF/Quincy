package graphs

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/graphs"
	"github.com/gin-gonic/gin"
)

// GetStandings returns html and js containing a bar chart of current total points per team.
func GetStandings(c *gin.Context) {
	data, err := graphs.GetStandingsData()
	if err != nil {
		resp := gin.H{
			"message": "failed to get standings data",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "standings.html", data)
	if err != nil {
		resp := gin.H{
			"message": "failed to render standings graph",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
}
