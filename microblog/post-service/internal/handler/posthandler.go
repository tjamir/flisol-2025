package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tjamir/flisol-2025/microblog/post-service/internal/kafka"
	"github.com/tjamir/flisol-2025/microblog/post-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/post-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCreatingPost = errors.New("erro ao criar post")
	ErrListingPosts = errors.New("erro ao listar posts")
)

type PostHandler struct {
	proto.UnimplementedPostServiceServer
	Repo *repository.DynamoPostRepository
	Publisher *kafka.PostPublisher
}

func NewPostHandler(repo *repository.DynamoPostRepository, publisher *kafka.PostPublisher) *PostHandler {
	return &PostHandler{Repo: repo, Publisher: publisher}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *proto.CreatePostRequest) (*proto.CreatePostResponse, error) {
	postID, err := h.Repo.CreatePost(ctx, req.UserId, req.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrCreatingPost, err).Error())
	}

	// Publica o evento post.created
	event := map[string]string{
		"post_id":    postID,
		"user_id":    req.UserId,
		"content":    req.Content,
		"created_at": time.Now().Format(time.RFC3339),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrCreatingPost, err).Error())
	}

	if err := h.Publisher.PublishPostCreated(ctx, eventBytes); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrCreatingPost, err).Error())
	}

	return &proto.CreatePostResponse{PostId: postID}, nil
}


func (h *PostHandler) ListPosts(ctx context.Context, req *proto.ListPostsRequest) (*proto.ListPostsResponse, error) {
	posts, nextCursor, err := h.Repo.ListPosts(ctx, req.UserId, req.Limit, req.Cursor)
	if err != nil {
		return nil, status.Error(codes.Internal, 
			fmt.Errorf("%w: %v", ErrListingPosts, err).Error())
	}

	var protoPosts []*proto.Post
	for _, p := range posts {
		protoPosts = append(protoPosts, &proto.Post{
			Id:        p.ID,
			UserId:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
		})
	}

	return &proto.ListPostsResponse{
		Posts:      protoPosts,
		NextCursor: nextCursor,
	}, nil
}
