package repository

import (
	"encoding/json"
	"errors"
	"time"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/models"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/util"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type pagedResponse[T any] struct {
	items []T
	page  int
	size  int
}

type BaseCRUD[T any] interface {
	Create(newEntity T) (T, error)
	// GetById(identifier string) (T, error)
	// GetAll(dest interface{}, conds ...interface{}) ([]T, error)
	// GetPaged(skip int, take int, dest interface{}, conds ...interface{}) (pagedResponse[T], error)
	// Update(updatedEntity T) (T, error)
	// Delete(identifier string) bool
}

type baseCRUD[T any] struct {
	db *gorm.DB
	repository IRepository[T]
}

func NewBaseCRUD[T any](database *gorm.DB) BaseCRUD[T] {
	return &baseCRUD[T]{
		db: database,
		repository: Repository[T](database),
	}
}

func (b *baseCRUD[T]) Create(obj T) (T, error) {
	if errs := validator.Validate(obj); errs != nil {
		return obj, errors.New(constants.INVALID_MODEL)
	}

	newEntity := &commonmodels.Entity[T]{
		Identifier: uuid.New().String(),
		CreatedOn: time.Now(),
		T: obj,
	}

	var entityToBeAdded commonmodels.Entity[T]

	encodedBytes, err := util.JsonEncoder(newEntity)

	if err != nil {
		return obj, err
	}

	if err := json.Unmarshal(encodedBytes, &entityToBeAdded); err != nil {
		return obj, err
	}

	tx := b.repository.Add(entityToBeAdded)

	if tx.Error != nil {
		return obj, tx.Error
	}

	return obj, nil
}