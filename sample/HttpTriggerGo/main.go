package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(request *rpc.RpcHttp) string {
	log.SetLevel(log.DebugLevel)

	msg := fmt.Sprintf("Hello, %s", request.Query["name"])
	log.Debugf(msg)

	return msg
}
