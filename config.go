package workr

import "go.uber.org/zap"

// Config for workr
type Config struct {
	Concurrency    int
	RedisHost      string
	RedisPort      int
	RedisPassword  string
	RedisDB        int
	RedisQueueName string
	Logger         *zap.Logger
}

// SetDefaults on config
func (cfg *Config) SetDefaults() {
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 5
	}

	if cfg.RedisHost == "" {
		cfg.RedisHost = "localhost"
	}

	if cfg.RedisPort == 0 {
		cfg.RedisPort = 6379
	}

	if cfg.Logger == nil {
		logger, _ := zap.NewProduction()
		cfg.Logger = logger
	}

	if cfg.RedisQueueName == "" {
		cfg.RedisQueueName = "wrkr_jobs_queue"
	}
}
