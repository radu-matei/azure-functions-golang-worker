package main

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(request *rpc.RpcHttp) []byte {
	log.SetLevel(log.DebugLevel)

	u := User{
		Name:          request.Query["name"],
		GeneratedName: fmt.Sprintf("%s-azfunc", request.Query["name"]),
	}

	log.Debugf("user: %v", u)

	b, err := json.Marshal(u)
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	return b

}

type User struct {
	Name          string
	GeneratedName string
}
