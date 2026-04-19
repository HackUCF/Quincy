package users

import (
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
	"github.com/gin-gonic/gin"
)

// UserListInfo represents a description of a userlist.
// This can be used by a frontend to create a PCR form that does not exposes usernames/passwords.
type UserListInfo struct {
	Name    types.UserListName `json:"name"`
	Domain  string             `json:"domain"`
	NetBIOS string             `json:"netbios"`
}

// GetUserLists is a route that returns descriptions of all userlists.
// Notably it does not actually return any usernames or passwords.
// This is intended for use to generate an obfuscated PCR form (like CCDC).
func GetUserLists(c *gin.Context) {
	cfg := config.Get(c)

	allInfo := make([]UserListInfo, 0, len(cfg.UserLists))

	for _, ul := range cfg.UserLists {
		info := UserListInfo{
			Name:    ul.Name,
			Domain:  ul.DomainName,
			NetBIOS: ul.NetBIOSName,
		}
		allInfo = append(allInfo, info)
	}

	c.JSON(http.StatusOK, allInfo)
}
