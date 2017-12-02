package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

func handleStreamingMessage(message *rpc.StreamingMessage) {
	switch m := message.Content.(type) {

	case *rpc.StreamingMessage_WorkerInitRequest:
		handleWorkerInitRequest(m)
	}
}

func handleWorkerInitRequest(message *rpc.StreamingMessage_WorkerInitRequest) {
	log.Debugf("received worker init request with host version %s ", message.WorkerInitRequest.HostVersion)
}
