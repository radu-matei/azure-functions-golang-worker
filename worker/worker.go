package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

var (
	functionMap = make(map[string]*rpc.RpcFunctionMetadata)
)

func handleStreamingMessage(message *rpc.StreamingMessage, client *Client, eventStream rpc.FunctionRpc_EventStreamClient) {
	switch m := message.Content.(type) {

	case *rpc.StreamingMessage_WorkerInitRequest:
		handleWorkerInitRequest(message.RequestId, m, client, eventStream)

	case *rpc.StreamingMessage_FunctionLoadRequest:
		handleFunctionLoadRequest(message.RequestId, m, client, eventStream)
	}
}

func handleWorkerInitRequest(requestID string, message *rpc.StreamingMessage_WorkerInitRequest, client *Client, eventStream rpc.FunctionRpc_EventStreamClient) {
	log.Debugf("received worker init request with host version %s ", message.WorkerInitRequest.HostVersion)

	workerInitResponse := &rpc.StreamingMessage{
		RequestId: requestID,
		Content: &rpc.StreamingMessage_WorkerInitResponse{
			WorkerInitResponse: &rpc.WorkerInitResponse{
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
			},
		},
	}

	if err := eventStream.Send(workerInitResponse); err != nil {
		log.Fatalf("Failed to send worker init response: %v", err)
	}
	log.Debugf("sent start worker init response : %v", workerInitResponse)
}

func handleFunctionLoadRequest(requestID string, message *rpc.StreamingMessage_FunctionLoadRequest, client *Client, eventStream rpc.FunctionRpc_EventStreamClient) {
	functionMap[message.FunctionLoadRequest.FunctionId] = message.FunctionLoadRequest.Metadata
	log.Debugf("added function to map: %s, %s", message.FunctionLoadRequest.FunctionId, functionMap[message.FunctionLoadRequest.FunctionId])

	functionLoadResponse := &rpc.StreamingMessage{
		RequestId: requestID,
		Content: &rpc.StreamingMessage_FunctionLoadResponse{
			FunctionLoadResponse: &rpc.FunctionLoadResponse{
				FunctionId: message.FunctionLoadRequest.FunctionId,
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
			},
		},
	}

	if err := eventStream.Send(functionLoadResponse); err != nil {
		log.Fatalf("Failed to send function load response: %v", err)
	}
	log.Debugf("sent function load response: %v", functionLoadResponse)
}
