package repository

import (
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
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
	GetFirst(obj interface{}) (T, error)
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

	tx := b.repository.Add(obj)

	if tx.Error != nil {
		return obj, tx.Error
	}

	return obj, nil
}

func (b *baseCRUD[T]) GetFirst(obj interface{}) (T, error) {
	var result T

	result, err := b.repository.FindFirst(obj)

	if err != nil {
		return result, err
	}

	return result, nil
}