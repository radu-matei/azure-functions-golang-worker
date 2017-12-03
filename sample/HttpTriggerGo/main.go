package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Run is the entrypoint to our Go Azure Function
func Run() {
	fmt.Println("Hello, Azure Functions!")
	log.SetLevel(log.DebugLevel)
	log.Debugf("logging in Azure Function using external Go dependency")
}
