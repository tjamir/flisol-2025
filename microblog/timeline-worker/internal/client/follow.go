package client

import (
	"context"

	"github.com/tjamir/flisol-2025/microblog/follow-service/proto"
	"google.golang.org/grpc"
)

type FollowClient struct {
	client proto.FollowServiceClient
}

func NewFollowClient(addr string) (*FollowClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // ou grpc.WithTransportCredentials(insecure.NewCredentials())
	if err != nil {
		return nil, err
	}

	return &FollowClient{
		client: proto.NewFollowServiceClient(conn),
	}, nil
}

func (f *FollowClient) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	resp, err := f.client.ListFollowers(ctx, &proto.ListFollowersRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.FollowerIds, nil
}
