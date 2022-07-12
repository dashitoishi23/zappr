package usermetadatamodels

import (
	"encoding/json"
	"time"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/google/uuid"
)

type UserMetadata struct {
	Identifier string `json:"identifier"`
	Metadata json.RawMessage `json:"metadata" gorm:"type:jsonb"`
	TenantIdentifier string `json:"tenantIdentifier"`
	EntityName string `json:"entityName"`
	CreatedOn time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

func (u *UserMetadata) InitFields(scope commonmodels.RequestScope){
	u.Identifier = uuid.New().String()

	u.CreatedOn = time.Now()
	u.ModifiedOn = time.Time{}
	u.TenantIdentifier = scope.UserTenant
}

func (u *UserMetadata) UpdateFields(createdOn time.Time, scope commonmodels.RequestScope) {
	u.CreatedOn = createdOn
	u.ModifiedOn = time.Now()

	u.TenantIdentifier = scope.UserTenant
}