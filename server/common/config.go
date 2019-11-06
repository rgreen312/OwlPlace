package common

import (
	"strings"
	"fmt"
	"strconv"
)


type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	ConsensusPort int    `json:"consensus_port"`
}


func IPToNodeId(ip_address string) int {
	// This function maps ip addresses to node-ids
	components := strings.Split(ip_address, ".")
	combined := fmt.Sprintf("%03s%03s", components[2], components[3])
	node_id, _ := strconv.Atoi(combined)
	return node_id
}
const (
	ApiPort int = 3001
)

