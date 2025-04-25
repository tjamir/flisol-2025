package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tjamir/flisol-2025/microblog/timeline-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/timeline-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrGettingTimeline = errors.New("erro ao buscar timeline")

type TimelineHandler struct {
	proto.UnimplementedTimelineServiceServer
	Repo *repository.CassandraTimelineRepository
}

func NewTimelineHandler(repo *repository.CassandraTimelineRepository) *TimelineHandler {
	return &TimelineHandler{Repo: repo}
}

func (h *TimelineHandler) GetTimeline(ctx context.Context, req *proto.GetTimelineRequest) (*proto.GetTimelineResponse, error) {
	var pagingState []byte
	if req.Cursor != "" {
		pagingState = []byte(req.Cursor)
	}

	posts, nextCursor, err := h.Repo.GetTimeline(ctx, req.UserId, int(req.Limit), pagingState)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrGettingTimeline, err).Error())
	}

	var protoPosts []*proto.Post
	for _, p := range posts {
		protoPosts = append(protoPosts, &proto.Post{
			Id:        p.PostID,
			UserId:    p.AuthorID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	var cursor string
	if nextCursor != nil {
		cursor = string(nextCursor)
	}

	return &proto.GetTimelineResponse{
		Posts:      protoPosts,
		NextCursor: cursor,
	}, nil
}
