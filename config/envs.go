package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth
	Database
	HttpServer
}

type HttpServer struct {
	Port string
	Env  string
}

type Database struct {
	User            string
	Pass            string
	Host            string
	Port            string
	Name            string
	SSLMode         string
	Timezone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Auth struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	TokenExpiry        time.Duration
	RefreshExpiry      time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default env")
	}

	return &Config{
		Auth: Auth{
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "defaultsecret"),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "defaultrefreshsecret"),
			TokenExpiry:        getDurationEnv("TOKEN_EXPIRY", time.Minute*15),
			RefreshExpiry:      getDurationEnv("REFRESH_EXPIRY", time.Hour*24*7),
		},
		Database: Database{
			User:            getEnv("DB_USER", "postgres"),
			Pass:            getEnv("DB_PASS", "password"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			Name:            getEnv("DB_NAME", "dbname"),
			SSLMode:         getEnv("DB_SSL_MODE", "require"),
			Timezone:        getEnv("DB_TIMEZONE", "UTC"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", time.Minute*5),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", time.Minute*5),
		},
		HttpServer: HttpServer{
			Port: getEnv("HTTP_PORT", "8080"),
			Env:  getEnv("HTTP_ENV", "development"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getIntEnv(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return intValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(strValue)
	if err != nil {
		log.Printf("Warning: Invalid duration value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}

	return durationValue
}
