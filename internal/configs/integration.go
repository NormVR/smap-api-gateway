package configs

import "os"

type IntegrationConfig struct {
	KafkaAddr     string
	RedisAddr     string
	RedisPassword string
}

func NewIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		KafkaAddr:     os.Getenv("KAFKA_BROKERS"),
		RedisAddr:     os.Getenv("REDIS_ADDRESS"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
