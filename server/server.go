package main

import (
	"os"

	"github.com/rgreen312/owlplace/server/apiserver"
	log "github.com/sirupsen/logrus"
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
	server, err := apiserver.NewApiServer(os.Getenv("MY_POD_IP"))
	if err != nil {
		log.Fatal(err)
	}
	server.ListenAndServe()
}
