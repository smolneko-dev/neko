package config

import "fmt"

const (
	defaultAddress = ":38575"
)

type Config struct {
	RunAddress string
}

func Load() (Config, error) {
	cfg := Config{
		RunAddress: defaultAddress,
	}

	fmt.Printf("\nstart with config:\n%+v\n\n", cfg)

	return cfg, nil
}
