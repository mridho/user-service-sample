package handler

import (
	"net/http"

	"user-service-sample/generated"
	"user-service-sample/repository"
	"user-service-sample/utils/authentication"
	"user-service-sample/utils/context_helper"
	"user-service-sample/utils/request_helper"
	"user-service-sample/utils/response"
	"user-service-sample/utils/string_helper"

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
	if messages := request_helper.BindAndValidateReqBody(ctx, s.Validator, &req); len(messages) > 0 {
		return response.StandardErrorResponse(ctx, http.StatusBadRequest, messages)
	}

	// get current user data
	user, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		Id: claims.Id,
	})
	if err != nil {
		ctx.Logger().Errorf("%s, failed GetUser by Id, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	// check phone number if req not empty & not the same with user current phone number
	reqPhoneNumber := string_helper.GetAndTrimPointerStringValue(req.PhoneNumber)
	if reqPhoneNumber != "" && reqPhoneNumber != user.PhoneNumber {
		if err := s.checkIsPhoneAlreadyRegistered(ctx, tracestr, reqPhoneNumber); err != nil {
			return err
		}
	}

	// update current user data
	if user.UpdateByReq(req) {
		if err := s.Repository.UpdateUser(ctx.Request().Context(), nil, user); err != nil {
			ctx.Logger().Errorf("%s, failed UpdateUser, err: %v", tracestr, err)
			return response.InternalErrorResponse(ctx)
		}
	}

	return ctx.JSON(http.StatusOK, generated.UserDataResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	})
}
