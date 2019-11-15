package apiserver

type MsgType int8

// Well defined Message types
const (
	Error        		MsgType = -1
	Open         		MsgType = 0
	DrawPixel    		MsgType = 1
	LoginUser    		MsgType = 2
	ChangeClientPixel  		MsgType = 3
	Image        		MsgType = 4
	Testing      		MsgType = 5
	DrawResponse 		MsgType = 6
	Close        		MsgType = 9
	VerificationFail    MsgType = 10
	CreateUser          MsgType = 11
)

type Msg struct {
	Type MsgType `json:"type"`
}

// Message generic message recieved through websocket
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

/*
	This message type is intended to be sent from
	the client to the server, signifying that the
	user would like to change a pixel on the canvas.
*/
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

type ImageMsg struct {
	Type         MsgType `json:"type"`
	FormatString string  `json:"formatString"`
}

type TestingMsg struct {
	Type MsgType `json:"type"`
	Msg string  `json:"msg"`
}

type DrawResponseMsg struct {
	Type   MsgType `json:"type"`
	Status int     `json:"status"`
}

type VerificationFailMsg struct {
	Type   MsgType `json:"type"`
	Status int     `json:"status"`
}

type CreateUserMsg struct {
	Type   MsgType `json:"type"`
	Status int     `json:"status"`
	Cooldown int   `json:"cooldown`
}
