package executor

import (
	"fmt"
	"plugin"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"github.com/radu-matei/azure-functions-golang-worker/runtime"
	"github.com/radu-matei/azure-functions-golang-worker/util"
)

var (
	functionMap = make(map[string]interface{})
	http        *runtime.HTTPRequest
)

// LoadMethod takes a .so object from the function's bin directory and loads it
func LoadMethod(request *rpc.FunctionLoadRequest) error {

	path := fmt.Sprintf("%s/bin/%s.so", request.Metadata.Directory, request.Metadata.Name)
	plugin, err := plugin.Open(path)
	if err != nil {
		log.Debugf("cannot get .so object from path %s: %v", path, err)
		return err
	}

	symbol, err := plugin.Lookup(request.Metadata.EntryPoint)
	if err != nil {
		log.Debugf("cannot look up symbol for entrypoint function %s: %v", request.Metadata.EntryPoint, err)
	}

	t := reflect.TypeOf(symbol)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("symbol is not func, but %v", t.Kind())
	}

	triggerType := t.In(0)
	if triggerType != reflect.TypeOf(http) {
		return fmt.Errorf("first argument not http request but %v", triggerType)
	}

	log.Debugf("entrypoint function type: %v, signature: %v", t.Kind(), t)

	f, ok := symbol.(func(*runtime.HTTPRequest, *runtime.Context) []byte)
	if !ok {
		log.Debug("incorrect function signature")
		return err
	}

	functionMap[request.FunctionId] = f

	return nil
}

// ExecuteMethod takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteMethod(request *rpc.InvocationRequest) (response *rpc.InvocationResponse) {
	var output []byte
	switch r := request.TriggerMetadata["req"].Data.(type) {
	case *rpc.TypedData_Http:
		h := util.ConvertToHTTPRequest(r.Http)
		ctx := &runtime.Context{
			FunctionName: request.FunctionId,
			InvocationID: request.InvocationId,
		}
		output = functionMap[request.FunctionId].(func(*runtime.HTTPRequest, *runtime.Context) []byte)(h, ctx)
	}

	log.Debugf("received output: %v", output)

	invocationResponse := &rpc.InvocationResponse{
		InvocationId: request.InvocationId,
		Result: &rpc.StatusResult{
			Status: rpc.StatusResult_Success,
		},
		ReturnValue: &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: string(output),
			},
		},
	}

	return invocationResponse
}
