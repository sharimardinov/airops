package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env string

	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	DSN             string
	MaxConns        int32
	ConnMaxLifetime time.Duration
}

func MustLoad() *Config {
	cfg := &Config{
		Env: getEnv("APP_ENV", "local"),

		HTTP: HTTPConfig{
			Addr:         getEnv("HTTP_ADDR", ":8080"),
			ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		},

		DB: DBConfig{
			DSN:             mustEnv("DB_DSN"),
			MaxConns:        int32(getInt("DB_MAX_CONNS", 10)),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		},
	}

	return cfg
}

// helpers

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return v
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
