package handler

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"
    "cache-service/internal/service"
)

type CacheHandler struct {
    service *service.CacheService
}

func NewCacheHandler(service *service.CacheService) *CacheHandler {
    return &CacheHandler{service: service}
}

type SetRequest struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
    TTL   int         `json:"ttl_seconds"`
}

func (h *CacheHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    
    // Health check
    if path == "/health" && r.Method == "GET" {
        h.HealthCheck(w, r)
        return
    }
    
    // POST /api/v1/cache - сохранить
    if path == "/api/v1/cache" && r.Method == "POST" {
        h.Set(w, r)
        return
    }
    
    // GET /api/v1/cache/get/{key} - получить
    if strings.HasPrefix(path, "/api/v1/cache/get/") && r.Method == "GET" {
        key := strings.TrimPrefix(path, "/api/v1/cache/get/")
        h.Get(w, r, key)
        return
    }
    
    // DELETE /api/v1/cache/del/{key} - удалить
    if strings.HasPrefix(path, "/api/v1/cache/del/") && r.Method == "DELETE" {
        key := strings.TrimPrefix(path, "/api/v1/cache/del/")
        h.Delete(w, r, key)
        return
    }
    
    // GET /api/v1/cache/exists/{key} - проверить существование
    if strings.HasPrefix(path, "/api/v1/cache/exists/") && r.Method == "GET" {
        key := strings.TrimPrefix(path, "/api/v1/cache/exists/")
        h.Exists(w, r, key)
        return
    }
    
    // POST /api/v1/cache/incr/{key} - инкремент
    if strings.HasPrefix(path, "/api/v1/cache/incr/") && r.Method == "POST" {
        key := strings.TrimPrefix(path, "/api/v1/cache/incr/")
        h.Increment(w, r, key)
        return
    }
    
    http.NotFound(w, r)
}

func (h *CacheHandler) Set(w http.ResponseWriter, r *http.Request) {
    var req SetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    ttl := time.Duration(req.TTL) * time.Second
    if ttl == 0 {
        ttl = 3600 * time.Second
    }
    
    ctx := r.Context()
    if err := h.service.Set(ctx, req.Key, req.Value, ttl); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok", "key": req.Key})
}

func (h *CacheHandler) Get(w http.ResponseWriter, r *http.Request, key string) {
    if key == "" {
        http.Error(w, "Key required", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    var value interface{}
    if err := h.service.Get(ctx, key, &value); err != nil {
        http.Error(w, "Key not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"key": key, "value": value})
}

func (h *CacheHandler) Delete(w http.ResponseWriter, r *http.Request, key string) {
    if key == "" {
        http.Error(w, "Key required", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    if err := h.service.Delete(ctx, key); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted", "key": key})
}

func (h *CacheHandler) Exists(w http.ResponseWriter, r *http.Request, key string) {
    if key == "" {
        http.Error(w, "Key required", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    exists, err := h.service.Exists(ctx, key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"key": key, "exists": exists})
}

func (h *CacheHandler) Increment(w http.ResponseWriter, r *http.Request, key string) {
    if key == "" {
        http.Error(w, "Key required", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    value, err := h.service.Increment(ctx, key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"key": key, "value": value})
}

func (h *CacheHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "service": "cache-service",
    })
}
