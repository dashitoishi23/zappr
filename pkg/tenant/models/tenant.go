package tenantmodels

import (
	"github.com/google/uuid"
)

type Tenant struct {
	Identifier string    `json:"identifier"`
	Name       string    `json:"name" validate:"nonzero"`
}

func(t *Tenant) InitFields() {
	t.Identifier = uuid.New().String()
}