package common

import (
	"time"
)

const (
	TimeFormat = time.RFC3339
	//Cooldown   = time.Duration(5 * time.Minute)
	Cooldown  = time.Duration(5 * time.Second)
	AlphaMask = 255
)

type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	ConsensusPort int    `json:"consensus_port"`
}
