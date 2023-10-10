package main

import (
	"time"

	env "github.com/caarlos0/env/v6"

	logger "github.com/w84thesun/logger"
)

type Config struct {
	Logger     logger.LoggingConfig
	ProbesPort string `env:"PROBES_PORT"`
	Prometheus Prometheus
	GRPC       GRPC
}

type Dispatcher struct {
	WorkersForSingleClient int `env:"WORKERS_FOR_SINGLE_CLIENT" envDefault:"50"`
	MaxMatchSubscriptions  int `env:"MAX_MATCH_SUBSCRIPTIONS" envDefault:"50"`
}

type Prometheus struct {
	Port string `env:"PROMETHEUS_PORT"`
}

type GRPC struct {
	ConnectTimeout     time.Duration `env:"GRPC_CONNECTION_TIMEOUT" envDefault:"2s"`
	RequestTimeout     time.Duration `env:"GRPC_REQUEST_TIMEOUT" envDefault:"2s"`
	ReconnectInterval  time.Duration `env:"GRPC_RECONNECT_INTERVAL" envDefault:"1s"`
	CheckStateInterval time.Duration `env:"GRPC_CHECK_STATE_INTERVAL" envDefault:"500ms"`
}

func parse() (*Config, error) {
	conf := Config{}
	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
