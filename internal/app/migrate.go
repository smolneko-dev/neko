//go:build migrate

package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func init() {
	cfg := DBConfig{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		User:     getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Name:     getEnv("DB_NAME"),
		SSLMode:  getEnv("DB_SSL_MODE"),
	}

	log.Printf("Migrate: started on host %s port %s database %s\n", cfg.Host, cfg.Port, cfg.Name)

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	var attempts = _defaultAttempts
	var err error
	var m *migrate.Migrate
	for attempts > 0 {
		m, err = migrate.New("file://migrations", url)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attemps left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}

func getEnv(key string) string {
	if env, ok := os.LookupEnv(key); ok && len(env) != 0 {
		return env
	} else {
		log.Fatalf("Migrate: environment variable not declared: %s", key)
	}
	return ""
}
