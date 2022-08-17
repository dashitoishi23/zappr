package usermodels

import "encoding/json"

type UpdateCurrentUser struct {
	Name              string          `json:"name" validate:"nonzero"`
	ProfilePictureURL string          `json:"profilePictureUrl" validate:"nonzero"`
	Metadata          json.RawMessage `json:"metadata" validate:"nonzero"`
	Email      string `json:"email" validate:"nonzero"`
}