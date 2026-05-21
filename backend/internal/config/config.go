package config

import "os"

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	WebhookURL string
	UploadDir  string
}

func Load() Config {
	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "fleetify"),
		DBPassword: getEnv("DB_PASSWORD", "fleetify_password"),
		DBName:     getEnv("DB_NAME", "fleetify_db"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		UploadDir:  getEnv("UPLOAD_DIR", "uploads"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
