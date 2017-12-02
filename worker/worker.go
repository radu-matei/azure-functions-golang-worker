package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

func handleStreamingMessage(message *rpc.StreamingMessage, client *Client, eventStream rpc.FunctionRpc_EventStreamClient) {
	switch m := message.Content.(type) {

	case *rpc.StreamingMessage_WorkerInitRequest:
		handleWorkerInitRequest(m, client, eventStream)
	}
}

func handleWorkerInitRequest(message *rpc.StreamingMessage_WorkerInitRequest, client *Client, eventStream rpc.FunctionRpc_EventStreamClient) {
	log.Debugf("received worker init request with host version %s ", message.WorkerInitRequest.HostVersion)

	workerInitResponse := &rpc.StreamingMessage{
		RequestId: client.Cfg.RequestID,
		Content: &rpc.StreamingMessage_WorkerInitResponse{
			WorkerInitResponse: &rpc.WorkerInitResponse{
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
			},
		},
	}

	if err := eventStream.Send(workerInitResponse); err != nil {
		log.Fatalf("Failed to worker init response: %v", err)
	}
	log.Debugf("sent start worker init response")
}
