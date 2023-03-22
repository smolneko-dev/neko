package v1

import (
	"errors"
	"strconv"

	"github.com/smolneko-team/smolneko/internal/model"
	"github.com/smolneko-team/smolneko/internal/usecase"
	"github.com/smolneko-team/smolneko/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type figuresRoutes struct {
	f   usecase.Figure
	img usecase.Images
	l   *logger.Logger
}

func newFiguresRoutes(handler fiber.Router, f usecase.Figure, img usecase.Images, l *logger.Logger) {
	r := &figuresRoutes{f, img, l}

	h := handler.Group("/figures")
	{
		h.Get("/", r.figures)
		h.Get("/:id", r.figure)
		h.Get("/:id/images", r.figureImages)
	}
}

type figuresResponse struct {
	Figures    []model.Figure `json:"data"`
	PrevCursor string         `json:"previous_cursor,omitempty"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

func (r *figuresRoutes) figures(c *fiber.Ctx) error {
	var count int

	if c.Query("count") == "" {
		count = 20
	} else if value, err := strconv.Atoi(c.Query("count")); err == nil {
		count = value
	} else {
		r.l.Error().Err(err).Msg("http - v1 - figures - count")

		return errorResponse(c, fiber.StatusBadRequest, "Query parameter 'count' is not an integer.")
	}
	if count <= 0 {
		r.l.Error().Msgf("count is negative or zero %d http - v1 - figures - count", count)
		return errorResponse(c, fiber.StatusBadRequest, "Query parameter 'count' is negative or zero.")
	}

	cursor := c.Query("cursor")
	figures, next, prev, err := r.f.Figures(c.UserContext(), count, cursor)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - figures")
		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(figuresResponse{figures, prev, next})
}

type figureResponse struct {
	Figure model.Figure `json:"data"`
}

func (r *figuresRoutes) figure(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err == nil || len(id) != 21 {
		r.l.Error().Err(errors.New("route parameter 'id' is not a nanoid(21)")).Msg("http - v1 - figure - id")
		return errorResponse(c, fiber.StatusBadRequest, "Route parameter 'id' is not a valid id.")
	}

	figure, err := r.f.Figure(c.UserContext(), id)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - figure")
		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(figureResponse{figure})
}

func (r *figuresRoutes) figureImages(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err == nil {
		r.l.Error().Err(errors.New("route parameter 'id' is not a string")).Msg("http - v1 - figureImages - id")
		return errorResponse(c, fiber.StatusBadRequest, "Route parameter 'id' is not a valid id.")
	}

	preview := c.Query("preview")
	images, err := r.img.Images(c.UserContext(), id, "figures", preview)
	if err != nil {
		r.l.Error().Err(err).Msg("http - v1 - figureImages")
		return errorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}
	count := len(images)

	return c.Status(fiber.StatusOK).JSON(imagesResponse{count, images})
}
