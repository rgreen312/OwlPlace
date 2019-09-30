package main

import (

	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/consensus"
)


func main() {
	// Make the backend channel that the api server and consensus module communicate with
	backend_channel := make(chan consensus.BackendMessage)
	// Start API listening asynchronously (TODO: pass in channel)
	server := apiserver.NewApiServer(backend_channel)
	go server.ListenAndServe()

	// Start consensus service
	consensus.MainConsensus(backend_channel)

}
