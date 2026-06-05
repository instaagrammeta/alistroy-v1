package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config aggregates all runtime configuration sourced from environment vars.
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
	Google   GoogleOAuthConfig
}

type PostgresConfig struct {
	Host, Port, User, Password, DBName, SSLMode string
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}

type RedisConfig struct {
	Host, Port, Password string
}

func (r RedisConfig) Addr() string { return r.Host + ":" + r.Port }

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
	Phone            string
	WhatsApp         string
	Telegram         string
	TelegramUsername string
}

type AdminBootstrap struct {
	Email    string
	Phone    string
	Password string
	Name     string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (g GoogleOAuthConfig) Enabled() bool {
	return g.ClientID != "" && g.ClientSecret != "" && g.RedirectURL != ""
}

// Load parses configuration from env (and optionally a local .env file).
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:       getStr("API_ENV", "development"),
		LogLevel:  getStr("API_LOG_LEVEL", "info"),
		Host:      getStr("API_HOST", "0.0.0.0"),
		Port:      getStr("API_PORT", "8080"),
		PublicURL: getStr("PUBLIC_URL", "http://localhost:3000"),
		Postgres: PostgresConfig{
			Host:     getStr("POSTGRES_HOST", "localhost"),
			Port:     getStr("POSTGRES_PORT", "5432"),
			User:     getStr("POSTGRES_USER", "alistroy"),
			Password: getStr("POSTGRES_PASSWORD", "alistroy"),
			DBName:   getStr("POSTGRES_DB", "alistroy"),
			SSLMode:  getStr("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getStr("REDIS_HOST", "localhost"),
			Port:     getStr("REDIS_PORT", "6379"),
			Password: getStr("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:     getStr("JWT_SECRET", ""),
			AccessTTL:  time.Duration(getInt("JWT_ACCESS_TTL_MINUTES", 120)) * time.Minute,
			RefreshTTL: time.Duration(getInt("JWT_REFRESH_TTL_HOURS", 720)) * time.Hour,
		},
		Upload: UploadConfig{
			Dir:        getStr("UPLOAD_DIR", "./uploads"),
			MaxSizeMB:  int64(getInt("UPLOAD_MAX_SIZE_MB", 50)),
			PublicBase: getStr("UPLOAD_PUBLIC_BASE", "/uploads"),
		},
		CORS: CORSConfig{AllowedOrigins: splitTrim(getStr("CORS_ALLOWED_ORIGINS", "*"))},
		Rate: RateConfig{
			RPS:   getInt("RATE_LIMIT_RPS", 30),
			Burst: getInt("RATE_LIMIT_BURST", 60),
		},
		Market: MarketplaceContacts{
			Phone:            getStr("MARKETPLACE_PHONE", ""),
			WhatsApp:         getStr("MARKETPLACE_WHATSAPP", ""),
			Telegram:         getStr("MARKETPLACE_TELEGRAM", ""),
			TelegramUsername: getStr("MARKETPLACE_TELEGRAM_USERNAME", ""),
		},
		Admin: AdminBootstrap{
			Email:    getStr("ADMIN_EMAIL", ""),
			Phone:    getStr("ADMIN_PHONE", ""),
			Password: getStr("ADMIN_PASSWORD", ""),
			Name:     getStr("ADMIN_NAME", "Admin"),
		},
		Google: GoogleOAuthConfig{
			ClientID:     getStr("GOOGLE_OAUTH_CLIENT_ID", ""),
			ClientSecret: getStr("GOOGLE_OAUTH_CLIENT_SECRET", ""),
			RedirectURL:  getStr("GOOGLE_OAUTH_REDIRECT_URL", ""),
		},
	}

	if len(cfg.JWT.Secret) < 16 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 16 characters")
	}
	return cfg, nil
}

func getStr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
func getInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
