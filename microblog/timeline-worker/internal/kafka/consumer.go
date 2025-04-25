package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler func(context.Context, []byte) error
}

func NewConsumer(broker, topic, groupID string, handler func(context.Context, []byte) error) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			continue
		}

		if err := c.handler(ctx, m.Value); err != nil {
			log.Printf("Erro ao processar mensagem: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
