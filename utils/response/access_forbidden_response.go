package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	AccessForbiddenErrorMsg = "unauthorized"
)

func AccessForbidden(ctx echo.Context) error {
	return SingleErrorResponse(ctx, http.StatusForbidden, AccessForbiddenErrorMsg)
}
