package storage

import (
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"gorm.io/gorm"
)

type GormRepo[T Model[E], E core.Entity] interface {
	core.Repository[E]
	Init() error
	getMappedField(field string) string
}

type Repository[T Model[E], E core.Entity] struct {
	StorageDb
	fieldMapper map[string]string
	mutex       sync.Mutex
}

func newRepository[T Model[E], E core.Entity](storage StorageDb) Repository[T, E] {
	return Repository[T, E]{StorageDb: storage}
}

func (r *Repository[T, E]) Init() error {
	// Initialize field mapper
	r.mutex.Lock()
	defer r.mutex.Unlock()

	stat := gorm.Statement{DB: r.Db()}

	var e E
	err := stat.Parse(e)
	if err != nil {
		return err
	}

	entityFields := core.ListFieldNames[E]()
	r.fieldMapper = make(map[string]string, len(entityFields))

	for f := range entityFields {
		modelFieldName := getMappedFieldForType[T](f)
		tableFieldName, ok := stat.Schema.FieldsByName[modelFieldName]
		if !ok {
			return fmt.Errorf("field %s does not exist in model %T", modelFieldName, e)
		}

		r.fieldMapper[f] = tableFieldName.DBName
	}

	return nil
}

func (r *Repository[T, E]) getMappedField(field string) string {
	if r.fieldMapper == nil {
		// Not initialized
		return ""
	}

	return r.fieldMapper[field]
}

func (r *Repository[T, E]) getOne(id uint) (*T, error) {
	var m T
	err := r.Db().First(&m, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrEntityNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *Repository[T, E]) getOneBy(field string, value any) (*T, error) {
	dbField := r.getMappedField(field)

	var m T
	err := r.Db().Where(fmt.Sprintf("%s = ?", dbField), value).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrEntityNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *Repository[T, E]) getAll() ([]T, error) {
	var ms []T
	err := r.Db().Find(&ms).Error

	return ms, err
}

func (r *Repository[T, E]) getAllBy(field string, value any) ([]T, error) {
	dbField := r.getMappedField(field)

	var ms []T
	err := r.Db().Where(fmt.Sprintf("%s = ?", dbField), value).Find(&ms).Error

	return ms, err
}

func (r *Repository[T, E]) create(m *T) (*T, error) {
	err := r.Db().Create(m).Error

	return m, err
}

func (r *Repository[T, E]) update(id uint, m *T) (*T, error) {
	err := r.Db().Where("id = ?", id).Updates(m).Error

	return m, err
}

func (r *Repository[T, E]) delete(id uint) error {
	var m T
	err := r.Db().Where("id = ?", id).Delete(&m).Error

	return err
}
