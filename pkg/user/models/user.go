package usermodels

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Identifier string `json:"identifier" gorm:"type:text"`
	Name       string `json:"name" validate:"nonzero" gorm:"type:text"`
	Email      string `json:"email" validate:"nonzero" gorm:"type:text"`
	Password   string `json:"password" gorm:"type:text"`
	IsADUser   bool   `json:"isAdUser" gorm:"type:bool"`
	Locale     string `json:"locale" validate:"nonzero" gorm:"type:text"`
	TenantIdentifier string `json:"tenantIdentifier" validate:"nonzero" gorm:"type:text"`
	CreatedOn time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn *time.Time `json:"modifiedOn"`
}

func(u *User) InitFields() {
	u.Identifier = uuid.New().String()

	u.CreatedOn = time.Now()
	u.ModifiedOn = nil
}

func(u *User) UpdateFields(createdOn time.Time){
	*u.ModifiedOn = time.Now()

	u.CreatedOn = createdOn
}