package core

type Repository[T Entity] interface {
	GetOne(id int) (*T, error)
	GetAll() ([]*T, error)
	GetAllBy(field string, value any) ([]*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id int) error
}
