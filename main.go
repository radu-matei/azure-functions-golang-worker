package main

import (
	"context"
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/worker"
)

var (
	flagDebug bool
	host      string
	port      int
	workerID  string
	requestID string
)

func init() {

	flag.BoolVar(&flagDebug, "debug", true, "enable verbose output")
	flag.StringVar(&host, "host", "127.0.0.1", "RPC Server Host")
	flag.IntVar(&port, "port", 0, "RPC Server Port")
	flag.StringVar(&workerID, "workerId", "", "RPC Server Worker ID")
	flag.StringVar(&requestID, "requestId", "", "Request ID")

	flag.Parse()

	if flagDebug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Debugf("attempting to start grpc connection to server %s:%d with worker id %s and request id %s", host, port, workerID, requestID)

	conn, err := worker.GetGRPCConnection(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer conn.Close()
	log.Debugf("started grpc connection...")

	cfg := &worker.ClientConfig{
		Host: host,
	}
	client := worker.NewClient(cfg, conn)
	client.StartEventStream(context.Background())
}
