package tenantmodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	Identifier string    `json:"identifier"`
	Name       string    `json:"name" validate:"nonzero"`
}

func(t *Tenant) BeforeCreate(tx *gorm.DB) (err error){
	t.Identifier = uuid.New().String()

	return
}