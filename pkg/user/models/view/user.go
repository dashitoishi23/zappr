package models

type User struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name" validate:"nonzero"`
	Email      string `json:"email" validate:"nonzero"`
	Password   string `json:"password" validate:"nonzero"`
}