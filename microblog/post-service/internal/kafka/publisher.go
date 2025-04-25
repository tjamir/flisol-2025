package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const PostCreatedTopic = "post.created"

type PostPublisher struct {
	writer *kafka.Writer
}

func NewPostPublisher(broker string) *PostPublisher {
	return &PostPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        PostCreatedTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *PostPublisher) EnsureTopic(ctx context.Context) error {
	log.Printf("kafka addr: %s\n", p.writer.Addr.String())
	conn, err := kafka.Dial("tcp", p.writer.Addr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	conn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))

	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{{
		Topic:             PostCreatedTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}}

	return conn.CreateTopics(topicConfigs...)
}

func (p *PostPublisher) PublishPostCreated(ctx context.Context, message []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339Nano)),
		Value: message,
	})
}

func (p *PostPublisher) Close() error {
	return p.writer.Close()
}
