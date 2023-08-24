package handler

import (
	"net/http"

	"user-service-sample/generated"
	"user-service-sample/repository"
	"user-service-sample/utils/context_helper"
	"user-service-sample/utils/password"
	"user-service-sample/utils/request_helper"
	"user-service-sample/utils/response"

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
	if messages := request_helper.BindAndValidateReqBody(ctx, s.Validator, &req); len(messages) > 0 {
		return response.StandardErrorResponse(ctx, http.StatusBadRequest, messages)
	}

	// check phone number
	if err := s.checkIsPhoneAlreadyRegistered(ctx, tracestr, req.PhoneNumber); err != nil {
		return err
	}

	hashedPassword, salt := password.SaltAndHashPassword(req.Password)

	out, err := s.Repository.InsertUser(ctx.Request().Context(), nil, repository.InsertUserInput{
		PhoneNumber:  req.PhoneNumber,
		FullName:     req.FullName,
		PasswordHash: hashedPassword,
		Salt:         salt,
	})
	if err != nil {
		ctx.Logger().Errorf("%s, failed InsertUser, err: %v", tracestr, err)
		return response.InternalErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusCreated, generated.RegisterResponse{
		Id:      out.Id,
		Message: pointer.String("user registration success"),
	})
}
