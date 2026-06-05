package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from the environment.
type Config struct {
	Env       string
	LogLevel  string
	Host      string
	Port      string
	PublicURL string

	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Upload   UploadConfig
	CORS     CORSConfig
	Rate     RateConfig
	Market   MarketplaceContacts
	Admin    AdminBootstrap
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type UploadConfig struct {
	Dir        string
	MaxSizeMB  int64
	PublicBase string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type RateConfig struct {
	RPS   int
	Burst int
}

type MarketplaceContacts struct {
	Phone    string
	WhatsApp string
}

type AdminBootstrap struct {
	Email    string
	Password string
	Name     string
}

// Load reads configuration from environment variables (and a .env file if present).
func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error if .env is missing

	cfg := &Config{
		Env:       getEnv("API_ENV", "development"),
		LogLevel:  getEnv("API_LOG_LEVEL", "info"),
		Host:      getEnv("API_HOST", "0.0.0.0"),
		Port:      getEnv("API_PORT", "8080"),
		PublicURL: getEnv("PUBLIC_URL", "http://localhost:3000"),
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "alistroy"),
			Password: getEnv("POSTGRES_PASSWORD", "alistroy"),
			DBName:   getEnv("POSTGRES_DB", "alistroy"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			AccessTTL:  time.Duration(getEnvInt("JWT_ACCESS_TTL_MINUTES", 60)) * time.Minute,
			RefreshTTL: time.Duration(getEnvInt("JWT_REFRESH_TTL_HOURS", 720)) * time.Hour,
		},
		Upload: UploadConfig{
			Dir:        getEnv("UPLOAD_DIR", "./uploads"),
			MaxSizeMB:  int64(getEnvInt("UPLOAD_MAX_SIZE_MB", 10)),
			PublicBase: getEnv("UPLOAD_PUBLIC_BASE", "/uploads"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitAndTrim(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		},
		Rate: RateConfig{
			RPS:   getEnvInt("RATE_LIMIT_RPS", 20),
			Burst: getEnvInt("RATE_LIMIT_BURST", 40),
		},
		Market: MarketplaceContacts{
			Phone:    getEnv("MARKETPLACE_PHONE", ""),
			WhatsApp: getEnv("MARKETPLACE_WHATSAPP", ""),
		},
		Admin: AdminBootstrap{
			Email:    getEnv("ADMIN_EMAIL", ""),
			Password: getEnv("ADMIN_PASSWORD", ""),
			Name:     getEnv("ADMIN_NAME", "Admin"),
		},
	}

	if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < 16 {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least 16 characters long")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
