package repository

import (
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type BaseCRUD[T any] interface {
	Create(newEntity T) (T, error)
	GetFirst(obj interface{}) (T, error)
	GetFirstAsync(obj interface{}, entity chan T, txnError chan error)
	GetAll() ([]T, error)
	Find(obj interface{}) ([]T, error)
	// GetPaged(skip int, take int, dest interface{}, conds ...interface{}) (pagedResponse[T], error)
	Update(updatedEntity T) (T, error)
	Delete(currentEntity *T, identifier string) (bool, error)
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
	result, err := b.repository.FindFirst(obj)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *baseCRUD[T]) GetFirstAsync(obj interface{}, entity chan T, txnError chan error){
	result, err := b.repository.FindFirst(obj)

	entity <- result
	txnError <- err
}

func (b *baseCRUD[T]) GetAll() ([]T, error) {
	result, err := b.repository.GetAll()

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *baseCRUD[T]) Find(obj interface{}) ([]T, error) {
	result, err := b.repository.Find(obj)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *baseCRUD[T]) Update(currentEntity T) (T, error) {
	result, err := b.repository.Update(currentEntity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *baseCRUD[T]) Delete(currentEntity *T, identifier string) (bool, error) {
	result, err := b.repository.Delete(currentEntity, identifier)

	if err != nil {
		return false, err
	}

	return result, nil
}