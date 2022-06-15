package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type IRepository[T any] interface {
	Add(T any) *gorm.DB
	// AddBulk(newEntity []T) *gorm.DB
	// Find() *gorm.DB
	// FindByIdentifer(identifier string) *gorm.DB
	// FindByConditions(dest interface{}, conds ...interface{}) *gorm.DB
	// Update(column string, value interface{}) *gorm.DB
	// Delete(identifier string) *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func Repository[T any](database *gorm.DB) IRepository[T]{
	return &repository{
		db: database,
	}
}

func(r *repository) Add(newEntity any) *gorm.DB{
	tx := r.db.Create(newEntity)
	fmt.Print(tx.RowsAffected)
	return tx
}




