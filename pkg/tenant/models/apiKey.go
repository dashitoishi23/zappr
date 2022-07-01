package tenantmodels

import (
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name"`
	Secret           string `json:"secret"`
	TenantIdentifier string `json:"tenantIdentifier"`
	CreatedOn        time.Time `json:"createdOn"`
	ModifiedOn 		time.Time `json:"modifiedOn"`
}

func (a *APIKey) InitFields(){
	a.Identifier = uuid.New().String()

	a.CreatedOn = time.Now()
	a.ModifiedOn = time.Time{}
}

func (a *APIKey) UpdateFields(createdOn time.Time){
	a.ModifiedOn = time.Now()

	a.CreatedOn = time.Time{}
}