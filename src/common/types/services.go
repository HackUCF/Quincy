package types

// ServiceTemplate contains *almost* everything an agent needs to run a check.
// It adds the host variable from the box, and a specific teams number.
// It is only missing a user, as this needs to be pulled from the db on demand.
type ServiceTemplate struct {
	Name      ServiceName  `json:"name"`
	CheckName CheckName    `json:"check"`
	UserList  UserListName `json:"user_list,omitempty,omitzero"`
	BoxName   BoxName      `json:"box"`
	// Arguments   map[string]any `json:"args,omitempty,omitzero"`

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
