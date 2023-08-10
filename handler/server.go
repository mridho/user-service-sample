package handler

import (
	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/structvalidator"
)

type Server struct {
	Validator  *structvalidator.StructValidator
	Config     *config.Config
	Repository repository.RepositoryInterface
}

type NewServerOptions struct {
	Config     *config.Config
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Validator: structvalidator.NewWithOptions(
			structvalidator.WithFieldTag("json"),
			structvalidator.WithCustomTranslation("startswith", "{0} should start with {1}"),
			structvalidator.WithCustomTranslation("required_without_all", "{0} is a required field when {1} not present"),
			structvalidator.WithPasswordValidationTag(),
		),
		Config:     opts.Config,
		Repository: opts.Repository,
	}
}
