package usermodels

import "encoding/json"

type UpdateUserMetadata struct {
	Identifier string `json:"identifier" gorm:"primaryKey"`
	Metadata   json.RawMessage `json:"metadata"`
	TenantIdentifier string `json:"tenantIdentifier"`
}