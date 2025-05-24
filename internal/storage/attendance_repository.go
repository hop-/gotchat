package storage

import "github.com/hop-/gotchat/internal/core"

type AttendanceRepository struct {
	StorageDb
}

func newAttendanceRepository(storage StorageDb) *AttendanceRepository {
	return &AttendanceRepository{storage}
}

func (a *AttendanceRepository) GetOne(id int) (*core.Attendance, error) {
	row := a.Db().QueryRow("SELECT id, user_id, channel_id, joined_at FROM attendances WHERE id = ?", id)
	if row == nil {
		return nil, ErrNotFound
	}

	var att core.Attendance

	err := row.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (a *AttendanceRepository) GetOneBy(field string, value any) (*core.Attendance, error) {
	if !isFieldExist[core.Attendance](field) {
		return nil, ErrFieldNotExist
	}

	row := a.Db().QueryRow("SELECT id, user_id, channel_id, joined_at FROM attendances WHERE "+field+" = ?", value)
	if row == nil {
		return nil, ErrNotFound
	}

	var att core.Attendance
	err := row.Scan(&att.Id, &att.UserId, &att.ChannelId, &att.JoinedAt)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (a *AttendanceRepository) GetAll() ([]*core.Attendance, error) {
	rows, err := a.Db().Query("SELECT id, user_id, channel_id, joined_at FROM attendances")
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

func (a *AttendanceRepository) GetAllBy(field string, value any) ([]*core.Attendance, error) {
	if !isFieldExist[core.Attendance](field) {
		return nil, ErrFieldNotExist
	}

	rows, err := a.Db().Query("SELECT id, user_id, channel_id, joined_at FROM attendances WHERE "+field+" = ?", value)
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

func (a *AttendanceRepository) Create(entity *core.Attendance) error {
	_, err := a.Db().Exec(
		"INSERT INTO attendances (user_id, channel_id, joined_at) VALUES (?, ?, ?)",
		entity.UserId,
		entity.ChannelId,
		entity.JoinedAt,
	)

	return err
}

func (a *AttendanceRepository) Update(entity *core.Attendance) error {
	_, err := a.Db().Exec(
		"UPDATE attendances SET user_id = ?, channel_id = ?, joined_at = ? WHERE id = ?",
		entity.UserId,
		entity.ChannelId,
		entity.JoinedAt,
		entity.Id,
	)

	return err
}

func (a *AttendanceRepository) Delete(id int) error {
	_, err := a.Db().Exec("DELETE FROM attendances WHERE id = ?", id)

	return err
}
