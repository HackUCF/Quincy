package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/db/users"
)

// GetAllUsers returns all users from all userlists and their up to date passwords for all teams.
func GetAllUsers(c *gin.Context) {
	users, err := users.GetAllUsers()
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
