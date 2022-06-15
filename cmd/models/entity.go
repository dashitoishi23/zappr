package commonmodels

import "time"

type Entity[T any] struct {
	Identifier string    `json:"identifier"`
	CreatedOn  time.Time `json:"createdOn" validate:"nonzero"`
	ModifiedOn time.Time `json:"modifiedOn" validate:"nonzero"`
	T          T
}