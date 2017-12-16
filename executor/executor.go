package executor

import (
	"encoding/json"
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
	"github.com/radu-matei/azure-functions-golang-worker/loader"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"github.com/radu-matei/azure-functions-golang-worker/util"
)

// ExecuteFunction takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteFunction(request *rpc.InvocationRequest) (response *rpc.InvocationResponse) {
	var output reflect.Value

	s := len(loader.LoadedFuncs[request.FunctionId].Out)
	out := make([]reflect.Value, s)
	for i := 0; i < s; i++ {
		out[i] = reflect.New(loader.LoadedFuncs[request.FunctionId].Out[i])
	}

	switch r := request.TriggerMetadata["req"].Data.(type) {
	case *rpc.TypedData_Http:
		h := util.ConvertToHTTPRequest(r.Http)
		ctx := &azfunc.Context{
			FunctionID:   request.FunctionId,
			InvocationID: request.InvocationId,
		}

		out = loader.LoadedFuncs[request.FunctionId].Func.Call([]reflect.Value{reflect.ValueOf(h), reflect.ValueOf(ctx)})
		output = out[0]

		log.Debugf("http request: %v", h)
	}

	b, err := json.Marshal(output.Interface())
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	log.Debugf("received output: %v", output)
	log.Debugf("encoded output: %v", string(b))

	invocationResponse := &rpc.InvocationResponse{
		InvocationId: request.InvocationId,
		Result: &rpc.StatusResult{
			Status: rpc.StatusResult_Success,
		},
		ReturnValue: &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: string(b),
			},
		},
	}

	return invocationResponse
}

// NewExecuteFunc is a new, improved version of ExecuteFunction
func NewExecuteFunc(req *rpc.InvocationRequest) (response *rpc.InvocationResponse, err error) {
	f, ok := loader.LoadedFuncs[req.FunctionId]
	if !ok {
		return nil, fmt.Errorf("function with functionID %v not loaded", req.FunctionId)
	}
	_ = paramsFromReq(req, f)

	return nil, nil
}

func paramsFromReq(req *rpc.InvocationRequest, f *azfunc.Func) map[string]reflect.Value {

	return nil
}
