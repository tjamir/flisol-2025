package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/tjamir/flisol-2025/microblog/follow-service/internal/kafka"
	"github.com/tjamir/flisol-2025/microblog/follow-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/follow-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFollowUser    = errors.New("erro ao seguir usuário")
	ErrGetFollowers  = errors.New("erro ao buscar followers")
)


type FollowHandler struct {
	proto.UnimplementedFollowServiceServer
	Repo      *repository.RedisFollowRepository
	Publisher *kafka.FollowPublisher
}

func NewFollowHandler(repo *repository.RedisFollowRepository, publisher *kafka.FollowPublisher) *FollowHandler {
	return &FollowHandler{Repo: repo, Publisher: publisher}
}

func (h *FollowHandler) Follow(ctx context.Context, req *proto.FollowRequest) (*proto.FollowResponse, error) {
	if err := h.Repo.Follow(ctx, req.FollowerId, req.FolloweeId); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrFollowUser, err).Error())
	}

	// Publicar o evento user.followed
	event := map[string]string{
		"follower_id":  req.FollowerId,
		"followee_id": req.FolloweeId,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrFollowUser, err).Error())
	}

	if err := h.Publisher.PublishUserFollowed(ctx, eventBytes); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrFollowUser, err).Error())
	}

	return &proto.FollowResponse{}, nil
}

func (h *FollowHandler) ListFollowers(ctx context.Context, req *proto.ListFollowersRequest) (*proto.ListFollowersResponse, error) {
	cursor := uint64(0)
	parse, err := strconv.ParseUint(req.Cursor, 10, 64)
	if err == nil {
		cursor = parse
	}
	followers, newCursor, err := h.Repo.ListFollowers(ctx, req.UserId, cursor, int64(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrGetFollowers, err).Error())
	}

	return &proto.ListFollowersResponse{
		FollowerIds: followers,
		NextCursor: fmt.Sprintf("%d", newCursor),
	}, nil
}

