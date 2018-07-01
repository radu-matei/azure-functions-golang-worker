package main

import (
	"strconv"

	"github.com/radu-matei/azure-functions-golang-worker/azfunc"
)

func fibonacci(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *azfunc.HTTPRequest, ctx *azfunc.Context) int {
	n, err := strconv.Atoi(req.Query["number"])
	if err != nil {
		ctx.Logger.Log("cannot convert query string to integer: %v", err)
	}

	return fibonacci(n)
}
