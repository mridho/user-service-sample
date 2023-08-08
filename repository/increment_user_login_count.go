package repository

import (
	"context"
	"database/sql"
	"time"
)

func (r *Repository) IncrementUserLoginCount(ctx context.Context, tx *sql.Tx, input User) (err error) {

	if input.Id == "" {
		return ErrInvalidInputParam
	}

	updatedAt := time.Now().UTC()

	query := `
		UPDATE users
		SET 
			updated_at = $2,
			login_count = login_count + 1
		WHERE id = $1
	`
	params := []interface{}{
		input.Id,
		updatedAt,
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, params...)
	} else {
		_, err = r.Db.ExecContext(ctx, query, params...)
	}

	return err
}
