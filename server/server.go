package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	log "github.com/sirupsen/logrus"
)

const (
	NAMESPACE        = "dev"
	OWLPLACE_NODEID  = "OWLPLACE_NODEID"
	OWLPLACE_ADDRESS = "OWLPLACE_ADDRESS"
)

func initLogging() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	membershipFile := flag.String("members", "", "Membership file containing a list of servers belonging to the cluster.")
	flag.Parse()

	// Initialize logrus
	initLogging()

	// Recover service address through environment variables:
	address := os.Getenv(OWLPLACE_ADDRESS)
	if address == "" {
		log.Fatalf("Provide owlplace address via environment variable '%s'", OWLPLACE_ADDRESS)
	}

	var membershipProvider consensus.MembershipProvider
	var nodeID uint64
	var err error

	// This indicates we'd like to use k8s as a discovery service.
	if *membershipFile == "" {
		nodeID, err = common.IPToNodeId(address)
        if err != nil {
            log.Fatalf("Error constructing nodeID from address: '%s'\nerr: %s", address, err.Error())
        }
		membershipProvider, err = consensus.NewKubernetesMembershipProvider(NAMESPACE)
	} else {
		// Recover nodeID and service address through environment variables:
		stringNodeID := os.Getenv(OWLPLACE_NODEID)
        nodeIDInt, err := strconv.Atoi(stringNodeID)
		nodeID = uint64(nodeIDInt)
		if err != nil {
			log.Fatalf("Invalid nodeID: '%s', provide via environment variable '%s'", stringNodeID, OWLPLACE_NODEID)
		}
		membershipProvider, err = consensus.StaticMembershipFromFile(*membershipFile)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Start server
	server, err := apiserver.NewApiServer(uint64(nodeID), address, membershipProvider)
	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
