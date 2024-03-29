package app

import (
	"context"
	"fmt"
	"time"

	"github.com/smolneko-dev/neko/config"
	v1 "github.com/smolneko-dev/neko/internal/controller/http/v1"
	"github.com/smolneko-dev/neko/internal/infrastructure/repo"
	"github.com/smolneko-dev/neko/internal/usecase"
	"github.com/smolneko-dev/neko/pkg/httpserver"
	"github.com/smolneko-dev/neko/pkg/logger"
	"github.com/smolneko-dev/neko/pkg/postgres"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level, cfg.App.StageStatus)

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	pg, err := postgres.New(context.Background(), url, 3, 5*time.Second)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Pool.Close()

	figuresUseCase := usecase.NewFigures(repo.NewFiguresRepo(pg, cfg.Storage))
	charactersUseCase := usecase.NewCharacters(repo.NewCharactersRepo(pg))
	imagesUseCase := usecase.NewImages(repo.NewImagesRepo(pg, cfg.Storage))

	routerCfg := v1.RouterConfig{
		Logger:   log,
		CorsUrls: cfg.WebUrls,
	}
	handler := fiber.New(httpserver.FiberConfig(cfg.StageStatus, cfg.App.Name))
	v1.NewRouter(handler, routerCfg, figuresUseCase, charactersUseCase, imagesUseCase)
	httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
}
