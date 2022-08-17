package usermodels

import "encoding/json"

type UpdateCurrentUser struct {
	Identifier 		  string 		`json:"identifier"`
	Name              string          `json:"name"`
	ProfilePictureURL string          `json:"profilePictureUrl"`
	Metadata          json.RawMessage `json:"metadata"`
	Email      string `json:"email"`
	TenantIdentifier string `json:"tenantIdentifier"`
}