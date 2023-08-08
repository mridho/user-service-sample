package handler

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/structvalidator"
	"github.com/golang/mock/gomock"
)

type serverMock struct {
	validator  *structvalidator.StructValidator
	repository *repository.MockRepositoryInterface
	cleanUp    func()

	server *Server
}

func setupServerMock(t *testing.T) *serverMock {

	ctrl := gomock.NewController(t)
	repository := repository.NewMockRepositoryInterface(ctrl)

	return &serverMock{
		validator: structvalidator.NewWithOptions(
			structvalidator.WithFieldTag("json"),
			structvalidator.WithCustomTranslation("startswith", "{0} should start with {1}"),
			structvalidator.WithPasswordValidationTag(),
		),
		repository: repository,
		cleanUp: func() {
			t.Helper()
			ctrl.Finish()
		},

		server: NewServer(NewServerOptions{
			Repository: repository,
		}),
	}
}
