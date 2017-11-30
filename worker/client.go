package worker

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"google.golang.org/grpc"
)

// ClientConfig contains all necessary configuration to connect to the Azure Functions Host
type ClientConfig struct {
	Host string
}

// Client that listens for events from the Azure Functions host and executes Golang methods
type Client struct {
	Cfg *ClientConfig
	RPC rpc.FunctionRpcClient
}

// NewClient returns a new instance of Client
func NewClient(cfg *ClientConfig, conn *grpc.ClientConn) *Client {
	log.Debugf("executing NewClient with config: host %s", cfg.Host)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Cfg: cfg,
		RPC: rpc.NewFunctionRpcClient(conn),
	}
}

// StartEventStream starts listening for messages from the Azure Functions Host
func (client *Client) StartEventStream(ctx context.Context, opts ...grpc.CallOption) error {
	log.Debugf("starting event stream..")

	eventStream, err := client.RPC.EventStream(ctx)
	if err != nil {
		log.Fatalf("Cannot get event stream: %v", err)
		return err
	}

	for {
		message, err := eventStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving stream: %v", err)
			return err
		}

		log.Debugf("received message: %s", message.Content)
	}

	return nil
}

//GetGRPCConnection returns a new grpc connection
func GetGRPCConnection(gothamHost string) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(gothamHost, grpc.WithInsecure()); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", gothamHost, err)
	}
	return conn, nil
}
