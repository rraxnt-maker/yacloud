package repository

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisRepository struct {
    client *redis.Client
}

func NewRedisRepository(addr, password string, db int) *RedisRepository {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    return &RedisRepository{
        client: client,
    }
}

func (r *RedisRepository) Ping(ctx context.Context) error {
    return r.client.Ping(ctx).Err()
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

func (r *RedisRepository) Delete(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
    result, err := r.client.Exists(ctx, key).Result()
    return result > 0, err
}

func (r *RedisRepository) Increment(ctx context.Context, key string) (int64, error) {
    return r.client.Incr(ctx, key).Result()
}

func (r *RedisRepository) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
    data, err := json.Marshal(value)
    if err != nil {
        return false, err
    }
    return r.client.SetNX(ctx, key, data, ttl).Result()
}

func (r *RedisRepository) Close() error {
    return r.client.Close()
}
