/*
Package users contains all the API routes involving scoring users.
This includes PCRs, as well as enumeration routes that can be used by frontends to create PCR pages.
*/
package users

import (
	"encoding/json"
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/users"
	"github.com/HackUCF/quincy/common/types"
	"github.com/gin-gonic/gin"
)

// PCR is the JSON object that must be posted to the server to complete a PCR.
type PCR struct {
	types.User
	UserListName types.UserListName `json:"user_list" example:"local users"`
	TeamNum      types.TeamNum      `json:"team_num"  example:"1"`
}

func userListExists(cfg *config.APIConfigSpec, name types.UserListName) bool {
	for _, ul := range cfg.UserLists {
		if name == ul.Name {
			return true
		}
	}
	return false
}

// SubmitPCR is a POST route for changing a scoring users password.
//
//	@Summary		Submit a password change request (PCR)
//	@Description	Updates the password for a scoring user in the specified userlist for the specified team.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			pcr	body		PCR		true	"Password change request"
//	@Success		200	{object}	object
//	@Failure		400	{object}	object
//	@Failure		500	{object}	object
//	@Router			/users [post]
func SubmitPCR(c *gin.Context) {
	db := conn.Get(c)
	cfg := config.Get(c)

	var pcr PCR
	err := json.NewDecoder(c.Request.Body).Decode(&pcr)
	if err != nil {
		resp := gin.H{
			"message": "failed to unmarshall json from request body",
			"error":   err,
			"pcr":     pcr,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !userListExists(cfg, pcr.UserListName) {
		resp := gin.H{
			"message":   "user list does not exist",
			"user_list": pcr.UserListName,
			"pcr":       pcr,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = users.UpdateUser(c.Request.Context(), db, pcr.UserListName, pcr.TeamNum, pcr.User)
	if err != nil {
		resp := gin.H{
			"message": "could not update userlist",
			"err":     err,
			"pcr":     pcr,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := gin.H{
		"message": "user updated",
		"pcr":     pcr,
	}
	c.JSON(http.StatusOK, resp)
}
