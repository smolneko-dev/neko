package v1

import (
	"fmt"
	"strings"

	_ "github.com/smolneko-dev/neko/docs"
	"github.com/smolneko-dev/neko/internal/usecase"
	"github.com/smolneko-dev/neko/pkg/logger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jaevor/go-nanoid"
	"github.com/rs/zerolog"
)

type RouterConfig struct {
	Logger   *logger.Logger
	CorsUrls string
}

// NewRouter -
// @title smolneko API
// @version 0.1.0
// @description https://smolneko.moe
// @contact.name Create an issue on GitHub.
// @contact.url https://github.com/smolneko-dev/neko/issues/new
// @license.name MIT License
// @license.url https://github.com/smolneko-dev/neko/blob/main/LICENSE
// @BasePath /v1
func NewRouter(app *fiber.App, cfg RouterConfig, f usecase.Figure, c usecase.Character, img usecase.Images) {
	corsCfg := cors.Config{
		Next:         nil,
		AllowOrigins: cfg.CorsUrls,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodHead,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders:     "User-Agent,Origin,Content-Type,Accept,Referrer",
		AllowCredentials: false,
		ExposeHeaders:    "Access-Control-Allow-Origin,Content-Type",
		MaxAge:           3600,
	}

	app.Use(
		recover.New(recover.Config{
			Next:             nil,
			EnableStackTrace: true,
		}),
		requestid.New(
			requestid.Config{
				Next: nil,
				// Get request id from header (X-Request-ID) if exists.
				// Header may be spoofed, must be behind a reverse proxy.
				Header: fiber.HeaderXRequestID,
				Generator: func() string {
					id, err := nanoid.Standard(21)
					if err != nil {
						return ""
					}
					return id()
				},
				ContextKey: "req_id",
			}),
		fiberzerolog.New(fiberzerolog.Config{
			Logger:   cfg.Logger.Logger,
			Fields:   []string{"status", "method", "url", "latency", "ip", "requestId", "error"},
			Messages: []string{"", "", ""},
		}),
		// TODO: refactor this
		func(c *fiber.Ctx) error {
			log := zerolog.Ctx(c.UserContext()).With().Interface("requestId", c.Locals("req_id")).Logger()
			cfg.Logger.Logger = &log
			return c.Next()
		},
		cors.New(corsCfg),
	)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	h := app.Group("/v1")
	{
		newFiguresRoutes(h, f, img, cfg.Logger)
		newCharactersRoutes(h, c, img, cfg.Logger)
	}

	// Not Found (404) error handler
	app.All("*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("API endpoint '%s' does not exist :(", c.OriginalURL())
		return c.Status(fiber.StatusNotFound).JSON(response{msg})
	})
}
