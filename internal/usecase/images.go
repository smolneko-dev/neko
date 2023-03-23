package usecase

import (
	"context"
	"fmt"

	"github.com/smolneko-team/neko/internal/model"
)

type ImagesUseCase struct {
	repo ImagesRepo
}

func NewImages(r ImagesRepo) *ImagesUseCase {
	return &ImagesUseCase{
		repo: r,
	}
}

func (uc ImagesUseCase) Images(ctx context.Context, id, entity, preview string) ([]model.Image, error) {
	images, err := uc.repo.GetImagesPathByEntityId(ctx, id, entity, preview)
	if err != nil {
		return images, fmt.Errorf("ImagesUseCase - Image - uc.repo.GetImagesPathByEntityId: %w", err)
	}

	return images, nil
}
