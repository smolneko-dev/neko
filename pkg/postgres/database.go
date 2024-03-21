package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(ctx context.Context, url string, attempts int, timeout time.Duration) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create new connection pool: %w", err)
	}

	err = retryWithAttempts(func() error {
		if err = pool.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping database: %w", err)
		}
		return nil
	}, attempts, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &DB{
		Pool:    pool,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func retryWithAttempts(fn func() error, attempts int, timeout time.Duration) error {
	var err error

	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(timeout)
			attempts--
			continue
		}
		return nil
	}

	return err
}
