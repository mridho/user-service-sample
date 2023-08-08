package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	InternalServerErrorMsg = "internal server error"
)

func InternalErrorResponse(ctx echo.Context) error {
	return SingleErrorResponse(ctx, http.StatusInternalServerError, InternalServerErrorMsg)
}
