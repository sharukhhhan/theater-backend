package app

import (
	"github.com/labstack/echo/v4"
	"log"
	"theater/config"
	"theater/internal/controllers/http/v1"
	"theater/internal/repository"
	"theater/internal/service"
	"theater/pkg/hasher"
	"theater/pkg/httpserver"
)

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error occurred while reading config: %s", err.Error())
	}

	db, err := ConnectToPostgres(cfg)
	if err != nil {
		log.Fatalf("error occurred while connecting to db: %s", err.Error())
	}

	repo := repository.NewRepository(db)

	serviceDep := service.Dependency{
		Repo:          repo,
		SignKey:       cfg.JWT.SignKey,
		TokenTTL:      cfg.JWT.TokenTTL,
		AdminUsername: cfg.Admin.Username,
		AdminPassword: cfg.Admin.Password,
		Hasher:        hasher.NewSHA1Hasher(cfg.Salt),
	}

	services := service.NewService(serviceDep)
	handler := echo.New()
	v1.NewRouter(handler, services)

	srv := new(httpserver.Server)
	if err := srv.Run(cfg.HTTP.Port, handler); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
