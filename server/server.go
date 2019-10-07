package main

import (
	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/consensus"
)

func main() {
	// Make the backend channel that the api server and consensus module communicate with
	api_to_backend_channel := make(chan consensus.BackendMessage)
	backend_to_api_channel := make(chan consensus.ConsensusMessage)
	// Start API listening asynchronously (TODO: pass in channel)
	server := apiserver.NewApiServer(api_to_backend_channel, backend_to_api_channel)
	go server.ListenAndServe()

	// Start consensus service
	consensus.MainConsensus(api_to_backend_channel, backend_to_api_channel)

}
