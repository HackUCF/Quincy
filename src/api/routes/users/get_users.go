package users

import (
	"net/http"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/api/db/users"
	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all users from all userlists and their up to date passwords for all teams.
func GetAllUsers(c *gin.Context) {
	db := conn.Get(c)

	users, err := users.GetAllUsers(db)
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
