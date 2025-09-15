package config

import (
	"github.com/caarlos0/env/v10"
)

// Config описывает параметры сервиса, управляемые через переменные окружения
type Config struct {
	KafkaBrokers     []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`
	KafkaGroupID     string   `env:"KAFKA_GROUP_ID" envDefault:"polyschedule-backend"`
	KafkaInputTopic  string   `env:"KAFKA_INPUT_TOPIC" envDefault:"schedule.requests"`
	KafkaOutputTopic string   `env:"KAFKA_OUTPUT_TOPIC" envDefault:"schedule.results"`
	HTTPAddr         string   `env:"HTTP_ADDR" envDefault:":8080"`
	LogLevel         string   `env:"LOG_LEVEL" envDefault:"info"`
}

// Load загружает конфигурацию из переменных окружения с дефолтами
func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
