package context_helper

import (
	"context"
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/labstack/echo/v4"
)

var (
	ErrRequestCanceled  = errors.New("request canceled")
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)

// check ctx.Err after ctx.Done channel is closed
func CheckCtxErr(ctx echo.Context) (err error) {
	select {
	case <-ctx.Request().Context().Done():
		switch ctx.Request().Context().Err() {
		case context.Canceled:
			response.SingleErrorResponse(ctx, http.StatusGone, ErrRequestCanceled.Error())
			return ErrRequestCanceled
		case context.DeadlineExceeded:
			response.SingleErrorResponse(ctx, http.StatusRequestTimeout, ErrDeadlineExceeded.Error())
			return ErrDeadlineExceeded
		default:
			return
		}
	default:
		return
	}
}
