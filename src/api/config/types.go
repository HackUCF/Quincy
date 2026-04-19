package config

import "github.com/HackUCF/quincy/common/types"

// APIConfigSpec defines the entire API configuration file.
type APIConfigSpec struct {
	NumTeams  types.TeamNum  `yaml:"num_teams"  mapstructure:"num_teams"  json:"num_teams"`
	DBFile    string         `yaml:"db_file"    mapstructure:"db_file"    json:"db_file"`
	Boxes     []BoxSpec      `yaml:"boxes"      mapstructure:"boxes"      json:"boxes"`
	UserLists []UserListSpec `yaml:"user_lists" mapstructure:"user_lists" json:"user_lists"`
	HTTP      HTTPSpec       `yaml:"http"       mapstructure:"http"       json:"http"`
}

// BoxSpec contains the config specification of one box.
// It is the parent to a list of services.
// The Name is enforced to be unique, but multiple boxes can have identical services or hosts.
type BoxSpec struct {
	Name     types.BoxName       `yaml:"name"     mapstructure:"name"     json:"name"`
	Host     string              `yaml:"host"     mapstructure:"host"     json:"host"`
	Services []types.ServiceSpec `yaml:"services" mapstructure:"services" json:"services"`
}

// UserListSpec contains the definition of a userlist.
// Userlists have a unique name.
// Optionally, they can have a Domain or NetBIOS name.
// These apply to every user in the list.
type UserListSpec struct {
	Name        types.UserListName `yaml:"name"    mapstructure:"name"    json:"name"`
	DomainName  string             `yaml:"domain"  mapstructure:"domain"  json:"domain,omitempty"`
	NetBIOSName string             `yaml:"netbios" mapstructure:"netbios" json:"netbios,omitempty"`
	Users       []UserSpec         `yaml:"users"   mapstructure:"users"   json:"users"`
}

// UserSpec is the config specification of a single user.
type UserSpec struct {
	Username string `yaml:"username" mapstructure:"username" json:"username"`
	Password string `yaml:"password" mapstructure:"password" json:"password"`
}

// HTTPSpec defines the settings applied to the API's listener.
type HTTPSpec struct {
	Host string `yaml:"host" mapstructure:"host" json:"host"`
	Port int    `yaml:"port" mapstructure:"port" json:"port"`
}
