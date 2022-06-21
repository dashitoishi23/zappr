package tenantmodels

import "time"

type SearchableTenant struct {
	Identifier string     `json:"identifier"`
	Name       string     `json:"name"`
	CreatedOn  time.Time  `json:"createdOn"`
	ModifiedOn *time.Time `json:"modifiedOn"`
}