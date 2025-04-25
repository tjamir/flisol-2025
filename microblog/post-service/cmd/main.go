package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tjamir/flisol-2025/microblog/post-service/internal/handler"
	"github.com/tjamir/flisol-2025/microblog/post-service/internal/kafka"
	"github.com/tjamir/flisol-2025/microblog/post-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/post-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	ctx := context.Background()
	repo, err := repository.NewDynamoPostRepository(ctx, endpoint)
	if err != nil {
		log.Fatalf("Erro ao conectar ao DynamoDB: %v", err)
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	publisher := kafka.NewPostPublisher(kafkaBroker)
	if err := publisher.EnsureTopic(ctx); err != nil {
		log.Fatalf("Erro ao garantir tópico Kafka: %v", err)
	}
	defer publisher.Close()

	startPostService(handler.NewPostHandler(repo, publisher))
}

func startPostService(h *handler.PostHandler) {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPostServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	log.Println("PostService rodando na porta :8082")

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Falha ao servir: %v", err)
		}
	}()

	// Espera sinal de interrupção
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Desligando PostService...")
	grpcServer.GracefulStop()
}
