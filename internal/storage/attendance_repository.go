package storage

import "github.com/hop-/gotchat/internal/core"

type AttendanceRepository struct {
	Repository[Attendance, core.Attendance]
}

func newAttendanceRepository(storage StorageDb) *AttendanceRepository {
	return &AttendanceRepository{newRepository[Attendance](storage)}
}

func (r *AttendanceRepository) Init() error {
	return r.Repository.Init()
}

func (r *AttendanceRepository) GetOne(id uint) (*core.Attendance, error) {
	m, err := r.getOne(id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AttendanceRepository) GetOneBy(field string, value any) (*core.Attendance, error) {
	if !core.IsFieldExist[core.Attendance](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	m, err := r.getOneBy(field, value)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AttendanceRepository) GetAll() ([]*core.Attendance, error) {
	ms, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var messages []*core.Attendance
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *AttendanceRepository) GetAllBy(field string, value any) ([]*core.Attendance, error) {
	if !core.IsFieldExist[core.Attendance](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	ms, err := r.getAllBy(field, value)
	if err != nil {
		return nil, err
	}

	var messages []*core.Attendance
	for _, m := range ms {
		messages = append(messages, m.ToEntity())
	}

	return messages, nil
}

func (r *AttendanceRepository) Create(entity *core.Attendance) (*core.Attendance, error) {
	m := FromEntityToAttendance(entity)
	_, err := r.create(m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AttendanceRepository) Update(entity *core.Attendance) (*core.Attendance, error) {
	m := FromEntityToAttendance(entity)
	_, err := r.update(m.Id, m)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *AttendanceRepository) Delete(id uint) error {
	return r.delete(id)
}
