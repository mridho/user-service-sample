package main

import (
	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	cfg, err := config.NewConfig("config.yml")
	if err != nil {
		e.Logger.Fatal(err)
	}
	var server generated.ServerInterface = newServer(cfg)

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer(cfg *config.Config) *handler.Server {
	// dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: cfg.DB.ToJdbcUrl(),
	})
	opts := handler.NewServerOptions{
		Config:     cfg,
		Repository: repo,
	}
	return handler.NewServer(opts)
}
