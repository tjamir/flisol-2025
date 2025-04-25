package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeline(t *testing.T) {
	ctx := context.Background()
	repo, err := NewCassandraTimelineRepository("localhost")
	assert.NoError(t, err)
	defer repo.Close()

	// Insere posts manualmente
	userID := "user-timeline"
	posts := []TimelinePost{
		{PostID: "post1", UserID: userID, AuthorID: "author1", Content: "Post 1", CreatedAt: time.Now()},
		{PostID: "post2", UserID: userID, AuthorID: "author2", Content: "Post 2", CreatedAt: time.Now().Add(-time.Minute)},
	}

	for _, p := range posts {
		err := repo.session.Query(
			`INSERT INTO timeline (user_id, post_id, author_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
			p.UserID, p.PostID, p.AuthorID, p.Content, p.CreatedAt,
		).WithContext(ctx).Exec()
		assert.NoError(t, err)
	}

	// Busca timeline
	result, _, err := repo.GetTimeline(ctx, userID, 10, nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Post 1", result[0].Content)
	assert.Equal(t, "Post 2", result[1].Content)
}
