package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	// OneSignal push notification
	OneSignalAppID  string
	OneSignalAPIKey string
	OneSignalAPIURL string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig = &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "finance_tracking"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		OneSignalAppID:  getEnv("ONESIGNAL_APP_ID", ""),
		OneSignalAPIKey: getEnv("ONESIGNAL_REST_API_KEY", ""),
		OneSignalAPIURL: getEnv("ONESIGNAL_API_URL", "https://onesignal.com/api/v1"),
	}

	return AppConfig
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
