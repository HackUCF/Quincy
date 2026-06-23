package config

import "github.com/HackUCF/quincy/common/types"

// APIConfigSpec defines the entire API configuration file.
type APIConfigSpec struct {
	NumTeams  types.TeamNum  `yaml:"num_teams"  mapstructure:"num_teams"  json:"num_teams"  example:"5"`
	Sinks     Sinks          `yaml:"sinks"      mapstructure:"sinks"      json:"sinks"`
	Boxes     []BoxSpec      `yaml:"boxes"      mapstructure:"boxes"      json:"boxes"`
	UserLists []UserListSpec `yaml:"user_lists" mapstructure:"user_lists" json:"user_lists"`
	HTTP      HTTPSpec       `yaml:"http"       mapstructure:"http"       json:"http"`
}

// BoxSpec contains the config specification of one box.
// It is the parent to a list of services.
// The Name is enforced to be unique, but multiple boxes can have identical services or hosts.
type BoxSpec struct {
	Name     types.BoxName       `yaml:"name"     mapstructure:"name"     json:"name"     example:"scrapyard"`
	Host     string              `yaml:"host"     mapstructure:"host"     json:"host"     example:"127.0.0.{}"`
	Services []types.ServiceSpec `yaml:"services" mapstructure:"services" json:"services"`
}

// UserListSpec contains the definition of a userlist.
// Userlists have a unique name.
// Optionally, they can have a Domain or NetBIOS name.
// These apply to every user in the list.
type UserListSpec struct {
	Name        types.UserListName `yaml:"name"    mapstructure:"name"    json:"name"               example:"local users"`
	DomainName  string             `yaml:"domain"  mapstructure:"domain"  json:"domain,omitempty"   example:"quin.cy"`
	NetBIOSName string             `yaml:"netbios" mapstructure:"netbios" json:"netbios,omitempty"  example:"QUIN"`
	Users       []UserSpec         `yaml:"users"   mapstructure:"users"   json:"users"`
}

// UserSpec is the config specification of a single user.
type UserSpec struct {
	Username string `yaml:"username" mapstructure:"username" json:"username" example:"geraldo"`
	Password string `yaml:"password" mapstructure:"password" json:"password" example:"BuyMyNFT1!"`
}

// HTTPSpec defines the settings applied to the API's listener.
type HTTPSpec struct {
	Host string `yaml:"host" mapstructure:"host" json:"host" example:"127.0.0.1"`
	Port int    `yaml:"port" mapstructure:"port" json:"port" example:"8888"`
}
