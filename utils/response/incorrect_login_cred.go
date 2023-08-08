package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	IncorrectLoginErrorMsg = "phone number or password is incorrect"
)

func IncorrectLoginCred(ctx echo.Context) error {
	return SingleErrorResponse(ctx, http.StatusBadRequest, IncorrectLoginErrorMsg)
}
