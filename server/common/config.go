package common

import (
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
