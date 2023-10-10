package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DependantServicePort string `env:"MOCK_SERVICE_PORT" envDefault:"9000"`
}

func main() {
	conf := Config{}
	err := env.Parse(&conf)
	if err != nil {
		log.Panicf("parse config: %v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	log.Println("stopping dependant mock...")
}
