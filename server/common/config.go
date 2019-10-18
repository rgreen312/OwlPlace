package common

type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	ConsensusPort int    `json:"consensus_port"`
}
