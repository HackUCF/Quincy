/*
Package users contains all the API routes involving scoring users.
This includes PCRs, as well as enumeration routes that can be used by frontends to create PCR pages.
*/
package users

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db/users"
	"github.com/HackUCF/Quincy/common/types"
)

// PCR is the JSON object that must be posted to the server to complete a PCR.
type PCR struct {
	types.User
	UserListID types.UserListID `json:"user_list"`
	TeamNum    types.TeamNum    `json:"team_num"`
}

// SubmitPCR is a POST route for changing a scoring users password.
// Takes the following input:
func SubmitPCR(c *gin.Context) {
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

	if !config.UserListExists(pcr.UserListID) {
		resp := gin.H{
			"message":   "user list does not exist",
			"user_list": pcr.UserListID,
			"pcr":       pcr,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = users.UpdateUser(pcr.UserListID, pcr.TeamNum, pcr.User)
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
