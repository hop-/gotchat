package core

type Repository[T Entity] interface {
	GetOne(id string) (T, error)
	GetAll() ([]T, error)
	Save(entity T) error
	Update(entity T) error
	Delete(id string) error
}
