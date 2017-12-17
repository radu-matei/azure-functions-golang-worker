package util

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// ConvertToHTTPRequest returns a formatted HTTPRequest from an rpc.HttpTrigger
func ConvertToHTTPRequest(r *rpc.RpcHttp) (*azfunc.HTTPRequest, error) {

	if r == nil {
		return nil, fmt.Errorf("cannot convert nil request")
	}

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
		return req, nil
	}

	switch d := r.Body.Data.(type) {
	case *rpc.TypedData_String_:
		req.Body = ioutil.NopCloser(bytes.NewBufferString(d.String_))
	}

	return req, nil
}

// ConvertToBlobInput returns a formatted BlobInput from an rpc.TypedData_String (as blob inputs are for now)
func ConvertToBlobInput(d *rpc.TypedData_String_) (*azfunc.BlobInput, error) {
	if d == nil {
		return nil, fmt.Errorf("cannot convert nil blob input")
	}

	return &azfunc.BlobInput{
		Data: d.String_,
	}, nil
}
