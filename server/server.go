package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/common"
)

func mapKeys(m map[int]*common.ServerConfig) []int {
	keys := make([]int, len(m))
	ptr := 0
	for key := range m {
		keys[ptr] = key
		ptr++
	}
	return keys
}

func initLogging() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	configFile := flag.String("config", "owlplace-docker.json", "Configuration file that contains a list of servers.")
	nodeID := flag.Int("nodeid", 1, "Index in the configuration file that represents this node.")

	flag.Parse()

	// Initialize logrus
	initLogging()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", *configFile)
	}

	var servers map[int]*common.ServerConfig

	err = json.Unmarshal([]byte(file), &servers)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %s\n%s", *configFile, err)
	}

	if _, ok := servers[*nodeID]; !ok {
		log.Fatalf("Requested nodeID is not found in the configuration file.  Valid nodeIDs: %+v", mapKeys(servers))
	}

	log.WithFields(log.Fields{
		"server config": servers,
		"nodeID":        *nodeID,
	}).Debug("joining dragonboat cluster")

	server, err := apiserver.NewApiServer(servers, *nodeID)
	if err != nil {
		log.WithFields(log.Fields{
			"server config": servers,
			"nodeID":        *nodeID,
		}).Fatal(err)
	}
	server.ListenAndServe()
}
