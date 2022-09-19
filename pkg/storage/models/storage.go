package storagemodels

import "time"

type StaticStorage struct {
	Identifier       string `json:"identifier"`
	URI              string `json:"uri"`
	TenantIdentifier string `json:"tenantIdentifier"`
	CreatedOn        time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}