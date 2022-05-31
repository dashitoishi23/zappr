package models

type User struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}