package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const UserFollowedTopic = "user.followed"

type FollowPublisher struct {
	writer *kafka.Writer
}

func NewFollowPublisher(broker string) *FollowPublisher {
	return &FollowPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        UserFollowedTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *FollowPublisher) EnsureTopic(ctx context.Context) error {
	log.Printf("kafka addr: %s\n", p.writer.Addr.String())
	conn, err := kafka.Dial("tcp", p.writer.Addr.String())
	if err != nil {
		return fmt.Errorf("erro ao conectar no Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("erro ao obter controller: %w", err)
	}
	conn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("erro ao conectar no controller: %w", err)
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{{
		Topic:             UserFollowedTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}}

	return conn.CreateTopics(topicConfigs...)
}

func (p *FollowPublisher) PublishUserFollowed(ctx context.Context, message []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339Nano)),
		Value: message,
	})
}

func (p *FollowPublisher) Close() error {
	return p.writer.Close()
}
