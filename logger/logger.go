package logger

// separate package in order to have private eventStream not visible from azfunc

import (
	"fmt"

	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// Logger exposes the functionality to send logs back to the runtime
type Logger struct {
	invocationID string
	eventStream  rpc.FunctionRpc_EventStreamClient
}

// NewLogger returns a new instance of type Logger to be used in user funcs
func NewLogger(e rpc.FunctionRpc_EventStreamClient, invocationID string) *Logger {
	return &Logger{
		invocationID: invocationID,
		eventStream:  e,
	}
}

// Log sends a log message to the runtime
func (l *Logger) Log(format string, args ...interface{}) error {

	log := &rpc.RpcLog{
		InvocationId: l.invocationID,
		Level:        rpc.RpcLog_Information,
		Message:      fmt.Sprintf(format, args...),
	}

	return l.eventStream.Send(&rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: log,
		},
	})
}

/*
type RpcLog struct {
	InvocationId string        `protobuf:"bytes,1,opt,name=invocation_id,json=invocationId" json:"invocation_id,omitempty"`
	Category     string        `protobuf:"bytes,2,opt,name=category" json:"category,omitempty"`
	Level        RpcLog_Level  `protobuf:"varint,3,opt,name=level,enum=FunctionRpc.RpcLog_Level" json:"level,omitempty"`
	Message      string        `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
	EventId      string        `protobuf:"bytes,5,opt,name=event_id,json=eventId" json:"event_id,omitempty"`
	Exception    *RpcException `protobuf:"bytes,6,opt,name=exception" json:"exception,omitempty"`
	// json serialized property bag, or could use a type scheme like map<string, TypedData>
	Properties string `protobuf:"bytes,7,opt,name=properties" json:"properties,omitempty"`
}
*/
