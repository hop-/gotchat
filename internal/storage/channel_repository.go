package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type ChannelRepository struct {
	StorageDb
}

func newChannelRepository(storage StorageDb) *ChannelRepository {
	return &ChannelRepository{storage}
}

func (r *ChannelRepository) GetOne(id int) (*core.Channel, error) {
	row := r.Db().QueryRow("SELECT id, unique_id, name FROM channels WHERE id = ?", id)
	if row == nil {
		return nil, core.ErrEntityNotFound
	}

	var ch core.Channel

	err := row.Scan(&ch.Id, &ch.UniqueId, &ch.Name)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

func (r *ChannelRepository) GetOneBy(field string, value any) (*core.Channel, error) {
	if !isFieldExist[core.Channel](field) {
		return nil, core.ErrEntityFieldNotExist
	}
	row := r.Db().QueryRow("SELECT id, unique_id, name FROM channels WHERE "+field+" = ?", value)
	if row == nil {
		return nil, core.ErrEntityNotFound
	}

	var ch core.Channel
	err := row.Scan(&ch.Id, &ch.UniqueId, &ch.Name)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

func (r *ChannelRepository) GetAll() ([]*core.Channel, error) {
	rows, err := queryWithRetry(r.Db(), "SELECT id, unique_id, name FROM channels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*core.Channel
	for rows.Next() {
		var ch core.Channel
		err := rows.Scan(&ch.Id, &ch.UniqueId, &ch.Name)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}

	return channels, nil
}

func (r *ChannelRepository) GetAllBy(field string, value any) ([]*core.Channel, error) {
	if !isFieldExist[core.Channel](field) {
		return nil, core.ErrEntityFieldNotExist
	}

	rows, err := queryWithRetry(r.Db(), "SELECT id, unique_id, name FROM channels where "+field+" = ?", value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*core.Channel
	for rows.Next() {
		var ch core.Channel
		err := rows.Scan(&ch.Id, &ch.UniqueId, &ch.Name)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}

	return channels, nil
}

func (r *ChannelRepository) Create(channel *core.Channel) error {
	_, err := execWithRetry(
		r.Db(),
		"INSERT INTO channels (unique_id, name) VALUES (?, ?)",
		channel.UniqueId,
		channel.Name,
	)

	return err
}

func (r *ChannelRepository) Update(channel *core.Channel) error {
	_, err := execWithRetry(
		r.Db(),
		"UPDATE channels SET unique_id = ?, name = ? WHERE id = ?",
		channel.UniqueId,
		channel.Name,
		channel.Id,
	)

	return err
}

func (r *ChannelRepository) Delete(id int) error {
	_, err := execWithRetry(r.Db(), "DELETE FROM channels WHERE id = ?", id)

	return err
}
