package app

import (
	"fmt"
	"time"

	"github.com/smolneko-team/smolneko/config"
	v1 "github.com/smolneko-team/smolneko/internal/controller/http/v1"
	"github.com/smolneko-team/smolneko/internal/infrastructure/repo"
	"github.com/smolneko-team/smolneko/internal/usecase"
	"github.com/smolneko-team/smolneko/pkg/httpserver"
	"github.com/smolneko-team/smolneko/pkg/logger"
	"github.com/smolneko-team/smolneko/pkg/postgres"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level, cfg.App.StageStatus)

	url := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	pg, err := postgres.New(url,
		postgres.MaxPoolSize(cfg.DB.PoolMax),
		postgres.ConnAttempts(10),
		postgres.ConnTimeout(5*time.Second),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	figuresUseCase := usecase.NewFigures(repo.NewFiguresRepo(pg, cfg.Storage))
	charactersUseCase := usecase.NewCharacters(repo.NewCharactersRepo(pg))
	imagesUseCase := usecase.NewImages(repo.NewImagesRepo(pg, cfg.Storage))

	routerCfg := v1.RouterConfig{
		Logger:   l,
		CorsUrls: cfg.WebUrls,
	}

	handler := fiber.New(httpserver.FiberConfig(cfg.StageStatus, cfg.App.Name))
	v1.NewRouter(handler, routerCfg, figuresUseCase, charactersUseCase, imagesUseCase)
	httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
}
