package handler

import (
	"database/sql"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/authentication"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/SawitProRecruitment/UserService/utils/string_helper"
	"github.com/labstack/echo/v4"
)

// Update user data
// (PATCH /v1/user)
func (s *Server) UpdateUser(ctx echo.Context) error {
	tracestr := "handler.UpdateUser"
	if err := context_helper.CheckCtxErr(ctx); err != nil {
		return err
	}

	claims, err := authentication.VerifyToken(ctx, s.Config.Secret)
	if err != nil {
		ctx.Logger().Infof("%s, VerifyToken failed, err: %+v", tracestr, err)
		return response.AccessForbidden(ctx)
	}

	var req generated.UpdateUserJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return response.SingleErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	if messages, _ := s.Validator.ValidateStruct(req); len(messages) > 0 {
		return response.StandardErrorResponse(ctx, http.StatusBadRequest, messages)
	}
	if req.FullName == nil && req.PhoneNumber == nil {
		return response.SingleErrorResponse(ctx, http.StatusBadRequest, "request need to have either fullName or phoneNumber")
	}

	// get current user data
	user, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		Id: claims.Id,
	})
	if err != nil {
		ctx.Logger().Errorf("%s, failed GetUser by Id, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	// check phone number
	reqPhoneNumber := string_helper.GetAndTrimPointerStringValue(req.PhoneNumber)
	if reqPhoneNumber != "" && reqPhoneNumber != user.PhoneNumber {
		uwp, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
			PhoneNumber: reqPhoneNumber,
		})
		if err != nil && err != sql.ErrNoRows {
			ctx.Logger().Errorf("%s, failed GetUser by PhoneNumber, err: %v", tracestr, err)
			return response.InternalErrorResponse(ctx)
		}
		if uwp.Id != "" {
			return response.PhoneAlreadyRegistered(ctx)
		}
	}

	// update current user data
	user.UpdateByReq(req)
	if err := s.Repository.UpdateUser(ctx.Request().Context(), nil, user); err != nil {
		ctx.Logger().Errorf("%s, failed UpdateUser, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, generated.UserDataResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	})
}
