package request_helper

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/utils/structvalidator"
	"github.com/labstack/echo/v4"
)

type requestBody interface {
	generated.LoginJSONRequestBody |
		generated.RegisterJSONRequestBody |
		generated.UpdateUserJSONRequestBody
}

// BindAndValidateReqBody binds the request body into 'reqPtr'(pointer to a req body struct)
// and then validate the req body struct according to each fields validation tags
func BindAndValidateReqBody[T requestBody](ctx echo.Context, validator *structvalidator.StructValidator, reqPtr *T) []string {
	if err := ctx.Bind(reqPtr); err != nil {
		return []string{err.Error()}
	}
	if messages, _ := validator.ValidateStruct(reqPtr); len(messages) > 0 {
		return messages
	}

	return nil
}
