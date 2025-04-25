package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCassandraTimelineRepository_CreatesSchemaAndInserts(t *testing.T) {
	repo, err := NewCassandraTimelineRepository("localhost:9042")
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	var fetchedCreatedAt time.Time

	userID := "test-user"
	postID := "test-post"
	authorID := "test-author"
	content := "Hello Timeline!"
	createdAt := time.Now()

	t.Cleanup(func() {
		defer repo.Close()
		// Cleanup
		deleteQuery := `DELETE FROM microblog.timeline WHERE user_id = ? AND created_at = ?`
		err = repo.session.Query(deleteQuery, userID, fetchedCreatedAt).Exec()
		assert.NoError(t, err)
		
	})

	// Inserir na timeline
	err = repo.AddToTimeline(userID, postID, authorID, content, createdAt)
	assert.NoError(t, err)

	// Verificar se o dado foi inserido
	var fetchedPostID, fetchedAuthorID, fetchedContent string
	

	query := `SELECT post_id, author_id, content, created_at FROM microblog.timeline WHERE user_id = ? LIMIT 1`
	err = repo.session.Query(query, userID).Scan(&fetchedPostID, &fetchedAuthorID, &fetchedContent, &fetchedCreatedAt)
	assert.NoError(t, err)
	assert.Equal(t, postID, fetchedPostID)
	assert.Equal(t, authorID, fetchedAuthorID)
	assert.Equal(t, content, fetchedContent)
	
}
