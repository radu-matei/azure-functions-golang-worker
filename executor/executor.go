package executor

import (
	"encoding/json"
	"fmt"
	"plugin"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"github.com/radu-matei/azure-functions-golang-worker/runtime"
	"github.com/radu-matei/azure-functions-golang-worker/util"
)

var (
	symbolMap = make(map[string]interface{})
	http      runtime.HTTPRequest
	tt        reflect.Type
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

	symbolMap[request.FunctionId] = symbol

	return nil
}

// ExecuteMethod takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteMethod(request *rpc.InvocationRequest) (response *rpc.InvocationResponse) {

	tt := reflect.TypeOf(symbolMap[request.FunctionId]).Out(0)
	var output = reflect.New(tt)

	switch r := request.TriggerMetadata["req"].Data.(type) {
	case *rpc.TypedData_Http:
		h := util.ConvertToHTTPRequest(r.Http)
		ctx := runtime.Context{
			FunctionID:   request.FunctionId,
			InvocationID: request.InvocationId,
		}

		f := reflect.ValueOf(symbolMap[request.FunctionId])

		y := f.Call([]reflect.Value{reflect.ValueOf(h), reflect.ValueOf(ctx)})
		log.Debugf("Function output types: %v", y)
		output = y[0]
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
