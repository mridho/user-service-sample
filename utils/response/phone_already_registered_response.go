package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	PhoneAlreadyRegisteredErrorMsg = "phone number already registered"
)

func PhoneAlreadyRegistered(ctx echo.Context) error {
	return SingleErrorResponse(ctx, http.StatusConflict, PhoneAlreadyRegisteredErrorMsg)
}
