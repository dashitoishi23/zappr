package masterrolemodels

import (
	"time"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Role struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name" validate:"nonzero"`
	CreatedOn        time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn       time.Time `json:"modifiedOn"`
	TenantIdentifier string `json:"tenantIdentifier" validate:"nonzero"`
	Scopes 			pq.StringArray `json:"scopes" validate:"nonzero" gorm:"type:text[]"`
}

func (r *Role) InitFields(requestScope commonmodels.RequestScope) {
	r.Identifier = uuid.New().String()

	r.CreatedOn = time.Now()
	r.ModifiedOn = time.Time{}
	r.TenantIdentifier = requestScope.UserTenant
}

func (r *Role) UpdateFields(createdOn time.Time, requestScope commonmodels.RequestScope){
	r.CreatedOn = createdOn
	r.ModifiedOn = time.Now()
	r.TenantIdentifier = requestScope.UserTenant
}