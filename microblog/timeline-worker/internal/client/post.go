package client

import (
	"context"

	"github.com/tjamir/flisol-2025/microblog/post-service/proto"
	"google.golang.org/grpc"
)

type PostClient struct {
	client proto.PostServiceClient
}

func NewPostClient(addr string) (*PostClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // ou grpc.WithTransportCredentials(insecure.NewCredentials())
	if err != nil {
		return nil, err
	}

	return &PostClient{
		client: proto.NewPostServiceClient(conn),
	}, nil
}

type Post struct {
	PostID    string
	AuthorID  string
	Content   string
	CreatedAt string
}

func (p *PostClient) ListPosts(ctx context.Context, userID string) ([]Post, error) {
	resp, err := p.client.ListPosts(ctx, &proto.ListPostsRequest{UserId: userID, Limit: 1000})
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, item := range resp.Posts {
		posts = append(posts, Post{
			PostID:    item.Id,
			AuthorID:  item.UserId,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		})
	}

	return posts, nil
}
