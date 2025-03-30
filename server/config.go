package server

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func defaultConfig() Config {
	return Config{
		Port:          "3000",
		LastFMAPIKey:  getLastFMAPIKey(),
		RateLimit:     50,
		RateLimitTime: "1m",
	}
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		slog.Error("Failed to open config file", slog.Any("err", err))
		slog.Info("Using default config")
		return defaultConfig()
	}

	var cfg Config
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		slog.Error("Failed to decode config file", slog.Any("err", err))
		slog.Info("Using default config")
		return defaultConfig()
	}
	return cfg
}

type Config struct {
	Port          string `yaml:"port"`
	LastFMAPIKey  string `yaml:"lastfm_api_key"`
	RateLimit     int    `yaml:"rate_limit"`
	RateLimitTime string `yaml:"rate_limit_time"`
}

func getLastFMAPIKey() string {
	envValue := os.Getenv("LASTFM_API_KEY")

	if strings.HasPrefix(envValue, "/") {
		data, err := os.ReadFile(envValue)
		if err != nil {
			log.Fatalf("Failed to read LastFM API key from secret file: %v", err)
		}
		return strings.TrimSpace(string(data))
	}

	if envValue == "" {
		log.Fatal("LASTFM_API_KEY is not set")
	}

	return envValue
}
