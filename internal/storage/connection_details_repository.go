package storage

import "github.com/hop-/gotchat/internal/core"

type ConnectionDetailsRepository struct {
	StorageDb
}

func newConnectionDetailsRepository(storage StorageDb) *ConnectionDetailsRepository {
	return &ConnectionDetailsRepository{
		StorageDb: storage,
	}
}

func (r *ConnectionDetailsRepository) GetOne(id int) (*core.ConnectionDetails, error) {
	row := r.Db().QueryRow("SELECT id, host_unique_id, client_unique_id, encryption_key, decryption_key, key_derivation_salt, created_at FROM connection_details WHERE id = ?", id)
	if row == nil {
		return nil, ErrNotFound
	}

	var details core.ConnectionDetails
	err := row.Scan(&details.Id, &details.HostUniqueId, &details.ClientUniqueId, &details.EncryptionKey, &details.DecryptionKey, &details.KeyDerivationSalt, &details.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *ConnectionDetailsRepository) GetOneBy(field string, value any) (*core.ConnectionDetails, error) {
	if !isFieldExist[core.ConnectionDetails](field) {
		return nil, ErrFieldNotExist
	}
	row := r.Db().QueryRow("SELECT id, host_unique_id, client_unique_id, encryption_key, decryption_key, key_derivation_salt, created_at FROM connection_details WHERE "+field+" = ?", value)
	if row == nil {
		return nil, ErrNotFound
	}

	var details core.ConnectionDetails
	err := row.Scan(&details.Id, &details.HostUniqueId, &details.ClientUniqueId, &details.EncryptionKey, &details.DecryptionKey, &details.KeyDerivationSalt, &details.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *ConnectionDetailsRepository) GetAll() ([]*core.ConnectionDetails, error) {
	rows, err := r.Db().Query("SELECT id, host_unique_id, client_unique_id, encryption_key, decryption_key, key_derivation_salt, created_at FROM connection_details")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detailsList []*core.ConnectionDetails
	for rows.Next() {
		var details core.ConnectionDetails
		err := rows.Scan(&details.Id, &details.HostUniqueId, &details.ClientUniqueId, &details.EncryptionKey, &details.DecryptionKey, &details.KeyDerivationSalt, &details.CreatedAt)
		if err != nil {
			return nil, err
		}
		detailsList = append(detailsList, &details)
	}

	return detailsList, nil
}

func (r *ConnectionDetailsRepository) GetAllBy(field string, value any) ([]*core.ConnectionDetails, error) {
	if !isFieldExist[core.ConnectionDetails](field) {
		return nil, ErrFieldNotExist
	}

	rows, err := r.Db().Query("SELECT id, host_unique_id, client_unique_id, encryption_key, decryption_key, key_derivation_salt, created_at FROM connection_details WHERE "+field+" = ?", value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detailsList []*core.ConnectionDetails
	for rows.Next() {
		var details core.ConnectionDetails
		err := rows.Scan(&details.Id, &details.HostUniqueId, &details.ClientUniqueId, &details.EncryptionKey, &details.DecryptionKey, &details.KeyDerivationSalt, &details.CreatedAt)
		if err != nil {
			return nil, err
		}
		detailsList = append(detailsList, &details)
	}

	return detailsList, nil
}

func (r *ConnectionDetailsRepository) Create(details *core.ConnectionDetails) error {
	_, err := r.Db().Exec(
		"INSERT INTO connection_details (host_unique_id, client_unique_id, encryption_key, decryption_key, key_derivation_salt, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		details.HostUniqueId,
		details.ClientUniqueId,
		details.EncryptionKey,
		details.DecryptionKey,
		details.KeyDerivationSalt,
		details.CreatedAt,
	)

	return err
}

func (r *ConnectionDetailsRepository) Update(details *core.ConnectionDetails) error {
	_, err := r.Db().Exec(
		"UPDATE connection_details SET host_unique_id = ?, client_unique_id = ?, encryption_key = ?, decryption_key = ?, key_derivation_salt = ?, created_at = ? WHERE id = ?",
		details.HostUniqueId,
		details.EncryptionKey,
		details.DecryptionKey,
		details.KeyDerivationSalt,
		details.CreatedAt,
		details.Id,
	)

	return err
}

func (r *ConnectionDetailsRepository) Delete(id int) error {
	_, err := r.Db().Exec("DELETE FROM connection_details WHERE id = ?", id)

	return err
}
