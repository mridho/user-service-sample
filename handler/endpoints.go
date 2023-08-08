package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	out, err := s.Repository.GetTestById(ctx.Request().Context(), repository.GetTestByIdInput{
		Id: params.Id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: "user not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %s", out.Name)

	return ctx.JSON(http.StatusOK, resp)
}
