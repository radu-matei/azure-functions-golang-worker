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
	Host      string
	Port      int
	WorkerID  string
	RequestID string
}

// Client that listens for events from the Azure Functions host and executes Golang methods
type Client struct {
	Cfg *ClientConfig
	RPC rpc.FunctionRpcClient
}

// NewClient returns a new instance of Client
func NewClient(cfg *ClientConfig, conn *grpc.ClientConn) *Client {
	log.Debugf("executing NewClient with config: host %s:%d, worker id %s, request id %s",
		cfg.Host, cfg.Port, cfg.WorkerID, cfg.RequestID)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Cfg: cfg,
		RPC: rpc.NewFunctionRpcClient(conn),
	}
}

// StartEventStream starts listening for messages from the Azure Functions Host
func (client *Client) StartEventStream(ctx context.Context, opts ...grpc.CallOption) {
	log.Debugf("starting event stream..")

	eventStream, err := client.RPC.EventStream(ctx)
	if err != nil {
		log.Fatalf("Cannot get event stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			message, err := eventStream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("error receiving stream: %v", err)
			}

			handleStreamingMessage(message, client, eventStream)
		}
	}()

	startStreamingMessage := &rpc.StreamingMessage{
		RequestId: client.Cfg.RequestID,
		Content: &rpc.StreamingMessage_StartStream{
			StartStream: &rpc.StartStream{
				WorkerId: client.Cfg.WorkerID,
			},
		},
	}

	if err := eventStream.Send(startStreamingMessage); err != nil {
		log.Fatalf("Failed to send start streaming request: %v", err)
	}
	log.Debugf("sent start streaming message to host")

	<-waitc

}

//GetGRPCConnection returns a new grpc connection
func GetGRPCConnection(host string) (conn *grpc.ClientConn, err error) {
	log.Debugf("trying to dial %s", host)
	if conn, err = grpc.Dial(host, grpc.WithInsecure()); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", host, err)
	}
	return conn, nil
}
