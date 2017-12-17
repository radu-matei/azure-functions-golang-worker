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

// ExecuteFunc takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteFunc(req *rpc.InvocationRequest) (response *rpc.InvocationResponse) {

	status := rpc.StatusResult_Success

	f, ok := loader.LoadedFuncs[req.FunctionId]
	if !ok {
		log.Debugf("function with functionID %v not loaded", req.FunctionId)
		status = rpc.StatusResult_Failure
	}
	params, err := getFuncParams(req, f)
	if err != nil {
		log.Debugf("cannot get params from request: %v", err)
		status = rpc.StatusResult_Failure
	}

	log.Debugf("params: %v", params)

	output := f.Func.Call(params)[0]

	b, err := json.Marshal(output.Interface())
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	return &rpc.InvocationResponse{
		InvocationId: req.InvocationId,
		Result: &rpc.StatusResult{
			Status: status,
		},
		ReturnValue: &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: string(b),
			},
		},
	}
}

// get final slice with params to call the function
func getFuncParams(req *rpc.InvocationRequest, f *azfunc.Func) ([]reflect.Value, error) {

	args := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		binding, ok := f.Bindings[input.Name]
		if ok {
			v, err := getValueFromBinding(input, binding)

			// TODO - investigate returning error from this function
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, err
			}
			args[input.Name] = v
		} else {
			return nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	ctx := &azfunc.Context{
		FunctionID:   req.FunctionId,
		InvocationID: req.InvocationId,
	}

	params := []reflect.Value{}
	for arg, t := range f.NamedInArgs {
		p, ok := args[arg]
		if ok {
			params = append(params, p)
		} else if t == reflect.TypeOf((*azfunc.Context)(nil)) {
			params = append(params, reflect.ValueOf(ctx))
		} else {
			return nil, fmt.Errorf("named argument not found")
		}
	}

	return params, nil
}

// TODO - add here cases for all bindings supported by Azure Functions
func getValueFromBinding(input *rpc.ParameterBinding, binding *rpc.BindingInfo) (reflect.Value, error) {

	switch binding.Type {
	case azfunc.HTTPTriggerType:
		switch r := input.Data.Data.(type) {
		case *rpc.TypedData_Http:
			h, err := util.ConvertToHTTPRequest(r.Http)
			if err != nil {
				return reflect.New(nil), err
			}
			return reflect.ValueOf(h), nil
		}

	case azfunc.BlobBindingType:
		switch d := input.Data.Data.(type) {
		case *rpc.TypedData_String_:
			b, err := util.ConvertToBlobInput(d)
			if err != nil {
				return reflect.New(nil), err
			}

			return reflect.ValueOf(b), nil
		}
	}
	return reflect.New(nil), fmt.Errorf("cannot handle binding %v", binding.Type)
}
