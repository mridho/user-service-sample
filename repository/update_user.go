package repository

import (
	"context"
	"database/sql"
	"time"
)

func (r *Repository) UpdateUser(ctx context.Context, tx *sql.Tx, input User) (err error) {

	if input.Id == "" {
		return ErrInvalidInputParam
	}

	updatedAt := time.Now().UTC()

	query := `
		UPDATE users
		SET 
			updated_at = $2,
			phone_number = $3,
			full_name = $4
		WHERE id = $1
	`
	params := []interface{}{
		input.Id,
		updatedAt,
		input.PhoneNumber,
		input.FullName,
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, params...)
	} else {
		_, err = r.Db.ExecContext(ctx, query, params...)
	}

	return err
}
