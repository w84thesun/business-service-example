package main

import (
	"context"
	logstd "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	stanio "github.com/nats-io/stan.go"
	logger "github.com/w84thesun/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const (
	minMatchSubscriptions     = 1
	minWorkersForSingleClient = 1
)

type Subscription struct {
	sub     stanio.Subscription
	subject string
}

func main() {
	conf, err := parse()
	if err != nil {
		logstd.Panicf("parse config: %v", err)
	}

	log, err := logger.New(conf.Logger)
	if err != nil {
		logstd.Panicf("new logger: %v", err)
	}

	// TODO: add K8S probes

	log.Infof("server started")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-ch:
		log.Infof("got signal %v", sig)
	}

	log.Infof("stopping server...")
}

func checkGrpcState(log logger.Logger, checkTime time.Duration, nameService string, grpcClient *grpc.ClientConn) {
	lastState := grpcClient.GetState()
	log.Infof("%v grpc state: %v", nameService, lastState)
	for {
		newState := grpcClient.GetState()
		if newState != lastState {
			log.Infof("%v grpc state: %v -> %v", nameService, lastState, newState)
			lastState = newState
		}
		time.Sleep(checkTime)
	}
}

func newGrpcClient(
	log logger.Logger, nameService string, grpcURI string,
	connectTimeout, reconnectInterval, checkStateInterval time.Duration,
) *grpc.ClientConn {
	log.Infof("creating new grpc %v client...", nameService)

	if grpcURI == "" {
		log.Panicf("empty grpcURI for %v client", nameService)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	// validate and sanitize config
	if reconnectInterval < 100*time.Millisecond {
		log.Warnf("%v reconnect interval %v less than 100ms, using 100ms", nameService, reconnectInterval)
		reconnectInterval = 100 * time.Millisecond
	}
	if checkStateInterval < 50*time.Millisecond {
		log.Warnf("%v check state interval %v less than 50ms, using 50ms", nameService, checkStateInterval)
		checkStateInterval = 50 * time.Millisecond
	}

	defaultBackoff := backoff.DefaultConfig
	defaultBackoff.MaxDelay = reconnectInterval

	grpcClient, err := grpc.DialContext(ctx, grpcURI, grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: defaultBackoff}))
	if err != nil {
		log.Panicf("%v grpc client: %v", nameService, err)
	}

	go checkGrpcState(log, checkStateInterval, nameService, grpcClient)

	return grpcClient
}

func closeGrpcClient(log logger.Logger, nameService string, grpcClient *grpc.ClientConn) {
	log.Infof("grpc %v client closing...", nameService)
	err := grpcClient.Close()
	if err != nil {
		log.Errorf("close grpc %v client: %v", nameService, err)
	}
	log.Infof("grpc %v client closed", nameService)
}

func unsubscribe(log logger.Logger, sub *Subscription) {
	err := sub.sub.Unsubscribe()
	if err != nil {
		log.Errorf("unsubscribe from %v: %v", sub.subject, err)
		return
	}
	log.Infof("success unsubscribe from %v", sub.subject)
}
