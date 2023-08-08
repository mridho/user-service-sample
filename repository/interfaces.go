// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
	"database/sql"
)

type RepositoryInterface interface {
	GetUser(ctx context.Context, input GetUserInput) (output User, err error)
	InsertUser(ctx context.Context, tx *sql.Tx, input InsertUserInput) (output InsertUserOutput, err error)
	IncrementUserLoginCount(ctx context.Context, tx *sql.Tx, input User) (err error)
	UpdateUser(ctx context.Context, tx *sql.Tx, input User) (err error)
}
