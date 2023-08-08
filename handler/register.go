package handler

import (
	"database/sql"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/password"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/labstack/echo/v4"
	"github.com/xorcare/pointer"
)

// Register new user to service
// (POST /v1/register)
func (s *Server) Register(ctx echo.Context) error {
	tracestr := "handler.Register"
	if err := context_helper.CheckCtxErr(ctx); err != nil {
		return err
	}

	var req generated.RegisterJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return response.SingleErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	if messages, _ := s.validator.ValidateStruct(req); len(messages) > 0 {
		return response.StandardErrorResponse(ctx, http.StatusBadRequest, messages)
	}

	// check phone number
	user, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil && err != sql.ErrNoRows {
		ctx.Logger().Errorf("%s, failed GetUser by PhoneNumber, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}
	if user.Id != "" {
		return response.PhoneAlreadyRegistered(ctx)
	}

	hashedPassword, salt := password.SaltAndHashPassword(req.Password)

	out, err := s.Repository.InsertUser(ctx.Request().Context(), nil, repository.InsertUserInput{
		PhoneNumber:  req.PhoneNumber,
		FullName:     req.FullName,
		PasswordHash: hashedPassword,
		Salt:         salt,
	})
	if err != nil {
		ctx.Logger().Errorf("%s, failed InsertNewUser, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusCreated, generated.RegisterResponse{
		Id:      out.Id,
		Message: pointer.String("user registration success"),
	})
}
