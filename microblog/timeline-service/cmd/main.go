package main

import (
	"log"
	"net"
	"os"

	"github.com/tjamir/flisol-2025/microblog/timeline-service/internal/handler"
	"github.com/tjamir/flisol-2025/microblog/timeline-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/timeline-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if cassandraHost == "" {
		cassandraHost = "localhost"
	}

	repo, err := repository.NewCassandraTimelineRepository(cassandraHost)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Cassandra: %v", err)
	}
	defer repo.Close()

	startTimelineService(handler.NewTimelineHandler(repo))
}

func startTimelineService(h *handler.TimelineHandler) {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterTimelineServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	log.Println("TimelineService rodando na porta :8084")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
