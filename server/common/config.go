package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TimeFormat = time.RFC3339
	//Cooldown   = time.Duration(5 * time.Minute)
	Cooldown      = time.Duration(15 * time.Second)
	AlphaMask     = 255
	ApiPort       = 3001
	ConsensusPort = 3010
)

func IPToNodeId(ip_address string) uint64 {
	// This function maps ip addresses to node-ids
	components := strings.Split(ip_address, ".")
	combined := fmt.Sprintf("%03s%03s", components[2], components[3])
	node_id, _ := strconv.Atoi(combined)
	return uint64(node_id)
}
