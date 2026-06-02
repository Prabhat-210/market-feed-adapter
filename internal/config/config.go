package config

import (
	"github.com/caarlos0/env/v11"
)

type FeedConfig struct {
	HeartbeatTimeoutMs int      `env:"HEARTBEAT_TIMEOUT_MS" envDefault:"5000"`
	ReconnectMaxMs     int      `env:"RECONNECT_MAX_MS"     envDefault:"30000"`
	AccessToken        string   `env:"ACCESS_TOKEN"     envDefault:"eyJ0eXAiOiJKV1QiLCJrZXlfaWQiOiJza192MS4wIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiI1U0NWNzciLCJqdGkiOiI2YTFlZDFkODJjNjRkNDZlNmMxNzg4OWEiLCJpc011bHRpQ2xpZW50IjpmYWxzZSwiaXNQbHVzUGxhbiI6ZmFsc2UsImlhdCI6MTc4MDQwNDY5NiwiaXNzIjoidWRhcGktZ2F0ZXdheS1zZXJ2aWNlIiwiZXhwIjoxNzgwNDM3NjAwfQ.lsQeVEA_D3GfidnyPgC_3QtC6Nssf6yvTCnYL3v-piA"`
	Instruments        []string `env:"INSTRUMENTS" envSeparator:"," envDefault:"NSE_INDEX|Nifty 50,NSE_INDEX|Nifty Bank"`
	Mode               string   `env:"MODE" envDefault:"full"`
	ConnectionChanSize int      `env:"CONNECTION_CHAN_SIZE"   envDefault:"1000"`
}
type Config struct {
	// service
	ServiceName string `env:"SERVICE_NAME" envDefault:"market-data-service"`
	Environment string `env:"ENVIRONMENT"  envDefault:"dev"`
	LogLevel    string `env:"LOG_LEVEL"    envDefault:"info"`

	FeedConfig
	// pipeline
	CollectorCount   int `env:"COLLECTOR_COUNT"   envDefault:"10"`
	InterpreterCount int `env:"INTERPRETER_COUNT" envDefault:"20"`
	BufferSize       int `env:"BUFFER_SIZE"       envDefault:"10000"`
	BatchMs          int `env:"BATCH_MS"          envDefault:"5"`

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
