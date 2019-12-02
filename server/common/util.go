package common

import (
	"fmt"
	"strconv"
	"strings"
)

// ReplacePort replaces the port at the end of an address.
func ReplacePort(address string, newPort int) string {
	banana := strings.Split(address, ":")
	return fmt.Sprintf("%s:%d", banana[0], newPort)
}

func IPToNodeId(ip_address string) uint64 {
	// This function maps ip addresses to node-ids
	components := strings.Split(ip_address, ".")
	combined := fmt.Sprintf("%03s%03s", components[2], components[3])
	node_id, _ := strconv.Atoi(combined)
	return uint64(node_id)
}
