package competition

import (
	"net/http"

	"github.com/HackUCF/Quincy/src/api/sinks/postgres/conn"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/gin-gonic/gin"
)

// PauseStatus reports the current state of both competition pauses.
//
//	@Summary		Get competition pause status
//	@Description	Reports both pause kinds at once — whether scoring is paused and whether checks are paused — each with the moment it entered its current state.
//	@Tags			comp
//	@Produce		json
//	@Success		200	{object}	types.PauseRecord
//	@Failure		500	{object}	object
//	@Failure		501	{object}	object
//	@Router			/comp/pause-status [get]
func PauseStatus(c *gin.Context) {

	db := conn.Get(c)

	record, err := pauses.GetPauseRecord(c.Request.Context(), db)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "failed to get pause status from db", "error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, record)
}
