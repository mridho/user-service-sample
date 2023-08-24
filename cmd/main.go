package main

import (
	"os"
	"user-service-sample/config"
	"user-service-sample/generated"
	"user-service-sample/handler"
	"user-service-sample/repository"

	"github.com/labstack/echo/v4"
)

var (
	e      *echo.Echo
	cfg    *config.Config
	err    error
	server generated.ServerInterface
)

func init() {
	e = echo.New()

	cfgFile := "/opt/config.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		cfgFile = "config.yml"
	}
	cfg, err = config.NewConfig(cfgFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	server = newServer(cfg)

	generated.RegisterHandlers(e, server)
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

func main() {
	if err := e.Start(":1323"); err != nil {
		e.Logger.Fatal(err)
	}
}
