package config

import (
	"os"
)

type Config struct {
	DB struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     string
		Schema   string
	}
	App struct {
		Port string
	}
}

func LoadConfig() *Config {
	cfg := &Config{}
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Schema = os.Getenv("DB_SCHEMA")
	cfg.App.Port = os.Getenv("PORT")

	return cfg
}
