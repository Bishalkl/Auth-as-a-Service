package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds all the configuration settings
type AppConfig struct {
	AppName                  string
	AppEnv                   string
	Port                     string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	JWTSecret                string
	AccessTokenExpireMinutes string
	RefreshTokenExpireHours  string
	MailAPIKey               string
	MailSender               string
}

// Config is the global config variable
var Config *AppConfig

// LoadEnv loads the environment variables and sets them in Config
func LoadEnv() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found. Using system envs.")
	}

	Config = &AppConfig{
		AppName:                  MustGetEnvOrDefault("APP_NAME", "AuthService"),
		AppEnv:                   MustGetEnvOrDefault("APP_ENV", "development"),
		Port:                     MustGetEnvOrDefault("PORT", "8080"),
		DBHost:                   MustGetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:                   MustGetEnvOrDefault("DB_PORT", "5432"),
		DBUser:                   MustGetEnvOrDefault("DB_USER", "postgres"),
		DBPassword:               MustGetEnvOrDefault("DB_PASSWORD", "password"),
		DBName:                   MustGetEnvOrDefault("DB_NAME", "auth_service"),
		RedisHost:                MustGetEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:                MustGetEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword:            MustGetEnvOrDefault("REDIS_PASSWORD", ""),
		JWTSecret:                MustGetEnvOrDefault("JWT_SECRET", "supersecretkey"),
		AccessTokenExpireMinutes: MustGetEnvOrDefault("ACCESS_TOKEN_EXPIRE_MINUTES", "15"),
		RefreshTokenExpireHours:  MustGetEnvOrDefault("REFRESH_TOKEN_EXPIRE_HOURS", "24"),
		MailAPIKey:               MustGetEnvOrDefault("MAIL_API_KEY", ""),
		MailSender:               MustGetEnvOrDefault("MAIL_SENDER", ""),
	}
}

// MustGetEnvOrDefault returns env var or fallback if not set
func MustGetEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
