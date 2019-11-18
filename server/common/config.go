package common

import (
	"time"
)

type MsgType int8

const (
	TimeFormat = time.RFC3339
)

type ServerConfig struct {
	Hostname      string `json:"hostname"`
	ApiPort       int    `json:"api_port"`
	ConsensusPort int    `json:"consensus_port"`
}


/*
	This message type is intended to be sent from
	the server to the client, notifying the user
	that a pixel was drawn by another user.
*/
type ChangeClientPixelMsg struct {
	Type   MsgType `json:"type"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	R      int     `json:"r"`
	G      int     `json:"g"`
	B      int     `json:"b"`
	UserID string  `json:"userID"`
}
