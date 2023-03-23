package repo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/smolneko-team/neko/config"
	"github.com/smolneko-team/neko/internal/model"
	"github.com/smolneko-team/neko/pkg/postgres"
)

type ImagesRepo struct {
	*postgres.Postgres
	config.Storage
}

func NewImagesRepo(pg *postgres.Postgres, s config.Storage) *ImagesRepo {
	return &ImagesRepo{pg, s}
}

// GetImagesPathByEntityId -
func (r *ImagesRepo) GetImagesPathByEntityId(ctx context.Context, id, entity, preview string) ([]model.Image, error) {
	query := r.Builder.Select("image_path, blurhash, preview")

	var err error
	var isPreview bool
	if preview != "" {
		isPreview, err = strconv.ParseBool(preview)
		if err != nil {
			return nil, fmt.Errorf("ImagesRepo - GetImagesPathByEntityId - ParseBool: %w", err)
		}
	}

	if entity == "figures" {
		query = query.
			From("figures_images").
			Where(sq.Eq{"figure_id": id})
	}
	if entity == "characters" {
		query = query.
			From("characters_images").
			Where(sq.Eq{"character_id": id})
	}

	if isPreview == true {
		query = query.
			Where(sq.Eq{"preview": isPreview}).
			OrderBy("created_at ASC").
			Limit(1)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ImagesRepo - GetImagesPathByEntityId - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ImagesRepo - GetImagesPathByEntityId - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	images := make([]model.Image, 0, 20)

	for rows.Next() {
		image := model.Image{}

		err = rows.Scan(
			&image.URL,
			&image.BlurHash,
			&image.IsPreview,
		)
		if err != nil {
			return nil, fmt.Errorf("ImagesRepo - GetImagesPathByEntityId - rows.Scan: %w", err)
		}
		images = append(images, image)
	}

	if len(images) > 0 {
		images, err = r.signURLs(images)
		if err != nil {
			return nil, fmt.Errorf("ImagesRepo - GetImagesPathByEntityId - SignURLs: %w", err)
		}
	}

	return images, nil
}

// signURLs -
func (r *ImagesRepo) signURLs(images []model.Image) ([]model.Image, error) {
	var keyBin, saltBin []byte
	var err error

	if keyBin, err = hex.DecodeString(r.ImgKey); err != nil {
		return nil, fmt.Errorf("can't decode key - %w", err)
	}

	if saltBin, err = hex.DecodeString(r.ImgSalt); err != nil {
		return nil, fmt.Errorf("can't decode salt - %w", err)
	}

	for i := 0; i < len(images); i++ {
		url := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("s3://%s/%s", r.Bucket, images[i].URL)))

		procOpts := ""
		path := fmt.Sprintf("%s/%s", procOpts, url)

		mac := hmac.New(sha256.New, keyBin)
		mac.Write(saltBin)
		mac.Write([]byte(path))

		signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		images[i].URL = fmt.Sprintf("%s/%s%s", r.ImgURL, signature, path)
	}
	return images, nil
}
