package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisFollowRepository struct {
	client *redis.Client
}

func NewRedisFollowRepository(addr string) *RedisFollowRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisFollowRepository{client: client}
}

func (r *RedisFollowRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	key := fmt.Sprintf("followers:%s", followeeID)
	return r.client.SAdd(ctx, key, followerID).Err()
}

func (r *RedisFollowRepository) ListFollowers(ctx context.Context, userID string, cursor uint64, count int64) ([]string, uint64, error) {
	key := fmt.Sprintf("followers:%s", userID)
	return r.client.SScan(ctx, key, cursor, "*", count).Result()
}
