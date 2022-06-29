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
	GetAllByTenant() ([]T, error)
	FindFirst(T interface{}) (T, error)
	Find(T interface{}) ([]T, error)
	GetPaged(T interface{}, page int, size int) (commonmodels.PagedResponse[T], error)
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

func(r *repository[T]) Update(newEntity T) (T, error) {
	tx := r.db.Updates(newEntity)

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

func(r *repository[T]) GetPaged(currentEntity interface{}, page int, size int) (commonmodels.PagedResponse[T], error) {
	skip := (page - 1) * size
	var result []T
	tx := r.db.Where(currentEntity).Find(&result)

	//Since somehow slices change length EVEN WHEN YOU READ THEM LMAO
	resultOfDBOp := len(result)

	//SAME INSANITY AS LINE 96
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



