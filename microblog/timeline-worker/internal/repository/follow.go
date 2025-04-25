package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type FollowRepository struct {
	client *redis.Client
}

func NewFollowRepository(addr string) *FollowRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &FollowRepository{client: client}
}

// GetFollowers retorna a lista de followers de um usuário
func (r *FollowRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	key := "followers:" + userID
	return r.client.SMembers(ctx, key).Result()
}
