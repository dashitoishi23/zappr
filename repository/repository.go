package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type IRepository[T any] interface {
	Add(T any) *gorm.DB
	// AddBulk(newEntity []T) *gorm.DB
	// Find() *gorm.DB
	FindFirst(T interface{}) (T, error)
	// FindByConditions(dest interface{}, conds ...interface{}) *gorm.DB
	// Update(column string, value interface{}) *gorm.DB
	// Delete(identifier string) *gorm.DB
}

type repository[T any] struct {
	db *gorm.DB
}

func Repository[T any](database *gorm.DB) IRepository[T]{
	return &repository[T]{
		db: database,
	}
}

func(r *repository[T]) Add(newEntity any) *gorm.DB{
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




