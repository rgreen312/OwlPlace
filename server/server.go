package main

import (
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/rgreen312/owlplace/server/apiserver"
)


func initLogging() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	// Initialize logrus
	initLogging()

	// Start server
	server, _ := apiserver.NewApiServer(os.Getenv("MY_POD_IP"))
	server.ListenAndServe()
}
