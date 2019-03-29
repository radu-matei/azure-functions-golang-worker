package main

import (
	"context"
	"flag"
	"math"

	log "github.com/sirupsen/logrus"

	"github.com/radu-matei/azure-functions-golang-worker/worker"
)

var (
	flagDebug            bool
	host                 string
	port                 int
	workerID             string
	requestID            string
	grpcMaxMessageLength int
)

func init() {

	flag.BoolVar(&flagDebug, "debug", true, "enable verbose output")
	flag.StringVar(&host, "host", "127.0.0.1", "RPC Server Host")
	flag.IntVar(&port, "port", 0, "RPC Server Port")
	flag.StringVar(&workerID, "workerId", "", "RPC Server Worker ID")
	flag.StringVar(&requestID, "requestId", "", "Request ID")
	flag.IntVar(&grpcMaxMessageLength, "grpcMaxMessageLength", math.MaxInt32, "Max message length")

	flag.Parse()

	if flagDebug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Debugf("attempting to start grpc connection to server %s:%d with worker id %s, request id %s and max message length %d", host, port, workerID, requestID, grpcMaxMessageLength)

	cfg := &worker.ClientConfig{
		Host:             host,
		Port:             port,
		WorkerID:         workerID,
		RequestID:        requestID,
		MaxMessageLength: grpcMaxMessageLength,
	}
	client := worker.New(cfg)
	if err := client.Connect(); err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer client.Disconnect()
	log.Debugf("started grpc connection...")
	client.StartEventStream(context.Background())
}
