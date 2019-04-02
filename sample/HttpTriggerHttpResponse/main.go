package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx *azfunc.Context, req *http.Request) (resp *http.Response) {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)

	name := req.URL.Query().Get("name")
	if name == "" {
		name = "anonymous"
	}

	respBody := fmt.Sprintf("Hello, %s!", name)

	resp = &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
		ContentLength: int64(len(respBody)),
		Request:       req,
		Header:        make(http.Header, 0),
	}
	return
}
