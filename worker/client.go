package worker

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// ClientConfig contains all necessary configuration to connect to the Azure Functions Host
type ClientConfig struct {
	Host             string
	Port             int
	WorkerID         string
	RequestID        string
	MaxMessageLength int
}

// Client that listens for events from the Azure Functions host and executes Golang methods
type Client struct {
	Cfg  *ClientConfig
	conn *grpc.ClientConn
}

// New returns a new instance of Client
func New(cfg *ClientConfig) *Client {
	log.Debugf("executing New with config: host %s:%d, worker id %s, request id %s",
		cfg.Host, cfg.Port, cfg.WorkerID, cfg.RequestID)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Cfg: cfg,
	}
}

// StartEventStream starts listening for messages from the Azure Functions Host
func (client *Client) StartEventStream(ctx context.Context, opts ...grpc.CallOption) {
	log.Debugf("starting event stream..")

	eventStream, err := rpc.NewFunctionRpcClient(client.conn).EventStream(ctx)
	if err != nil {
		log.Fatalf("cannot get event stream: %v", err)
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
		log.Fatalf("failed to send start streaming request: %v", err)
	}
	log.Debugf("sent start streaming message to host")

	<-waitc

}

// Connect tries to establish a grpc connection with the server
func (client *Client) Connect(opts ...grpc.DialOption) (err error) {
	log.Debugf("attempting to start grpc connection to server %s:%d with worker id %s and request id %s", client.Cfg.Host, client.Cfg.Port, client.Cfg.WorkerID, client.Cfg.RequestID)
	opts = append(opts, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(client.Cfg.MaxMessageLength)))

	conn, err := getGRPCConnection(fmt.Sprintf("%s:%d", client.Cfg.Host, client.Cfg.Port), opts)
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
		return
	}

	client.conn = conn
	log.Debugf("started grpc connection...")
	return
}

// Disconnect closes the connection to the server
func (client *Client) Disconnect() error {
	return client.conn.Close()
}

//getGRPCConnection returns a new grpc connection
func getGRPCConnection(host string, opts []grpc.DialOption) (conn *grpc.ClientConn, err error) {
	log.Debugf("trying to dial %s", host)
	if conn, err = grpc.Dial(host, opts...); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", host, err)
	}
	return conn, nil
}
