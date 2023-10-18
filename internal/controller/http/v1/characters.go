package v1

import (
	"errors"
	"strconv"

	"github.com/smolneko-dev/neko/internal/model"
	"github.com/smolneko-dev/neko/internal/usecase"
	"github.com/smolneko-dev/neko/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type charactersRoutes struct {
	c   usecase.Character
	img usecase.Images
	l   *logger.Logger
}

func newCharactersRoutes(handler fiber.Router, c usecase.Character, img usecase.Images, l *logger.Logger) {
	r := &charactersRoutes{c, img, l}

	h := handler.Group("/characters")
	{
		h.Get("/", r.characters)
		h.Get("/:id", r.character)
		h.Get("/:id/images", r.characterImages)
	}
}

type charactersResponse struct {
	Characters []model.Character `json:"data"`
	PrevCursor string            `json:"previous_cursor,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
}

func (r *charactersRoutes) characters(c *fiber.Ctx) error {
	var count int

	if c.Query("count") == "" {
		count = 20
	} else if value, err := strconv.Atoi(c.Query("count")); err == nil {
		count = value
	} else {
		r.l.Error().Err(err).Msg("http - v1 - characters - count")
		return errorResponse(c, fiber.StatusBadRequest, "Query parameter 'count' is not an integer.")
	}
	if count <= 0 {
		r.l.Error().Msgf("count is negative or zero %d http - v1 - characters - count", count)
		return errorResponse(c, fiber.StatusBadRequest, "Query parameter 'count' is negative or zero.")
	}

	cursor := c.Query("cursor")
	characters, next, prev, err := r.c.Characters(c.UserContext(), count, cursor)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - characters")
		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(charactersResponse{characters, prev, next})
}

type characterResponse struct {
	Character model.Character `json:"data"`
}

func (r *charactersRoutes) character(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err == nil {
		r.l.Error().Err(errors.New("route parameter 'id' is not a string")).Msg("http - v1 - character - id")
		return errorResponse(c, fiber.StatusBadRequest, "Route parameter 'id' is not a valid id.")
	}

	character, err := r.c.Character(c.UserContext(), id)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - character")

		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(characterResponse{character})
}

func (r *charactersRoutes) characterImages(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err == nil {
		r.l.Error().Err(errors.New("route parameter 'id' is not a string")).Msg("http - v1 - characterImages - id")
		return errorResponse(c, fiber.StatusBadRequest, "Route parameter 'id' is not a valid id.")
	}

	preview := c.Query("preview")

	images, err := r.img.Images(c.UserContext(), id, "characters", preview)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - characterImages")
		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	count := len(images)

	return c.Status(fiber.StatusOK).JSON(imagesResponse{count, images})
}
