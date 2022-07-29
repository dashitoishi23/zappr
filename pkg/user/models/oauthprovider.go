package usermodels

import (
	"encoding/json"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	Identifier string `json:"identifier" gorm:"primaryKey"`
	Name       string `json:"name"`
	Metadata   json.RawMessage `json:"metadata" gorm:"type:jsonb"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

func (o *OAuthProvider) InitFields() {
	o.Identifier = uuid.New().String()
}