package v1

import "github.com/smolneko-team/neko/internal/model"

type imagesResponse struct {
	Count  int           `json:"count"`
	Images []model.Image `json:"data"`
}
