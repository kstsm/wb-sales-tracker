package config

import (
	"os"

	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Postgres Postgres
}

type Server struct {
	Host string
	Port int
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	Ssl      string
}

func GetConfig() Config {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		slog.Fatal("Failed to read .env file", "error", err)
		os.Exit(1)
	}

	return Config{
		Server: Server{
			Host: viper.GetString("SRV_HOST"),
			Port: viper.GetInt("SRV_PORT"),
		},
		Postgres: Postgres{
			Username: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			DBName:   viper.GetString("POSTGRES_DB"),
			Ssl:      viper.GetString("POSTGRES_SSL"),
		},
	}
}
