package graphs

import (
	"net/http"

	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/graphs"
	"github.com/HackUCF/quincy/common/log"
	"github.com/gin-gonic/gin"
)

// GetScores returns an HTML chart of scores over time.
//
//	@Summary		Scores over time chart
//	@Description	Returns an embeddable HTML/JS fragment rendering a chart of score trends over time.
//	@Tags			graphs
//	@Produce		html
//	@Success		200		{string}	string
//	@Failure		500		{object}	object
//	@Failure		501		{object}	object
//	@Router			/graphs/scores [get]
func GetScores(c *gin.Context) {
	db := conn.Get(c)

	data, err := graphs.GetScoresData(c.Request.Context(), db)
	if err != nil {
		resp := gin.H{
			"message": "failed to get scores data",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		log.Panic(err.Error())
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "scores.html", data)
	if err != nil {
		resp := gin.H{
			"message": "failed to render scores graph",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		log.Panic(err.Error())
		return
	}
}
