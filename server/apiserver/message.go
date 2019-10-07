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
	Msg
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Color  string `json:"color"`
	UserID string `json:"userID"`
}

type CreateUserMsg struct {
	Msg
	Id string `json:"id"`
}
