package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tjamir/flisol-2025/microblog/follow-service/internal/handler"
	"github.com/tjamir/flisol-2025/microblog/follow-service/internal/kafka"
	"github.com/tjamir/flisol-2025/microblog/follow-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/follow-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	repo := repository.NewRedisFollowRepository(redisAddr)

	// Kafka
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	publisher := kafka.NewFollowPublisher(kafkaBroker)
	if err := publisher.EnsureTopic(ctx); err != nil {
		log.Fatalf("Erro ao garantir tópico Kafka: %v", err)
	}
	defer publisher.Close()

	startFollowService(handler.NewFollowHandler(repo, publisher))
}

func startFollowService(h *handler.FollowHandler) {
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("Falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterFollowServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	log.Println("FollowService rodando na porta :8083")

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

	log.Println("Desligando FollowService...")
	grpcServer.GracefulStop()
}
