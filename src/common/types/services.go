package types

// ServiceSpec contains the config specification of one service.
// This always belongs to a box, but the object is not linked directly.
// Service Names are only unique per individual box.
type ServiceSpec struct {
	Name      ServiceName  `yaml:"name"      mapstructure:"name"      json:"name"`
	CheckName CheckName    `yaml:"check"     mapstructure:"check"     json:"check"`
	UserList  UserListName `yaml:"user_list" mapstructure:"user_list" json:"user_list,omitempty"`
	Timeout   float64      `yaml:"timeout"   mapstructure:"timeout"   json:"timeout,omitempty"`
}

// ServiceTemplate contains *almost* everything an agent needs to run a check.
// It adds variables from the box, and a specific team number.
// It is only missing a user, as this needs to be pulled from the db on demand.
type ServiceTemplate struct {
	ServiceSpec

	BoxName BoxName `json:"box"`
	Host    string  `json:"host"`
	TeamNum TeamNum `json:"team_num"`
}

// User contains all the information relating to a single user.
type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DomainName  string `json:"domain,omitempty,omitzero"`
	NetBIOSName string `json:"netbios,omitempty,omitzero"`
}

// Service is a fully rendered service ready to be run by the agent.
// Is just a ServiceTemplate with an added user.
type Service struct {
	ServiceTemplate
	User *User `json:"user,omitempty,omitzero"`
}
