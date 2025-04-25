package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/client"
	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/repository"
)

type EventHandler struct {
	FollowClient   *client.FollowClient
	PostClient     *client.PostClient
	TimelineRepo   *repository.CassandraTimelineRepository
}

func NewEventHandler(followClient *client.FollowClient, postClient *client.PostClient, timelineRepo *repository.CassandraTimelineRepository) *EventHandler {
	return &EventHandler{
		FollowClient: followClient,
		PostClient:   postClient,
		TimelineRepo: timelineRepo,
	}
}

func (h *EventHandler) HandlePostCreated(ctx context.Context, msg []byte) error {
	var event struct {
		PostID    string `json:"post_id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, event.CreatedAt)
	if err != nil {
		return err
	}

	// Buscar os followers via FollowClient
	followers, err := h.FollowClient.GetFollowers(ctx, event.UserID)
	if err != nil {
		return err
	}

	// Inserir post na timeline de cada follower
	for _, followerID := range followers {
		err := h.TimelineRepo.AddToTimeline(followerID, event.PostID, event.UserID, event.Content, createdAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *EventHandler) HandleUserFollowed(ctx context.Context, msg []byte) error {
	var event struct {
		FollowerID string `json:"follower_id"`
		FolloweeID string `json:"followee_id"`
	}
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	// Buscar posts antigos via PostClient
	posts, err := h.PostClient.ListPosts(ctx, event.FolloweeID)
	if err != nil {
		return err
	}

	// Adicionar os posts antigos na timeline do follower
	for _, post := range posts {
		createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			return err
		}

		err = h.TimelineRepo.AddToTimeline(event.FollowerID, post.PostID, post.AuthorID, post.Content, createdAt)
		if err != nil {
			return err
		}
	}

	return nil
}
