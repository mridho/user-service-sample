package response

import (
	"user-service-sample/generated"

	"github.com/labstack/echo/v4"
)

func StandardErrorResponse(ctx echo.Context, code int, messages []string) error {
	return ctx.JSON(code, generated.ErrorResponse{
		Messages: messages,
	})
}

func SingleErrorResponse(ctx echo.Context, code int, message string) error {
	return StandardErrorResponse(ctx, code, []string{message})
}
