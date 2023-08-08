package repository

import (
	"context"
)

func (r *Repository) GetUser(ctx context.Context, input GetUserInput) (output User, err error) {
	q := `
	SELECT
		id,
		created_at,
		updated_at,
		deleted_at,
		phone_number,
		full_name,
		password_hash,
		salt,
		login_count
	FROM users
	`
	var param interface{}
	if input.Id != "" {
		q += `
		WHERE id = $1
		`
		param = input.Id
	} else if input.PhoneNumber != "" {
		q += `
		WHERE phone_number = $1
		`
		param = input.PhoneNumber
	}

	if param == nil {
		return output, ErrInvalidInputParam
	}

	err = r.Db.QueryRowContext(ctx, q, param).Scan(
		&output.Id,
		&output.CreatedAt,
		&output.UpdatedAt,
		&output.DeletedAt,
		&output.PhoneNumber,
		&output.FullName,
		&output.PasswordHash,
		&output.Salt,
		&output.LoginCount,
	)
	if err != nil {
		return output, err
	}
	return output, nil
}
