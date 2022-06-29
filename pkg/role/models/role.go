package masterrolemodels

import (
	"time"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
	"github.com/google/uuid"
)

type Role struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name" validate:"nonzero"`
	CreatedOn        time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn       time.Time `json:"modifiedOn"`
	TenantIdentifier string `json:"tenantIdentifier" validate:"nonzero"`
}

func (r *Role) InitFields() {
	r.Identifier = uuid.New().String()

	r.CreatedOn = time.Now()
	r.ModifiedOn = time.Time{}
	r.TenantIdentifier = state.GetState().GetUserContext().UserTenant
}

func (r *Role) UpdateFields(createdOn time.Time){
	r.CreatedOn = createdOn
	r.ModifiedOn = time.Now()
	r.TenantIdentifier = state.GetState().UserContext.UserTenant
}