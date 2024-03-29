package tenantmodels

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	Identifier string    `json:"identifier" gorm:"primaryKey"`
	Name       string    `json:"name" validate:"nonzero"`
	AdminEmail string `json:"adminEmail" validate:"nonzero"`
	CreatedOn time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

func(t *Tenant) InitFields() {
	t.Identifier = uuid.New().String()

	t.CreatedOn = time.Now()
	t.ModifiedOn = time.Time{}
}

func(t *Tenant) UpdateFields(createdOn time.Time){
	t.ModifiedOn = time.Now()

	t.CreatedOn = createdOn
}