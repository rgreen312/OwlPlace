package common

type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	WebsocketPort int    `json:"websocket_port"`
	ConsensusPort int    `json:"consensus_port"`
}
