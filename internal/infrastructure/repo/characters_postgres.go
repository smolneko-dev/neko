package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/smolneko-dev/neko/internal/model"
	"github.com/smolneko-dev/neko/pkg/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharactersRepo struct {
	*postgres.Postgres
}

func NewCharactersRepo(pg *postgres.Postgres) *CharactersRepo {
	return &CharactersRepo{pg}
}

func (r *CharactersRepo) GetCharacterById(ctx context.Context, id string) (model.Character, error) {
	character := model.Character{}

	lang := "en"

	sql, args, err := r.Builder.
		Select("id, birth_at, created_at, updated_at, is_draft").
		Column("COALESCE(name ->  ?  #>>'{}', name -> ?  #>>'{}')", lang, "en").
		Column("COALESCE(description -> ? #>>'{}', description -> ? #>>'{}') as description", lang, "en").
		From("characters").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return character, fmt.Errorf("CharactersRepo - GetCharacterById - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)
	if err != nil {
		return character, fmt.Errorf("CharactersRepo - GetCharacterById - r.Pool.QueryRow: %w", err)
	}

	err = row.Scan(
		&character.ID,
		&character.BirthAt,
		&character.CreatedAt,
		&character.UpdatedAt,
		&character.IsDraft,
		&character.Name,
		&character.Description,
	)
	if err != nil {
		return character, fmt.Errorf("CharactersRepo - GetCharacterById - row.Scan: %w", err)
	}

	return character, nil
}

func (r *CharactersRepo) GetCharacters(ctx context.Context, count int, cursor string) ([]model.Character, string, string, error) {
	if count > 50 {
		count = 50
	}

	lang := "en"
	columns := r.Builder.
		Select("id, birth_at, created_at, updated_at, is_draft").
		Column("COALESCE(name -> ? #>>'{}', name -> ? #>>'{}') as name", lang, "en").
		Column("COALESCE(description -> ? #>>'{}', description -> ? #>>'{}') as description", lang, "en").
		From("characters")

	query := r.Builder.Select("characters_cols.*")

	var created time.Time
	var id, suffix string
	var err error
	if cursor != "" {
		created, id, suffix, err = decodeCursor(cursor)
		if err != nil {
			return nil, "", "", fmt.Errorf("CharactersRepo - GetCharacters - decodeCursor : %w", err)
		}

		if suffix == "next" {
			columns = columns.
				Where(sq.LtOrEq{
					"created_at": created,
				}).Where(sq.Or{
				sq.Lt{
					"created_at": created,
				},
				sq.Lt{
					"id": id,
				},
			}).OrderBy("created_at DESC, id DESC")
		} else {
			columns = columns.
				Where(sq.GtOrEq{
					"created_at": created,
				}).Where(sq.Or{
				sq.Gt{
					"created_at": created,
				},
				sq.Gt{
					"id": id,
				},
			}).OrderBy("created_at ASC, id ASC")
			query = query.OrderBy("created_at DESC, id DESC")
		}
	} else {
		columns = columns.OrderBy("created_at DESC, id DESC")
	}
	columns = columns.Limit(uint64(count) + 1)

	query = query.FromSelect(columns, "characters_cols").
		OrderBy("created_at DESC, id DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, "", "", fmt.Errorf("CharactersRepo - GetCharacters - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", "", fmt.Errorf("CharactersRepo - GetCharacters - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	characters := make([]model.Character, 0, _defaultEntityCap)

	for rows.Next() {
		character := model.Character{}
		var date pgtype.Date

		err = rows.Scan(
			&character.ID,
			&date,
			&character.CreatedAt,
			&character.UpdatedAt,
			&character.IsDraft,
			&character.Name,
			&character.Description,
		)
		if err != nil {
			return nil, "", "", fmt.Errorf("CharactersRepo - GetCharacters - rows.Scan: %w", err)
		}
		if date.Time.IsZero() {
			character.BirthAt = ""
		} else {
			character.BirthAt = date.Time.Format("2006-01-02")
		}

		characters = append(characters, character)
	}

	hasNextPage := false
	if len(characters) >= count+1 {
		hasNextPage = true
	}

	if hasNextPage == true {
		if suffix == "next" || suffix == "" {
			characters = characters[:len(characters)-1]
		} else {
			characters = characters[1:]
		}
	}

	var previousCursor, nextCursor string
	if (len(characters) > 0 && hasNextPage == true) || (len(characters) > 0 && suffix == "prev") {
		nextCursor = encodeCursor(characters[len(characters)-1].CreatedAt, characters[len(characters)-1].ID, "next")
	}
	if ((suffix == "prev" || cursor != "") && len(characters) > 0 && hasNextPage == true) ||
		(suffix == "next" && len(characters) > 0 && len(characters) < count) {
		previousCursor = encodeCursor(characters[0].CreatedAt, characters[0].ID, "prev")
	}

	return characters, nextCursor, previousCursor, nil
}
