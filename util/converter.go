package util

import (
	"bytes"
	"io/ioutil"

	"github.com/radu-matei/azure-functions-golang-worker/rpc"
	"github.com/radu-matei/azure-functions-golang-worker/runtime"
)

// ConvertToHTTPRequest returns a formatted HTTPRequest from an rpc.HttpTrigger
func ConvertToHTTPRequest(r *rpc.RpcHttp) *runtime.HTTPRequest {

	req := &runtime.HTTPRequest{
		Method:     r.Method,
		URL:        r.Url,
		Headers:    r.Headers,
		Params:     r.Params,
		StatusCode: r.StatusCode,
		Query:      r.Query,
		IsRaw:      r.IsRaw,
	}

	switch d := r.Body.Data.(type) {
	case *rpc.TypedData_String_:
		req.Body = ioutil.NopCloser(bytes.NewBufferString(d.String_))
	}

	return req
}
