package handler

import (
	"database/sql"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/authentication"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/password"
	"github.com/SawitProRecruitment/UserService/utils/request_helper"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/labstack/echo/v4"
)

// Log In as registered user, will return JWT
// (POST /v1/login)
func (s *Server) Login(ctx echo.Context) error {
	tracestr := "handler.Login"
	if err := context_helper.CheckCtxErr(ctx); err != nil {
		return err
	}

	var req generated.LoginJSONRequestBody
	if err := request_helper.BindAndValidateReqBody(ctx, s.Validator, &req); err != nil {
		return err
	}

	user, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.IncorrectLoginCred(ctx)
		}
		ctx.Logger().Errorf("%s, failed GetUser by PhoneNumber, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	if !password.CheckPassword(req.Password, user.PasswordHash, user.Salt) {
		return response.IncorrectLoginCred(ctx)
	}

	// user & password correct
	token, err := authentication.GenerateSignedToken(s.Config.Secret, user)
	if err != nil {
		ctx.Logger().Errorf("%s, failed GenerateSignedToken, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	// increment login count
	if err := s.Repository.IncrementUserLoginCount(ctx.Request().Context(), nil, user); err != nil {
		ctx.Logger().Infof("%s, failed IncrementUserLoginCount, err: %v", tracestr, err)
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Id:    user.Id,
		Token: token,
	})
}
