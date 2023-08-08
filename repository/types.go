// This file contains types that are used in the repository layer.
package repository

import (
	"errors"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/utils/string_helper"
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

func (u *User) UpdateByReq(req generated.UpdateUserJSONRequestBody) {
	if u == nil {
		return
	}

	if reqFullName := string_helper.GetAndTrimPointerStringValue(req.FullName); reqFullName != "" {
		u.FullName = reqFullName
	}

	if reqPhoneNumber := string_helper.GetAndTrimPointerStringValue(req.PhoneNumber); reqPhoneNumber != "" {
		u.PhoneNumber = reqPhoneNumber
	}
}
