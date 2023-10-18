package v1

import "github.com/smolneko-dev/neko/internal/model"

type imagesResponse struct {
	Count  int           `json:"count"`
	Images []model.Image `json:"data"`
}
