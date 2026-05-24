package users

import (
	"net/http"

	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/users"
	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all users from all userlists and their up to date passwords for all teams.
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
