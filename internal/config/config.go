package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST            string
	DB_PASSWORD        string
	DB_PORT            int64
	DB_NAME            string
	DB_USER            string
	DB_SSLMODE         string
	JWTsecret          string
	REDIS_ADDR         string
	REDIS_PASSWORD     string
	CORSAllowedOrigins []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &Config{}
	config.DB_NAME = getEnv("DB_NAME", "loci_db")
	config.DB_HOST = getEnv("DB_HOST", "localhost")
	config.DB_PASSWORD = getEnv("DB_PASSWORD", "secretpassword")
	config.DB_PORT = int64(getEnvAsInt("DB_PORT", 5432))
	config.DB_USER = getEnv("DB_USER", "admin")
	config.DB_SSLMODE = getEnv("DB_SSLMODE", "disable")
	config.JWTsecret = getEnv("JWTSecret", "mydogiscalledrufus")
	config.REDIS_ADDR = getEnv("REDIS_ADDR", "localhost:6379")
	config.REDIS_PASSWORD = getEnv("REDIS_PASSWORD", "secretpassword")
	//config.CORSOrigin = getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3001")
	config.CORSAllowedOrigins = getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3001"}, ",")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, separator)
}
