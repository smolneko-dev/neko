package v1

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/smolneko-team/smolneko/internal/usecase"
	"github.com/smolneko-team/smolneko/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func NewRouter(app *fiber.App, webUrls string, l logger.Interface, f usecase.Figure, c usecase.Character, img usecase.Images) {
	log.Println(webUrls)

	corsCfg := cors.Config{
		Next:         nil,
		AllowOrigins: webUrls,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodHead,
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
		cors.New(corsCfg),
		fLogger.New(fLogger.ConfigDefault),
	)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	h := app.Group("/v1")
	{
		newFiguresRoutes(h, f, img, l)
		newCharactersRoutes(h, c, img, l)
	}

	// Not Found (404) error handler
	app.All("*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("API endpoint '%s' does not exist :(", c.OriginalURL())
		return c.Status(fiber.StatusNotFound).JSON(response{msg})
	})
}
