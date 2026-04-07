package service

import (
    "context"
    "time"
    "cache-service/internal/repository"
)

type CacheService struct {
    repo *repository.RedisRepository
}

func NewCacheService(repo *repository.RedisRepository) *CacheService {
    return &CacheService{repo: repo}
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    return s.repo.Set(ctx, key, value, ttl)
}

func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
    return s.repo.Get(ctx, key, dest)
}

func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
    return s.repo.Delete(ctx, keys...)
}

func (s *CacheService) Exists(ctx context.Context, key string) (bool, error) {
    return s.repo.Exists(ctx, key)
}

func (s *CacheService) Increment(ctx context.Context, key string) (int64, error) {
    return s.repo.Increment(ctx, key)
}

func (s *CacheService) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
    return s.repo.SetNX(ctx, key, value, ttl)
}
