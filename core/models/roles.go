package models

type Roles struct {
	Role string `json:"status_name"`
}

type RolesSelection struct {
	Roles []Roles
}
