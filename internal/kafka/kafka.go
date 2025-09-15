package kafka

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	kgo "github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key, value []byte) error

type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topic   string
	Handler MessageHandler
}

type Consumer struct {
	reader  *kgo.Reader
	handler MessageHandler
	stopped chan struct{}
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	r := kgo.NewReader(kgo.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: r, handler: cfg.Handler, stopped: make(chan struct{})}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Info().Msg("kafka consumer started")
	defer close(c.stopped)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info().Msg("consumer context done")
				return
			}
			log.Error().Err(err).Msg("read message failed")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if c.handler != nil {
			if err := c.handler(ctx, m.Key, m.Value); err != nil {
				log.Error().Err(err).Msg("handler failed")
			}
		}
	}
}

func (c *Consumer) Stop() {
	_ = c.reader.Close()
	<-c.stopped
}

type Producer struct {
	writer *kgo.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kgo.Writer{
		Addr:         kgo.TCP(brokers...),
		Topic:        topic,
		RequiredAcks: kgo.RequireOne,
		Balancer:     &kgo.LeastBytes{},
	}
	return &Producer{writer: w}
}

func (p *Producer) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kgo.Message{Key: key, Value: value})
}

func (p *Producer) Close(ctx context.Context) error {
	return p.writer.Close()
}
