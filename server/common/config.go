package common

import (
	"time"
)

const (
	TimeFormat = time.RFC3339
)

type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	ConsensusPort int    `json:"consensus_port"`
}
