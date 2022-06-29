package commonmodels

import (
	"time"

	"github.com/google/uuid"
)

type Entity[T any] struct {
	Identifier string `json:"identifier"`
	Entity T `json:"entity"`
	CreatedOn  time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

func (e *Entity[T]) InitFields(){
	e.Identifier = uuid.NewString()

	e.CreatedOn = time.Now()
	e.ModifiedOn = time.Time{}
}

func (e *Entity[T]) UpdateFields(createdOn time.Time){
	e.ModifiedOn = time.Now()

	e.CreatedOn = createdOn
}