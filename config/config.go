package config

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort string
    DatabaseURL string
}

func Load() *Config {
    // Load .env file
    _ = godotenv.Load()

    return &Config{
        ServerPort: getEnv("SERVER_PORT", "3000"),
        DatabaseURL : getEnv("DATABASE_URL", ""),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}