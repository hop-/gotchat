package core

type Repository[T Entity] interface {
	GetOne(id uint) (*T, error)
	GetOneBy(field string, value any) (*T, error)
	GetAll() ([]*T, error)
	GetAllBy(field string, value any) ([]*T, error)
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Delete(id uint) error
}
