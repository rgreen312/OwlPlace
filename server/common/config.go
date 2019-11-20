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

// Well defined Message types
const (
	Error             MsgType = -1
	Open              MsgType = 0
	DrawPixel         MsgType = 1
	LoginUser         MsgType = 2
	ChangeClientPixel MsgType = 3
	Image             MsgType = 4
	Testing           MsgType = 5
	DrawResponse      MsgType = 6
	Close             MsgType = 9
	VerificationFail  MsgType = 10
	CreateUser        MsgType = 11
)


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
