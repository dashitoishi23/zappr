package repository

import (
	"fmt"
	"math"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
	"gorm.io/gorm"
)

type IRepository[T any] interface {
	Add(T T) *gorm.DB
	GetAll() ([]T, error)
	GetAllByAssociation(associationName string) ([]T, error)
	GetAllByTenant() ([]T, error)
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

func(r *repository[T]) GetAllByTenant() ([]T, error) {
	var result []T
	tx := r.db.Where(map[string]interface{}{
		"TenantIdentifier": state.GetState().UserContext.UserTenant,
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

func (r *repository[T]) ExecuteRawQuery(sql string, values ...interface{}) (bool, error) {
	tx := r.db.Exec(sql, values)
	
	if tx.Error != nil {
		return false, tx.Error
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



