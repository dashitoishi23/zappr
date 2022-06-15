package repository

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