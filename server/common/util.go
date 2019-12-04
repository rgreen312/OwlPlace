package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SetAddressPort replaces the port at the end of an address.
func SetAddressPort(address string, newPort int) string {
	banana := strings.Split(address, ":")
	// No need to check the banana's length, if no port is attached it'll just
	// contain the hostname.
	return fmt.Sprintf("%s:%d", banana[0], newPort)
}

// IPToNodeId maps IP Addresses to nodeIDs
func IPToNodeId(ip_address string) (uint64, error) {
	components := strings.Split(ip_address, ".")
	if len(components) == 4 {
		combined := fmt.Sprintf("%03s%03s", components[2], components[3])
		node_id, _ := strconv.Atoi(combined)
		return uint64(node_id), nil
	} else {
		return 0, os.ErrExist
	}
}
