package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/runtime"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(request *runtime.HTTPRequest, ctx *runtime.Context) []byte {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s", ctx.FunctionID, ctx.InvocationID)

	u := User{
		Name:          request.Query["name"],
		GeneratedName: fmt.Sprintf("%s-azfunc", request.Query["name"]),
	}

	b, err := json.Marshal(u)
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	if request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Debugf("cannot read body")
		}
		log.Debugf("request body: %v", string(body))
	}

	return b
}

type User struct {
	Name          string
	GeneratedName string
}
