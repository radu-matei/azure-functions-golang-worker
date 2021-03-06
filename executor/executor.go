package executor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
	"github.com/radu-matei/azure-functions-golang-worker/loader"
	"github.com/radu-matei/azure-functions-golang-worker/logger"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"github.com/radu-matei/azure-functions-golang-worker/util"
)

// ExecuteFunc takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteFunc(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) (response *rpc.InvocationResponse) {

	log.Debugf("\n\n\nInvocation Request: %v", req)

	response = &rpc.InvocationResponse{
		InvocationId: req.InvocationId,
		Result: &rpc.StatusResult{
			Status: rpc.StatusResult_Failure,
		},
	}

	f, ok := loader.LoadedFuncs[req.FunctionId]
	if !ok {
		log.Debugf("function with functionID %v not loaded", req.FunctionId)
		return
	}
	params, outBindings, err := getFinalParams(req, f, eventStream)
	if err != nil {
		log.Debugf("cannot get params from request: %v", err)
		return
	}

	log.Debugf("params: %v", params)
	log.Debugf("out bindings: %v", outBindings)

	output := f.Func.Call(params)[0]

	switch v := output.Interface().(type) {
	case *http.Response:
		resp, err := encodeResponse(v)
		if err != nil {
			log.Debugf("cannot encode response: %v", err)
			return
		}
		response.ReturnValue = &rpc.TypedData{
			Data: &rpc.TypedData_Http{
				Http: resp,
			},
		}
		response.OutputData = []*rpc.ParameterBinding{
			{
				Name: "handler",
				Data: &rpc.TypedData{
					Data: &rpc.TypedData_Http{
						Http: resp,
					},
				},
			},
		}
		response.Result.Status = rpc.StatusResult_Success
		return
	default:
		b, err := json.Marshal(output.Interface())
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
			return
		}
		response.ReturnValue = &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: string(b),
			},
		}
	}

	outputData := make([]*rpc.ParameterBinding, len(outBindings))
	i := 0
	for k, v := range outBindings {
		b, err := json.Marshal(v.Interface())
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
			return
		}

		outputData[i] = &rpc.ParameterBinding{
			Name: k,
			Data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: string(b),
				},
			},
		}
	}
	response.OutputData = outputData
	response.Result.Status = rpc.StatusResult_Success
	return
}

func getFinalParams(req *rpc.InvocationRequest, f *azfunc.Func, eventStream rpc.FunctionRpc_EventStreamClient) ([]reflect.Value, map[string]reflect.Value, error) {
	args := make(map[string]reflect.Value)
	outBindings := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		binding, ok := f.Bindings[input.Name]
		if ok {
			v, err := getValueFromBinding(input, binding)
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, nil, err
			}
			args[input.Name] = v
		} else {
			return nil, nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	ctx := &azfunc.Context{
		FunctionID:   req.FunctionId,
		InvocationID: req.InvocationId,
		Logger:       logger.NewLogger(eventStream, req.InvocationId),
	}

	log.Debugf("args map: %v", args)

	params := make([]reflect.Value, len(f.NamedInArgs))
	i := 0
	for _, v := range f.NamedInArgs {
		p, ok := args[v.Name]
		if ok {
			params[i] = p
			i++
		} else if v.Type == reflect.TypeOf((*azfunc.Context)(nil)) {
			params[i] = reflect.ValueOf(ctx)
			i++
		} else {
			b, ok := f.Bindings[v.Name]
			if ok {
				if b.Direction == rpc.BindingInfo_out {
					o, err := getOutBinding(b)
					if err != nil {
						return nil, nil, fmt.Errorf("cannot get out binding %s: %v", v.Name, err)
					}

					params[i] = o
					outBindings[v.Name] = o
					i++
				}
			}
		}
	}

	return params, outBindings, nil
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

func getOutBinding(b *rpc.BindingInfo) (reflect.Value, error) {
	switch b.Type {
	case azfunc.BlobBindingType:
		b := &azfunc.Blob{
			Data: "",
		}
		return reflect.ValueOf(b), nil
	case azfunc.HTTPBindingType:
		return reflect.ValueOf(&http.Response{}), nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot handle binding %v", b.Type)
	}
}

// encodeResponse returns a protobuf Http type from a *http.Response
func encodeResponse(r *http.Response) (*rpc.RpcHttp, error) {
	resp := &rpc.RpcHttp{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	resp.Body = &rpc.TypedData{
		Data: &rpc.TypedData_String_{
			String_: string(body),
		},
	}
	// For now return body as a string also
	resp.Body = &rpc.TypedData{
		Data: &rpc.TypedData_String_{
			String_: string(body),
		},
	}
	resp.Headers = make(map[string]string, len(r.Header))
	for key, value := range r.Header {
		resp.Headers[key] = strings.Join(value, ",")
	}
	resp.StatusCode = strconv.Itoa(r.StatusCode)
	return resp, nil
}
