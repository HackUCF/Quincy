package users

import (
	"net/http"

	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/api/sinks/postgres/users"
	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all users from all userlists and their up to date passwords for all teams.
//
//	@Summary		Get all scoring users
//	@Description	Returns every user in every userlist for every team, including current passwords. Shape: {"team": {"userlist": [User]}}.
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	map[string]map[string][]types.User
//	@Failure		400	{object}	object
//	@Router			/users [get]
func GetAllUsers(c *gin.Context) {
	db := conn.Get(c)

	users, err := users.GetAllUsers(c.Request.Context(), db)
	if err != nil {
		resp := gin.H{
			"message": "failed to get all users",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, users)
}
