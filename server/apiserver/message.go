package apiserver

type MsgType uint8

const (
	DrawPixel  MsgType = 0
	CreateUser MsgType = 1
)

type Msg struct {
	Type MsgType `json:"type"`
}

type DrawPixelMsg struct {
	Type   MsgType `json:"type"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	R      int     `json:"r"`
	G      int     `json:"g"`
	B      int     `json:"b"`
	UserID string  `json:"userID"`
}

type CreateUserMsg struct {
	Type MsgType `json:"type"`
	Id   string  `json:"id"`
}
