package main

import (
	"fmt"
	"net/http"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx *azfunc.Context, req *http.Request) User {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)

	name := req.URL.Query().Get("name")
	if name == "" {
		name = "anonymous"
	}

	u := User{
		Name:          name,
		GeneratedName: fmt.Sprintf("%s-azfunc", name),
	}

	return u
}

// User mocks any struct (or pointer to struct) you might want to return
type User struct {
	Name          string
	GeneratedName string
}
