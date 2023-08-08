// This file contains types that are used in the repository layer.
package repository

import (
	"errors"
	"time"
)

var (
	ErrInvalidInputParam = errors.New("invalid input param")
)

type InsertUserInput struct {
	PhoneNumber  string
	FullName     string
	PasswordHash string
	Salt         string
}

type InsertUserOutput struct {
	Id string
}

type GetUserInput struct {
	Id          string
	PhoneNumber string
}

type User struct {
	Id           string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	PhoneNumber  string
	FullName     string
	PasswordHash string
	Salt         string
	LoginCount   uint32
}
