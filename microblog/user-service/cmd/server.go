package main

import (
	"log"
	"net"

	"github.com/tjamir/flisol-2025/microblog/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServiceInterface proto.UserServiceServer

func StartUserService(service proto.UserServiceServer) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	log.Println("UserService rodando na porta :8081")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("falha ao servir: %v", err)
	}

	// graceful shutdown omitido aqui por simplicidade no teste
}
