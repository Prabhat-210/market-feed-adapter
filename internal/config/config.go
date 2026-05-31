package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	// service
	ServiceName string `env:"SERVICE_NAME" envDefault:"market-data-service"`
	Environment string `env:"ENVIRONMENT"  envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL"    envDefault:"info"`

	// pipeline
	CollectorCount   int `env:"COLLECTOR_COUNT"   envDefault:"10"`
	InterpreterCount int `env:"INTERPRETER_COUNT" envDefault:"20"`
	BufferSize       int `env:"BUFFER_SIZE"       envDefault:"10000"`
	BatchMs          int `env:"BATCH_MS"          envDefault:"5"`

	// feed
	FeedURL            string `env:"FEED_URL"`
	HeartbeatTimeoutMs int    `env:"HEARTBEAT_TIMEOUT_MS" envDefault:"5000"`
	ReconnectMaxMs     int    `env:"RECONNECT_MAX_MS"     envDefault:"30000"`

	// kafka
	KafkaBrokers       string `env:"KAFKA_BROKERS"`
	KafkaTopic         string `env:"KAFKA_TOPIC"         envDefault:"price.updated"`
	KafkaDLQTopic      string `env:"KAFKA_DLQ_TOPIC"     envDefault:"price.dlq"`
	KafkaRetryAttempts int    `env:"KAFKA_RETRY_ATTEMPTS" envDefault:"3"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
