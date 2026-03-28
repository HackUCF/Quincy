package config

import "github.com/HackUCF/Quincy/common/types"

// APIConfigSpec defines the entire API configuration file.
type APIConfigSpec struct {
	NumTeams  types.TeamNum  `yaml:"num_teams,required"  json:"num_teams"`
	DBFile    string         `yaml:"db_file,required"    json:"db_file"`
	Boxes     []BoxSpec      `yaml:"boxes,required"      json:"boxes"`
	UserLists []UserListSpec `yaml:"user_lists"          json:"user_lists"`
	HTTP      HTTPSpec       `yaml:"http,required"       json:"http"`
}

// BoxSpec contains the config specification of one box.
// It is the parent to a list of services.
// The Name is enforced to be unique, but multiple boxes can have identical services or hosts.
type BoxSpec struct {
	Name     types.BoxName `yaml:"name,required"     json:"name"`
	Host     string        `yaml:"host,required"     json:"host"`
	Services []ServiceSpec `yaml:"services,required" json:"services"`
}

// ServiceSpec contains the config specification of one service.
// This always belongs to a box, but the object is not linked directly.
// Service Names are only unique per individual boxes.
// Multiple boxes can have SSH, but one box can't have two SSH-es
// (SSH2 would be f
type ServiceSpec struct {
	Name      types.ServiceName  `yaml:"name,required"           json:"name"`
	CheckName types.CheckName    `yaml:"check,required"          json:"check"`
	UserList  types.UserListName `yaml:"user_list,omitempty"     json:"user_list,omitempty"`
	// Arguments   map[string]any    `yaml,json:"args"`
}

// UserListSpec contains the definition of a userlist.
// Userlists have a unique name.
// Optionally, they can have a Domain or NetBIOS name.
// These apply to every user in the list.
type UserListSpec struct {
	Name        types.UserListName `yaml:"name,required"         json:"name"`
	DomainName  string             `yaml:"domain,omitempty"      json:"domain,omitempty"`
	NetBIOSName string             `yaml:"netbios,omitempty"     json:"netbios,omitempty"`
	Users       []UserSpec         `yaml:"users,required"        json:"users"`
}

// UserSpec is the config specification of a single user.
type UserSpec struct {
	Username string `yaml:"username,required" json:"username"`
	Password string `yaml:"password,required" json:"password"`
}

// HTTPSpec defines the settings applied to the APIs listener.
type HTTPSpec struct {
	Host string `yaml:"host,required" json:"host"`
	Port int    `yaml:"port,required" json:"port"`
}
