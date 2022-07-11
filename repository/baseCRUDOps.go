package repository

import (
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"gopkg.in/validator.v2"
)

type BaseCRUD[T any] interface {
	Create(newEntity T) (T, error)
	GetFirst(obj interface{}) (T, error)
	GetFirstAsync(obj interface{}, entity chan T, txnError chan error)
	GetAll() ([]T, error)
	GetAllByTenant() ([]T, error)
	Find(obj interface{}) ([]T, error)
	GetPagedAsync(obj interface{}, page int, size int, pagedResponse chan commonmodels.PagedResponse[T], 
		txnError chan error)
	Update(updatedEntity T) (T, error)
	Delete(currentEntity interface{}) (bool, error)
	QueryRawSql(sql string, conditions ...interface{}) ([]T, error)
}

type baseCRUD[T any] struct {
	repository IRepository[T]
}

func NewBaseCRUD[T any](repository IRepository[T]) BaseCRUD[T] {
	return &baseCRUD[T]{
		repository: repository,
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

func (b *baseCRUD[T]) GetAllByTenant() ([]T, error) {
	result, err := b.repository.GetAllByTenant()

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

func (b *baseCRUD[T]) Delete(currentEntity interface{}) (bool, error) {
	result, err := b.repository.Delete(currentEntity)

	if err != nil {
		return false, err
	}

	return result, nil
}

func (b *baseCRUD[T]) GetPagedAsync(obj interface{}, page int, size int, pagedResponse chan commonmodels.PagedResponse[T], 
	txnError chan error) {

	result, err := b.repository.GetPaged(obj, page, size)

	pagedResponse <- result
	txnError <- err
}

func (b *baseCRUD[T]) QueryRawSql(sql string, conditions ...interface{}) ([]T, error) {
	result, err := b.repository.QueryRawSql(sql, conditions)

	return result, err
}