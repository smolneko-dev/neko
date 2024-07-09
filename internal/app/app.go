package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/smolneko-dev/neko/internal/config"
	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
	cfg config.Config
}

func New(_ context.Context, cfg config.Config) (App, error) {
	return App{
		log: zap.Must(zap.NewProduction()),
		cfg: cfg,
	}, nil
}

func (a App) Run(_ context.Context) error {
	defer func() {
		err := a.log.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Println(err)
		}
	}()

	r := fiber.New()

	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	a.log.Info("neko started", zap.String("addr", a.cfg.RunAddress))

	err := r.Listen(a.cfg.RunAddress)
	if err != nil {
		return fmt.Errorf("http router: %w", err)
	}

	return nil
}
