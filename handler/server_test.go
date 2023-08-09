package handler

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/test_helper"
	"github.com/golang/mock/gomock"
)

type serverMock struct {
	config     *config.Config
	repository *repository.MockRepositoryInterface
	cleanUp    func()

	server *Server
}

func setupServerMock(t *testing.T) *serverMock {

	ctrl := gomock.NewController(t)
	repository := repository.NewMockRepositoryInterface(ctrl)

	mockConfig := &config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Database: "user_service_db",
		},
		Secret: config.SecretConfig{
			RsaPrivatePem: test_helper.TestRsaPrivatePem,
			RsaPublicPem:  test_helper.TestRsaPublicPem,
		},
	}

	return &serverMock{
		config:     mockConfig,
		repository: repository,
		cleanUp: func() {
			t.Helper()
			ctrl.Finish()
		},

		server: NewServer(NewServerOptions{
			Config:     mockConfig,
			Repository: repository,
		}),
	}
}
