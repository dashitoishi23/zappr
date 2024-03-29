package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"gorm.io/gorm"
)

type IRepository[T any] interface {
	AddRange(T []T) *gorm.DB
	Add(T T) *gorm.DB
	GetAll() ([]T, error)
	GetAllByAssociation(associationName string) ([]T, error)
	GetAllByTenant(ctx context.Context) ([]T, error)
	FindFirst(T interface{}) (T, error)
	FindFirstByAssociation(associationName string, T interface{}) (T, error)
	Find(T interface{}) ([]T, error)
	FindByAssociation(associationName string, T interface{}) ([]T, error)
	GetPaged(T interface{}, page int, size int) (commonmodels.PagedResponse[T], error)
	Update(T T) (T, error)
	Delete(currentEntity interface{}) (bool, error)
	ExecuteRawQuery(sql string, values ...interface{}) (bool, error)
	QueryRawSql(sql string, conditions ...interface{}) ([]T, error)
	QueryRawSqlPaged(sql string, page int, size int, conditions ...interface{}) (commonmodels.PagedResponse[T], error)
	GetTransaction() *gorm.DB
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
	tx := r.db.Create(&newEntity)
	fmt.Print(tx.RowsAffected)
	return tx
}

func(r *repository[T]) AddRange(newEntity []T) *gorm.DB{
	r.db.Begin()

	for i:=0; i<len(newEntity); i++ {
		tx:= r.db.Create(newEntity[i])

		if tx.Error != nil {
			return tx
		}
	}

	tx := r.db.Commit()

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

func(r *repository[T]) FindFirstByAssociation(associationName string, currentEntity interface{}) (T, error) {
	var result T
	tx := r.db.Joins(associationName).First(&result, currentEntity)

	return result, tx.Error
}

func(r *repository[T]) GetAll() ([]T, error) {
	var result []T
	tx := r.db.Find(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	return result, nil
}

func(r *repository[T]) GetAllByAssociation(associationName string) ([]T, error) {
	var result []T
	tx := r.db.Joins(associationName).Find(&result)

	return result, tx.Error
}

func(r *repository[T]) GetAllByTenant(ctx context.Context) ([]T, error) {
	var result []T
	tx := r.db.Where(map[string]interface{}{
		"TenantIdentifier": ctx.Value("requestScope").(commonmodels.RequestScope).UserTenant,
	}).Find(&result)

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

func(r *repository[T]) FindByAssociation(associationName string, currentEntity interface{}) ([]T, error) {
	var result []T
	tx := r.db.Joins(associationName).Find(&result, currentEntity)

	return result, tx.Error
}

func(r *repository[T]) Update(newEntity T) (T, error) {
	tx := r.db.Updates(newEntity)

	if tx.Error != nil {
		return newEntity, tx.Error
	}

	return newEntity, nil
}

func(r *repository[T]) Delete(currentEntity interface{}) (bool, error) {
	tx := r.db.Delete(currentEntity)

	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil
}

func(r *repository[T]) GetPaged(currentEntity interface{}, page int, size int) (commonmodels.PagedResponse[T], error) {
	skip := (page - 1) * size
	var result []T
	tx := r.db.Where(currentEntity).Find(&result)

	//Since somehow slices change length EVEN WHEN YOU READ THEM LMAO
	resultOfDBOp := len(result)

	if skip > resultOfDBOp {
		skip = resultOfDBOp
	}

	end := skip + size

	if end > resultOfDBOp {
		end = resultOfDBOp
	}

	return commonmodels.PagedResponse[T]{
		Items: result[skip:end],
		Page: page,
		Size: size,
		Pages: int(math.Ceil(float64(len(result))/float64(size))),
	}, tx.Error	

}

func (r *repository[T]) ExecuteRawQuery(sql string, values ...interface{}) (bool, error) {
	tx := r.db.Exec(sql, values)
	
	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, errors.New(constants.RECORD_NOT_FOUND)
	}

	return true, nil
}

func (r *repository[T]) QueryRawSql(sql string, conditions ...interface{}) ([]T, error) {
	var res []T
	tx := r.db.Raw(sql, conditions...).Scan(&res)

	return res, tx.Error
}

func (r *repository[T]) QueryRawSqlPaged(sql string, page int, size int, 
	conditions ...interface{}) (commonmodels.PagedResponse[T], error) {

	var result []T
	tx := r.db.Raw(sql, conditions...).Scan(&result)

	skip := (page - 1) * size

	resultOfDBOp := len(result)

	if skip > resultOfDBOp {
		skip = resultOfDBOp
	}

	end := skip + size

	if end > resultOfDBOp {
		end = resultOfDBOp
	}


	if tx.Error != nil {
		return commonmodels.PagedResponse[T]{
			Items: result,
			Page: page,
			Size: size,
			Pages: 0,
		}, tx.Error
	}

	return commonmodels.PagedResponse[T]{
		Items: result[skip:end],
		Page: page,
		Size: size,
		Pages: int(math.Ceil(float64(len(result))/float64(size))),
	}, nil	


}

func (r *repository[T]) GetTransaction() *gorm.DB {
	tx := r.db.Begin()

	return tx
}



