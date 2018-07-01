package main

import (
	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *azfunc.HTTPRequest, outBlob *azfunc.Blob, ctx *azfunc.Context) RequestInfo {

	r := RequestInfo{
		UA:   req.Headers["user-agent"],
		Host: req.Headers["host"],
	}

	outBlob.Data = r.UA

	return r
}

// RequestInfo contains information about the request
type RequestInfo struct {
	UA   string
	Host string
}
