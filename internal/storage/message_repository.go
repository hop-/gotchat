package storage

import (
	"github.com/hop-/gotchat/internal/core"
)

type MessageRepository struct {
	StorageDb
}

func newMessageRepository(storage StorageDb) *MessageRepository {
	return &MessageRepository{storage}
}

func (r *MessageRepository) GetOne(id int) (*core.Message, error) {
	row := r.Db().QueryRow("SELECT * FROM messages WHERE id = ?", id)
	if row == nil {
		return nil, ErrNotFound
	}
	var message core.Message
	err := row.Scan(&message.Id, &message.UserId, &message.ChannelId, &message.Text, &message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) GetOneBy(field string, value any) (*core.Message, error) {
	if !isFieldExist[core.Channel](field) {
		return nil, ErrFieldNotExist
	}

	row := r.Db().QueryRow("SELECT * FROM messages WHERE "+field+" = ?", value)
	if row == nil {
		return nil, ErrNotFound
	}

	var message core.Message
	err := row.Scan(&message.Id, &message.UserId, &message.ChannelId, &message.Text, &message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) GetAll() ([]*core.Message, error) {
	rows, err := r.Db().Query("SELECT * FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*core.Message
	for rows.Next() {
		var message core.Message
		err := rows.Scan(&message.Id, &message.UserId, &message.ChannelId, &message.Text, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) GetAllBy(field string, value any) ([]*core.Message, error) {
	if !isFieldExist[core.Channel](field) {
		return nil, ErrFieldNotExist
	}

	rows, err := r.Db().Query("SELECT * FROM messages WHERE "+field+" = ?", value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*core.Message
	for rows.Next() {
		var message core.Message
		err := rows.Scan(&message.Id, &message.UserId, &message.ChannelId, &message.Text, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) Create(entity *core.Message) error {
	_, err := r.Db().Exec(
		"INSERT INTO messages (user_id, channel_id, text, created_at) VALUES (?, ?, ?, ?)",
		entity.UserId,
		entity.ChannelId,
		entity.Text,
		entity.CreatedAt,
	)

	return err
}

func (r *MessageRepository) Update(entity *core.Message) error {
	_, err := r.Db().Exec(
		"UPDATE messages SET user_id = ?, channel_id = ?, text = ?, created_at = ? WHERE id = ?",
		entity.UserId,
		entity.ChannelId,
		entity.Text,
		entity.CreatedAt,
		entity.Id,
	)

	return err
}

func (r *MessageRepository) Delete(id int) error {
	_, err := r.Db().Exec("DELETE FROM messages WHERE id = ?", id)

	return err
}
