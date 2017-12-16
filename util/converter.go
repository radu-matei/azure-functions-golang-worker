package util

import (
	"bytes"
	"io/ioutil"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// ConvertToHTTPRequest returns a formatted HTTPRequest from an rpc.HttpTrigger
func ConvertToHTTPRequest(r *rpc.RpcHttp) *azfunc.HTTPRequest {

	req := &azfunc.HTTPRequest{
		Method:     r.Method,
		URL:        r.Url,
		Headers:    r.Headers,
		Params:     r.Params,
		StatusCode: r.StatusCode,
		Query:      r.Query,
		IsRaw:      r.IsRaw,
	}

	if r.Body == nil {
		return req
	}

	switch d := r.Body.Data.(type) {
	case *rpc.TypedData_String_:
		req.Body = ioutil.NopCloser(bytes.NewBufferString(d.String_))
	}

	return req
}
