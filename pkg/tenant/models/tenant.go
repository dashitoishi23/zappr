package tenantmodels

import "time"

type Tenant struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name" validate:"nonzero"`
	CreatedOn  time.Time   `json:"createdOn" validate:"nonzero"`
	ModifiedOn time.Time 	`json:"modifiedOn"`
}