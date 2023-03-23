package usecase

import (
	"context"
	"fmt"

	"github.com/smolneko-team/neko/internal/model"
)

type FiguresUseCase struct {
	repo FiguresRepo
}

func NewFigures(r FiguresRepo) *FiguresUseCase {
	return &FiguresUseCase{
		repo: r,
	}
}

func (uc FiguresUseCase) Figures(ctx context.Context, count int, cursor string) ([]model.Figure, string, string, error) {
	figures, next, prev, err := uc.repo.GetFigures(ctx, count, cursor)
	if err != nil {
		return nil, "", "", fmt.Errorf("FiguresUseCase - Figures - f.repo.GetFigures: %w", err)
	}

	return figures, next, prev, nil
}

func (uc FiguresUseCase) Figure(ctx context.Context, id string) (model.Figure, error) {
	figure, err := uc.repo.GetFigureById(ctx, id)
	if err != nil {
		return figure, fmt.Errorf("FiguresUseCase - Figure - f.repo.GetFigure: %w", err)
	}

	return figure, nil
}
