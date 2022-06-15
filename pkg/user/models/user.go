package usermodels

import "github.com/google/uuid"

type User struct {
	Identifier string `json:"identifier" gorm:"type:text"`
	Name       string `json:"name" validate:"nonzero" gorm:"type:text"`
	Email      string `json:"email" validate:"nonzero" gorm:"type:text"`
	Password   string `json:"password" gorm:"type:text"`
	IsADUser   bool   `json:"isAdUser" gorm:"type:bool"`
	Locale     string `json:"locale" validate:"nonzero" gorm:"type:text"`
}

func (u *User) BeforeCreate(err error) {
	u.Identifier = uuid.New().String()

	return
}