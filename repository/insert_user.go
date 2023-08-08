package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (r *Repository) InsertUser(ctx context.Context, tx *sql.Tx, input InsertUserInput) (output InsertUserOutput, err error) {
	id := uuid.NewString()
	query := `
		INSERT INTO users (id, phone_number, full_name, password_hash, salt)
		VALUES ( $1, $2, $3, $4, $5)
	`
	params := []interface{}{
		id,
		input.PhoneNumber,
		input.FullName,
		input.PasswordHash,
		input.Salt,
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, params...)
	} else {
		_, err = r.Db.ExecContext(ctx, query, params...)
	}
	if err != nil {
		return InsertUserOutput{}, err
	}

	return InsertUserOutput{
		Id: id,
	}, nil
}
