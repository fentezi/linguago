package repositories

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client, ctx context.Context) *RedisRepository {
	return &RedisRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisRepository) Get(key string) (string, error) {
	const op = "repositories.RedisRepository.Get"
	res, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (r *RedisRepository) Set(key string, value string) error {
	const op = "repositories.RedisRepository.Set"
	err := r.client.Set(r.ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisRepository) Delete(key string) error {
	const op = "repositories.RedisRepository.Delete"
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
