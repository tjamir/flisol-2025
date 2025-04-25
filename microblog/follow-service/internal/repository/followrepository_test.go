package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFollowAndListFollowers(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	ctx := context.Background()
	repo := NewRedisFollowRepository(addr)

	followeeID := "user-valentina"
	followerIDs := []string{"follower1", "follower2", "follower3"}

	// Cleanup antes
	repo.client.Del(ctx, "followers:"+followeeID)

	// Executa Follow
	for _, followerID := range followerIDs {
		err := repo.Follow(ctx, followerID, followeeID)
		assert.NoError(t, err)
	}

	// Lista Followers
	var allFollowers []string
	var cursor uint64
	for {
		results, nextCursor, err := repo.ListFollowers(ctx, followeeID, cursor, 2)
		assert.NoError(t, err)
		allFollowers = append(allFollowers, results...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	assert.ElementsMatch(t, followerIDs, allFollowers)
}
