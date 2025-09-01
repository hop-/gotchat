package storage

import "github.com/hop-/gotchat/internal/core"

type AttendanceRepository struct {
	StorageDb
}

func newAttendanceRepository(storage StorageDb) *AttendanceRepository {
	return &AttendanceRepository{storage}
}

func (r *AttendanceRepository) GetOne(id int) (*core.Attendance, error) {
	row := r.Db().QueryRow("SELECT id, user_id, channel_id, joined_at FROM attendances WHERE id = ?", id)
	if row == nil {
		return nil, core.ErrEntityNotFound
	}

	var att core.Attendance

	err := row.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (r *AttendanceRepository) GetOneBy(field string, value any) (*core.Attendance, error) {
	if !isFieldExist[core.Attendance](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	row := r.Db().QueryRow("SELECT id, user_id, channel_id, joined_at FROM attendances WHERE "+field+" = ?", value)
	if row == nil {
		return nil, core.ErrEntityNotFound
	}

	var att core.Attendance
	err := row.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (r *AttendanceRepository) GetAll() ([]*core.Attendance, error) {
	rows, err := queryWithRetry(r.Db(), "SELECT id, user_id, channel_id, joined_at FROM attendances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*core.Attendance
	for rows.Next() {
		var att core.Attendance
		err := rows.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &att)
	}

	return attendances, nil
}

func (r *AttendanceRepository) GetAllBy(field string, value any) ([]*core.Attendance, error) {
	if !isFieldExist[core.Attendance](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	rows, err := queryWithRetry(r.Db(), "SELECT id, user_id, channel_id, joined_at FROM attendances WHERE "+field+" = ?", value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*core.Attendance
	for rows.Next() {
		var att core.Attendance
		err := rows.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &att)
	}

	return attendances, nil
}

func (r *AttendanceRepository) Create(entity *core.Attendance) error {
	_, err := execWithRetry(
		r.Db(),
		"INSERT INTO attendances (user_id, channel_id, joined_at) VALUES (?, ?, ?)",
		entity.UserId,
		entity.ChannelId,
		entity.JoinedAt,
	)

	return err
}

func (r *AttendanceRepository) Update(entity *core.Attendance) error {
	_, err := execWithRetry(
		r.Db(),
		"UPDATE attendances SET user_id = ?, channel_id = ?, joined_at = ? WHERE id = ?",
		entity.UserId,
		entity.ChannelId,
		entity.JoinedAt,
		entity.Id,
	)

	return err
}

func (r *AttendanceRepository) Delete(id int) error {
	_, err := execWithRetry(r.Db(), "DELETE FROM attendances WHERE id = ?", id)

	return err
}
