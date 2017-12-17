package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *azfunc.HTTPRequest, myBlob *azfunc.BlobInput, ctx *azfunc.Context) BlobData {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s", ctx.FunctionID, ctx.InvocationID)

	d := BlobData{
		Name: req.Query["name"],
		Data: myBlob.Data,
	}

	return d
}

// BlobData mocks any struct (or pointer to struct) you might want to return
type BlobData struct {
	Name string
	Data string
}
