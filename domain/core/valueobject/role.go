package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Role string
type Roles []Role

const (
	UserRole       Role = "user"
	AdminRole      Role = "admin"
	SupervisorRole Role = "supervisor"
)

func (r Role) String() string {
	return string(r)
}

func ConvertRolesToString(roles []Role) []string {
	stringRoles := make([]string, len(roles))

	for i, role := range roles {
		stringRoles[i] = string(role)
	}
	return stringRoles
}

func ConvertStringToRoles(roles []string) []Role {
	roleEnums := make([]Role, len(roles))

	for i, role := range roles {
		roleEnums[i] = Role(role)
	}
	return roleEnums
}

func (roles Roles) Value() (driver.Value, error) {
	return json.Marshal(roles)
}

func (roles Roles) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(data, &roles)
}
