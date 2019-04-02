package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// ConvertToHTTPRequest returns a formatted HTTPRequest from an rpc.HttpTrigger
func ConvertToHTTPRequest(r *rpc.RpcHttp) (*http.Request, error) {

	if r == nil {
		return nil, fmt.Errorf("cannot convert nil request")
	}

	var body io.Reader
	if r.Body != nil {
		switch r := r.Body.Data.(type) {
		case *rpc.TypedData_String_:
			body = ioutil.NopCloser(bytes.NewBufferString(r.String_))
		}
	}

	req, err := http.NewRequest(r.GetMethod(), r.GetUrl(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range r.GetHeaders() {
		req.Header.Set(key, value)
	}

	return req, nil
}

// ConvertToBlobInput returns a formatted BlobInput from an rpc.TypedData_String (as blob inputs are for now)
func ConvertToBlobInput(d *rpc.TypedData_String_) (*azfunc.Blob, error) {
	if d == nil {
		return nil, fmt.Errorf("cannot convert nil blob input")
	}

	return &azfunc.Blob{
		Data: d.String_,
	}, nil
}
