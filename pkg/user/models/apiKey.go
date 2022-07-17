package usermodels

import (
	"time"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type APIKey struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name"`
	Secret           string `json:"secret"`
	TenantIdentifier string `json:"tenantIdentifier"`
	UserIdentifier string `json:"userIdentifier"`
	Scopes 			pq.StringArray `json:"scopes" gorm:"type:text[]"`
	CreatedOn        time.Time `json:"createdOn"`
	ModifiedOn 		time.Time `json:"modifiedOn"`
}

func (a *APIKey) InitFields(requestScope commonmodels.RequestScope){
	a.Identifier = uuid.New().String()

	a.TenantIdentifier = requestScope.UserTenant

	a.CreatedOn = time.Now()
	a.ModifiedOn = time.Time{}
}

func (a *APIKey) UpdateFields(createdOn time.Time){
	a.ModifiedOn = time.Now()

	a.CreatedOn = time.Time{}
}