package app

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const addr = ":38575"

func Run() {
	log := zap.Must(zap.NewProduction())
	defer log.Sync()

	r := fiber.New()

	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Info("neko started", zap.String("addr", addr))

	err := r.Listen(addr)
	if err != nil {
		panic(err)
	}
}
