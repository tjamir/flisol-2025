package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tjamir/flisol-2025/microblog/user-service/proto"
	"github.com/tjamir/flisol-2025/microblog/commons/test"
)

// fakeUserService implementa o serviço gRPC com lógica simples
type fakeUserService struct {
	proto.UnimplementedUserServiceServer
}

func (f *fakeUserService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return &proto.RegisterResponse{
		Id:       "fake-id",
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func TestRegisterIntegration(t *testing.T) {
	server := test.StartGRPCServer(func(s *grpc.Server) {
		proto.RegisterUserServiceServer(s, &fakeUserService{})
	})

	ctx := context.Background()

	clientConn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(server.Dialer()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer clientConn.Close()

	client := proto.NewUserServiceClient(clientConn)

	resp, err := client.Register(ctx, &proto.RegisterRequest{
		Username: "valentina",
		Email:    "valen@example.com",
		Password: "123456",
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "valentina", resp.Username)
	assert.Equal(t, "valen@example.com", resp.Email)
}

