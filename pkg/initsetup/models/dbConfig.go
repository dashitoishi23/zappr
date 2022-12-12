package initsetupmodels

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DBConfig struct {
	Identifier string `json:"identifier"`
	Config json.RawMessage `json:"config" gorm:"type:jsonb"`
	CreatedOn        time.Time `json:"createdOn"`
	ModifiedOn 		 time.Time `json:"modifiedOn"`
}

func (r *DBConfig) InitFields() {
	r.Identifier = uuid.New().String()

	r.CreatedOn = time.Now()
	r.ModifiedOn = time.Time{}
}

func (r *DBConfig) UpdateFields(createdOn time.Time){
	r.CreatedOn = createdOn
	r.ModifiedOn = time.Now()
}