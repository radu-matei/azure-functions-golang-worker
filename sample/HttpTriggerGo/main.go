package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/runtime"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req runtime.HTTPRequest, ctx runtime.Context) User {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s", ctx.FunctionID, ctx.InvocationID)

	u := User{
		Name:          req.Query["name"],
		GeneratedName: fmt.Sprintf("%s-azfunc", req.Query["name"]),
	}

	return u
}

// User mocks any struct (or pointer to struct) you might want to return
type User struct {
	Name          string
	GeneratedName string
}
