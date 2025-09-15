package polyschedule

import (
	"context"
	"net/http"

	"github.com/Koshsky/polyschedule-backend/internal/config"
	"github.com/Koshsky/polyschedule-backend/internal/httpserver"
	"github.com/Koshsky/polyschedule-backend/internal/kafka"
	"github.com/Koshsky/polyschedule-backend/internal/processor"
)

// Service инкапсулирует запуск HTTP и Kafka компонентов.
type Service struct {
	cfg      config.Config
	consumer *kafka.Consumer
	producer *kafka.Producer
	httpSrv  *httpserver.Server
}

// New создаёт сервис из переданной конфигурации.
func New(cfg config.Config) *Service {
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaOutputTopic)
	proc := processor.NewProcessor(producer)
	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: cfg.KafkaBrokers,
		GroupID: cfg.KafkaGroupID,
		Topic:   cfg.KafkaInputTopic,
		Handler: proc.Process,
	})
	httpSrv := httpserver.New(cfg.HTTPAddr)
	return &Service{cfg: cfg, consumer: consumer, producer: producer, httpSrv: httpSrv}
}

// Start запускает HTTP и consumer.
func (s *Service) Start(ctx context.Context) (<-chan error, error) {
	errCh := make(chan error, 1)
	go func() { errCh <- s.httpSrv.Start() }()
	go s.consumer.Start(ctx)
	return errCh, nil
}

// Stop останавливает HTTP и consumer, закрывает producer.
func (s *Service) Stop(ctx context.Context) {
	_ = s.httpSrv.Stop(ctx)
	s.consumer.Stop()
	_ = s.producer.Close(context.Background())
}

// HTTPServer возвращает http.Server для интеграции (например, mux замены).
func (s *Service) HTTPServer() *http.Server { return s.httpSrv.Server() }
