package common

import (
	"time"
	"strings"
	"fmt"
	"strconv"
	"os"
)

const (
	TimeFormat = time.RFC3339
	//Cooldown   = time.Duration(5 * time.Minute)
	Cooldown  = time.Duration(15 * time.Second)
	AlphaMask = 255
	ApiPort = 3001
	ConsensusPort = 3010
)

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
	UserLoginResponse MsgType = 11
)


func IPToNodeId(ip_address string) (uint64, error) {
	// This function maps ip addresses to node-ids
	components := strings.Split(ip_address, ".")
	if(len(components) == 4){
		combined := fmt.Sprintf("%03s%03s", components[2], components[3])
		node_id, _ := strconv.Atoi(combined)
		return uint64(node_id), nil
	} else {
		return 0, os.ErrExist
	}
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
