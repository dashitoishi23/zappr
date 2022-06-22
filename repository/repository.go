package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type IRepository[T any] interface {
	Add(T T) *gorm.DB
	// AddBulk(newEntity []T) *gorm.DB
	GetAll() ([]T, error)
	FindFirst(T interface{}) (T, error)
	Find(T interface{}) ([]T, error)
	Update(T T) (T, error)
	Delete(currentEntity *T, identifier string) (bool, error)
}

type repository[T any] struct {
	db *gorm.DB
}

func Repository[T any](database *gorm.DB) IRepository[T]{
	return &repository[T]{
		db: database,
	}
}

func(r *repository[T]) Add(newEntity T) *gorm.DB{
	tx := r.db.Create(newEntity)
	fmt.Print(tx.RowsAffected)
	return tx
}

func(r *repository[T]) FindFirst(currentEntity interface{}) (T, error){
	var result T
	tx := r.db.Where(currentEntity).First(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	return result, nil

}

func(r *repository[T]) GetAll() ([]T, error) {
	var result []T
	tx := r.db.Find(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	return result, nil
}

func(r *repository[T]) Find(currentEntity interface{}) ([]T, error) {
	var result []T
	tx := r.db.Where(currentEntity).Find(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	return result, nil
}

func(r *repository[T]) Update(newEntity T) (T, error) {
	tx := r.db.Save(newEntity)

	if tx.Error != nil {
		return newEntity, tx.Error
	}

	return newEntity, nil
}

func(r *repository[T]) Delete(currentEntity *T, identifier string) (bool, error) {
	tx := r.db.Delete(currentEntity, identifier)

	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil
}




