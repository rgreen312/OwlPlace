package apiserver

type MsgType uint8

// Well defined Message types
const (
	Open      MsgType = 0
	DrawPixel MsgType = 1
	LoginUser MsgType = 2
	Image     MsgType = 4
	Testing   MsgType = 5
	Close     MsgType = 9
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

type LoginUserMsg struct {
	Type MsgType `json:"type"`
	Id   string  `json:"id"`
}

type ImageMsg struct {
	Type         MsgType `json:"type"`
	FormatString string  `json:"formatString"`
}

type TestingMsg struct {
	Type MsgType `json:"type"`
	Msg  string  `json:"msg"`
}
