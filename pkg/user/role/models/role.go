package rolemodels

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Identifier     string    `json:"identifier"`
	Name           string    `json:"name" validate:"nonzero"`
	UserIdentifier string    `json:"userIdentifier" validate:"nonzero"`
	RoleIdentifier string 	 `json:"roleIdentifier" validate:"nonzero"`
	CreatedOn      time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn 		time.Time 	`json:"modifiedOn"`
}

func(r *Role) InitFields() {
	r.Identifier = uuid.New().String()

	r.CreatedOn = time.Now()
	r.ModifiedOn = time.Time{}
}

func(r *Role) UpdateFields(createdOn time.Time){
	r.ModifiedOn = time.Now()

	r.CreatedOn = createdOn
}