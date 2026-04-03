package graphs

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/api/db/graphs"
	"github.com/HackUCF/Quincy/common/log"
	"github.com/gin-gonic/gin"
)

func GetScores(c *gin.Context) {
	db := conn.Get(c)

	data, err := graphs.GetScoresData(db)
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
