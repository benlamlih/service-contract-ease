package config

import (
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App struct {
		Port         string `koanf:"port"`
		FrontendURL  string `koanf:"frontend_url"`
		OtelEndpoint string `koanf:"otel_endpoint"`
	} `koanf:"app"`

	DB struct {
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		Name     string `koanf:"name"`
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		Schema   string `koanf:"schema"`
	} `koanf:"db"`
}

var k = koanf.New(".")

func LoadConfig() *Config {
	envFile := "dev"
	if val := os.Getenv("APP_ENV"); val != "" {
		envFile = val
	}

	configFile := "config." + envFile + ".yaml"
	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		log.Fatalf("error loading config.dev.yaml: %v", err)
	}

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		log.Fatalf("error loading env vars: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return &cfg
}
