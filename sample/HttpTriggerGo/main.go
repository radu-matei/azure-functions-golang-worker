package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(request *rpc.RpcHttp) {
	log.SetLevel(log.DebugLevel)
	log.Debugf("received function invocation request: %v", request)

}

func main() {
	// this will be built using -buildmode=pluign, so there is no need for main
	// only here to silence VSCode errors - don't use, will be removed
}
