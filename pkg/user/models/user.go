package usermodels

import (
	"time"

	rolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/role/models"
	"github.com/google/uuid"
)

type User struct {
	Identifier string `json:"identifier" `
	Name       string `json:"name" validate:"nonzero" `
	Email      string `json:"email" validate:"nonzero" `
	Password   string `json:"password" `
	IsADUser   bool   `json:"isAdUser"`
	Locale     string `json:"locale" validate:"nonzero" `
	TenantIdentifier string `json:"tenantIdentifier" validate:"nonzero" `
	Role 			rolemodels.Role  `json:"role" gorm:"foreignKey:UserIdentifier"`
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