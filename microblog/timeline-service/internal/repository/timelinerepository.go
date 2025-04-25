package repository

import (
	"context"
	"time"

	"github.com/gocql/gocql"
)

type CassandraTimelineRepository struct {
	session *gocql.Session
}

type TimelinePost struct {
	PostID    string
	UserID    string
	AuthorID  string
	Content   string
	CreatedAt time.Time
}

func NewCassandraTimelineRepository(host string) (*CassandraTimelineRepository, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "system" // Conecta sem keyspace pra criar depois
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Cria keyspace e tabela
	if err := setupSchema(session); err != nil {
		return nil, err
	}

	// Abre nova sessão no keyspace microblog
	cluster.Keyspace = "microblog"
	newSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraTimelineRepository{session: newSession}, nil
}

func setupSchema(session *gocql.Session) error {
	// Cria keyspace
	if err := session.Query(`
		CREATE KEYSPACE IF NOT EXISTS microblog 
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
	`).Exec(); err != nil {
		return err
	}

	// Cria tabela timeline
	return session.Query(`
		CREATE TABLE IF NOT EXISTS microblog.timeline (
			user_id text,
			created_at timestamp,
			post_id text,
			author_id text,
			content text,
			PRIMARY KEY (user_id, created_at)
		) WITH CLUSTERING ORDER BY (created_at DESC)
	`).Exec()
}

func (r *CassandraTimelineRepository) GetTimeline(ctx context.Context, userID string, limit int, pagingState []byte) ([]TimelinePost, []byte, error) {
	query := r.session.Query(
		`SELECT post_id, author_id, content, created_at FROM timeline WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`,
		userID, limit,
	).WithContext(ctx)

	if pagingState != nil {
		query = query.PageState(pagingState)
	}

	iter := query.Iter()

	var posts []TimelinePost
	var post TimelinePost

	for iter.Scan(&post.PostID, &post.AuthorID, &post.Content, &post.CreatedAt) {
		post.UserID = userID
		posts = append(posts, post)
	}

	if err := iter.Close(); err != nil {
		return nil, nil, err
	}

	return posts, iter.PageState(), nil
}

func (r *CassandraTimelineRepository) Close() {
	r.session.Close()
}
