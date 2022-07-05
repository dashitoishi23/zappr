package userrolemodels

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRole struct {
	Identifier     string    `json:"identifier" gorm:"primaryKey"`
	UserIdentifier string    `json:"userIdentifier"`
	RoleIdentifier string    `json:"roleIdentifier"`
	Scopes 			pq.StringArray `json:"scopes" gorm:"type:text[]"`
	CreatedOn      time.Time `json:"createdOn"`
	ModifiedOn     time.Time `json:"modifiedOn"`
}

func (u *UserRole) InitFields() {
	u.Identifier = uuid.New().String()

	u.CreatedOn = time.Now()
	u.ModifiedOn = time.Time{}
}

func (u *UserRole) UpdateFields(createdOn time.Time) {
	u.ModifiedOn = time.Now()

	u.CreatedOn = createdOn
}