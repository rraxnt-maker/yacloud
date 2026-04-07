package main

import (
    "context"
    "log"
    "net/http"
    "time"
    "cache-service/internal/config"
    "cache-service/internal/handler"
    "cache-service/internal/repository"
    "cache-service/internal/service"
)

func main() {
    cfg := config.Load()
    
    redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
    redisRepo := repository.NewRedisRepository(redisAddr, cfg.RedisPassword, cfg.RedisDB)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := redisRepo.Ping(ctx); err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    
    log.Println("Connected to Redis successfully")
    
    cacheService := service.NewCacheService(redisRepo)
    cacheHandler := handler.NewCacheHandler(cacheService)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", cacheHandler.HandleRequest)
    
    port := cfg.ServerPort
    log.Printf("Cache Service starting on port %s", port)
    log.Printf("Available endpoints:")
    log.Printf("  GET  /health")
    log.Printf("  POST /api/v1/cache")
    log.Printf("  GET  /api/v1/cache/get/{key}")
    log.Printf("  DELETE /api/v1/cache/del/{key}")
    log.Printf("  GET  /api/v1/cache/exists/{key}")
    log.Printf("  POST /api/v1/cache/incr/{key}")
    
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
