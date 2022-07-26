package usermodels

import (
	"encoding/json"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Metadata   json.RawMessage `json:"metadata" gorm:"type:jsonb"` 
}

func (o *OAuthProvider) InitFields() {
	o.Identifier = uuid.New().String()
}