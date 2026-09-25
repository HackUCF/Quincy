package competition

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/gin-gonic/gin"
)

func PauseStatus(c *gin.Context) {

	db := conn.Get(c)

	record, err := pauses.GetPauseRecord(c.Request.Context(), db)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "failed to get pause status from db", "error": err},
		)
		return
	}

	c.JSON(http.StatusOK, record)
}
