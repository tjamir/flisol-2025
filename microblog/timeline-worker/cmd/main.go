package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/client"
	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/handler"
	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/kafka"
	"github.com/tjamir/flisol-2025/microblog/timeline-worker/internal/repository"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configuração de serviços
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}

	followClient, _ := client.NewFollowClient("follow-service:8083")
	postClient, _ := client.NewPostClient("post-service:8082")
	timelineRepo, _ := repository.NewCassandraTimelineRepository("cassandra:9042")
	defer timelineRepo.Close()

	eventHandler := handler.NewEventHandler(followClient, postClient, timelineRepo)

	// Consumers Kafka
	postConsumer := kafka.NewConsumer(kafkaBroker, "post.created", "worker-group", eventHandler.HandlePostCreated)
	userFollowedConsumer := kafka.NewConsumer(kafkaBroker, "user.followed", "worker-group", eventHandler.HandleUserFollowed)

	// Start consumers
	go postConsumer.Start(ctx)
	go userFollowedConsumer.Start(ctx)

	log.Println("Worker rodando...")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Desligando worker...")
	postConsumer.Close()
	userFollowedConsumer.Close()
}
