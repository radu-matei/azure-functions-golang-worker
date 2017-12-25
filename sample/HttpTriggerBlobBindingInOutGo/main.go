package main

import (
	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *azfunc.HTTPRequest, inBlob *azfunc.Blob, outBlob *azfunc.Blob, ctx *azfunc.Context) BlobData {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)

	d := BlobData{
		Name: req.Query["name"],
		Data: inBlob.Data,
	}

	outBlob.Data = "Leeeet's hope this doesn't miserably fail..."

	return d
}

// BlobData mocks any struct (or pointer to struct) you might want to return
type BlobData struct {
	Name string
	Data string
}
