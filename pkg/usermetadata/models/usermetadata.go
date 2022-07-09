package usermetadatamodels

import (
	"time"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
	"github.com/google/uuid"
)

type UserMetadata struct {
	Identifier string `json:"identifier"`
	Metadata map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	TenantIdentifier string `json:"tenantIdentifier"`
	EntityName string `json:"entityName"`
	CreatedOn time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

func (u *UserMetadata) InitFields(){
	u.Identifier = uuid.New().String()

	u.CreatedOn = time.Now()
	u.ModifiedOn = time.Time{}
	u.TenantIdentifier = state.GetState().UserContext.UserTenant
}

func (u *UserMetadata) UpdateFields(createdOn time.Time) {
	u.CreatedOn = createdOn
	u.ModifiedOn = time.Now()

	u.TenantIdentifier = state.GetState().UserContext.UserTenant
}