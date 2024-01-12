package v1

import "github.com/gofiber/fiber/v2"

type response struct {
	Error string `json:"error"`
}

func errorResponse(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(response{msg})
}
