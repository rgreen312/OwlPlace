package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/common"
)

func mapKeys(m map[int]*common.ServerConfig) []int {
	keys := make([]int, len(m))
	ptr := 0
	for key, _ := range m {
		keys[ptr] = key
		ptr++
	}
	return keys
}

func main() {
	configFile := flag.String("config", "owlplace-docker.json", "Configuration file that contains a list of servers.")
	nodeId := flag.Int("nodeid", 1, "Index in the configuration file that represents this node.")

	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", *configFile)
	}

	var servers map[int]*common.ServerConfig

	err = json.Unmarshal([]byte(file), &servers)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %s\n%s", *configFile, err)
	}

	if _, ok := servers[*nodeId]; !ok {
		log.Fatalf("Requested nodeId is not found in the configuration file.  Valid nodeIds: %+v", mapKeys(servers))
	}

	log.Printf("Joining Dragonboat cluster with configuration:\n%+v", servers)

	server := apiserver.NewApiServer(servers, *nodeId)
	server.SetupRoutes()
	// http.ListenAndServe(fmt.Sprintf(":%d", api.port), nil)
}
