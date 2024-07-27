package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultAddress = ":38575"
)

type Config struct {
	Server
}

type Server struct {
	Address string `mapstructure:"SERVER_ADDRESS"`
}

func Load() (Config, error) {
	viper.SetEnvPrefix("SN")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("config file not found, environment variables and default values will be used")
		} else {
			return Config{}, fmt.Errorf("can't read config file: %w", err)
		}
	} else {
		log.Println("config file found and successfully parsed")
	}

	viper.SetDefault("server.address", defaultAddress)

	cfg := Config{
		Server: Server{
			Address: viper.GetString("server.address"),
		},
	}

	log.Printf("start with config: %+v\n\n ", cfg)

	return cfg, nil
}
