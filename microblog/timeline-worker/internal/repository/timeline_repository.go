package repository

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

const (
	keyspace       = "microblog"
	timelineTable  = "timeline"
)

type CassandraTimelineRepository struct {
	session *gocql.Session
}

func NewCassandraTimelineRepository(host string) (*CassandraTimelineRepository, error) {
	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second

	// Conecta no cluster
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no Cassandra: %w", err)
	}

	// Garante keyspace e tabela
	if err := setupSchema(session); err != nil {
		session.Close()
		return nil, err
	}

	return &CassandraTimelineRepository{session: session}, nil
}

func (r *CassandraTimelineRepository) Close() {
	r.session.Close()
}

func (r *CassandraTimelineRepository) AddToTimeline(userID, postID, authorID, content string, createdAt time.Time) error {
	query := fmt.Sprintf(`INSERT INTO %s.%s (user_id, post_id, author_id, content, created_at) VALUES (?, ?, ?, ?, ?)`, keyspace, timelineTable)
	return r.session.Query(query, userID, postID, authorID, content, createdAt).Exec()
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