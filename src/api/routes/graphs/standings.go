package graphs

import (
	"net/http"

	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/graphs"
	"github.com/gin-gonic/gin"
)

// GetStandings returns html and js containing a bar chart of current total points per team.
//
//	@Summary		Standings bar chart
//	@Description	Returns an embeddable HTML/JS fragment showing current total points per team as a bar chart.
//	@Tags			graphs
//	@Produce		html
//	@Success		200		{string}	string
//	@Failure		500		{object}	object
//	@Failure		501		{object}	object
//	@Router			/graphs/standings [get]
func GetStandings(c *gin.Context) {
	db := conn.Get(c)

	data, err := graphs.GetStandingsData(c.Request.Context(), db)
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
