package graphs

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/graphs"
	"github.com/gin-gonic/gin"
)

// GetHeatmap returns html and js containing a heatmap of historical uptime percentage per team per box/service.
func GetHeatmap(c *gin.Context) {
	data, err := graphs.GetHeatmapData()
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
