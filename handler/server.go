package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/structvalidator"
)

type Server struct {
	validator  *structvalidator.StructValidator
	Repository repository.RepositoryInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		validator: structvalidator.NewWithOptions(
			structvalidator.WithFieldTag("json"),
			structvalidator.WithCustomTranslation("startswith", "{0} should start with {1}"),
			structvalidator.WithPasswordValidationTag(),
		),
		Repository: opts.Repository,
	}
}
