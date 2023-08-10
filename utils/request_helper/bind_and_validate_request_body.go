package request_helper

import (
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/SawitProRecruitment/UserService/utils/structvalidator"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidReqBody = errors.New("invalid req body")
)

// BindAndValidateReqBody binds the request body into 'reqPtr'(pointer to a req body struct)
// and then validate the req body struct according to each fields validation tags
func BindAndValidateReqBody[T any](ctx echo.Context, validator *structvalidator.StructValidator, reqPtr *T) error {
	if err := ctx.Bind(reqPtr); err != nil {
		response.SingleErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return err
	}
	if messages, _ := validator.ValidateStruct(reqPtr); len(messages) > 0 {
		response.StandardErrorResponse(ctx, http.StatusBadRequest, messages)
		return ErrInvalidReqBody
	}

	return nil
}
