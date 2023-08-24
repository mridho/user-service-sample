// This file contains types that are used in the repository layer.
package repository

import (
	"errors"
	"time"

	"user-service-sample/generated"
	"user-service-sample/utils/string_helper"
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

func (u *User) UpdateByReq(req generated.UpdateUserJSONRequestBody) bool {
	if u == nil {
		return false
	}

	updated := false

	if reqFullName := string_helper.GetAndTrimPointerStringValue(req.FullName); reqFullName != "" && reqFullName != u.FullName {
		u.FullName = reqFullName
		updated = true
	}

	if reqPhoneNumber := string_helper.GetAndTrimPointerStringValue(req.PhoneNumber); reqPhoneNumber != "" && reqPhoneNumber != u.PhoneNumber {
		u.PhoneNumber = reqPhoneNumber
		updated = true
	}

	return updated
}
