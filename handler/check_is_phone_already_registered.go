package handler

import (
	"database/sql"
	"errors"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/labstack/echo/v4"
)

func (s *Server) checkIsPhoneAlreadyRegistered(ctx echo.Context, tracestr string, phoneNumber string) error {
	uwp, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		PhoneNumber: phoneNumber,
	})
	if err != nil && err != sql.ErrNoRows {
		ctx.Logger().Errorf("%s, failed GetUser by PhoneNumber, err: %v", tracestr, err)
		response.InternalErrorResponse(ctx)
		return err
	}
	if uwp.Id != "" {
		response.PhoneAlreadyRegistered(ctx)
		return errors.New(response.PhoneAlreadyRegisteredErrorMsg)
	}

	return nil
}
