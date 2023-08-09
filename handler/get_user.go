package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/authentication"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/labstack/echo/v4"
)

// Retrieve user detail
// (GET /v1/user)
func (s *Server) GetUser(ctx echo.Context) error {
	tracestr := "handler.GetUser"
	if err := context_helper.CheckCtxErr(ctx); err != nil {
		return err
	}

	claims, err := authentication.VerifyToken(ctx, s.Config.Secret)
	if err != nil {
		ctx.Logger().Infof("%s, VerifyToken failed, err: %+v", tracestr, err)
		return response.AccessForbidden(ctx)
	}

	user, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		Id: claims.Id,
	})
	if err != nil {
		ctx.Logger().Errorf("%s, failed GetUser by Id err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, generated.UserDataResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	})
}
