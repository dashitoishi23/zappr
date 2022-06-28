package rolemodels

import "time"

type SearchableRole struct {
	Identifier     string    `json:"identifier"`
	Name           string    `json:"name"`
	UserIdentifier string    `json:"userIdentifier"`
	RoleIdentifier string	 `json:"roleIdentifier"`
	CreatedOn      time.Time `json:"createdOn"`
	ModifiedOn     time.Time `json:"modifiedOn"`
}