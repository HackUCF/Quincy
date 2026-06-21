package graphs

import (
	"net/http"

	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/graphs"
	"github.com/gin-gonic/gin"
)

// GetHeatmap returns html and js containing a heatmap of historical uptime percentage per team per box/service.
//
//	@Summary		Uptime heatmap
//	@Description	Returns an embeddable HTML/JS fragment showing historical uptime percentage as a heatmap.
//	@Tags			graphs
//	@Produce		html
//	@Success		200		{string}	string
//	@Failure		500		{object}	object
//	@Router			/graphs/heatmap [get]
func GetHeatmap(c *gin.Context) {
	db := conn.Get(c)

	data, err := graphs.GetHeatmapData(c.Request.Context(), db)
	if err != nil {
		resp := gin.H{
			"message": "failed to get heatmap data",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "heatmap.html", data)
	if err != nil {
		resp := gin.H{
			"message": "failed to render heatmap",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
}
