package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    ServerPort     string
    RedisHost      string
    RedisPort      string
    RedisPassword  string
    RedisDB        int
    DefaultTTL     time.Duration
    MaxMemory      string
    EvictionPolicy string
}

func Load() *Config {
    return &Config{
        ServerPort:     getEnv("SERVER_PORT", "8082"),
        RedisHost:      getEnv("REDIS_HOST", "redis"),
        RedisPort:      getEnv("REDIS_PORT", "6379"),
        RedisPassword:  getEnv("REDIS_PASSWORD", ""),
        RedisDB:        getEnvAsInt("REDIS_DB", 0),
        DefaultTTL:     time.Duration(getEnvAsInt("CACHE_TTL_SECONDS", 3600)) * time.Second,
        MaxMemory:      getEnv("REDIS_MAX_MEMORY", "256mb"),
        EvictionPolicy: getEnv("REDIS_EVICTION_POLICY", "allkeys-lru"),
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
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
