package apiserver

import (
	"net/http"
	"strconv"
	"github.com/pkg/errors"
	"github.com/rgreen312/owlplace/server/common"
)


type Msg struct {
	Type common.MsgType `json:"type"`
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
	Type   common.MsgType `json:"type"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	R      int     `json:"r"`
	G      int     `json:"g"`
	B      int     `json:"b"`
	A      int     `json:"a"`
	UserID string  `json:"userID"`
}

func NewDrawPixelMsg(req *http.Request) (*DrawPixelMsg, error) {

	intExtractor := func(queryParam string) (int, error) {
		val, err := strconv.ParseUint(req.URL.Query().Get(queryParam), 10, 8)
		return int(val), err
	}

	x, err := intExtractor("X")
	if err != nil {
		return nil, errors.Wrap(err, "extracting x")
	}
	y, err := intExtractor("Y")
	if err != nil {
		return nil, errors.Wrap(err, "extracting y")
	}
	r, err := intExtractor("R")
	if err != nil {
		return nil, errors.Wrap(err, "extracting r")
	}
	g, err := intExtractor("G")
	if err != nil {
		return nil, errors.Wrap(err, "extracting g")
	}
	b, err := intExtractor("B")
	if err != nil {
		return nil, errors.Wrap(err, "extracting b")
	}
	// TODO(backend team): determine if you'd like to have the alpha mask required
	//a, err := intExtractor("A")
	//if err != nil {
	//return nil, errors.Wrap(err, "extracting a")
	//}
	a := 255

	// TODO(backend team): add user id parsing
	return &DrawPixelMsg{
		Type: common.DrawPixel,
		X:    x,
		Y:    y,
		R:    r,
		G:    g,
		B:    b,
		A:    a,
	}, nil
}

type LoginUserMsg struct {
	Type common.MsgType `json:"type"`
	Email   string  `json:"email"`
}

type ImageMsg struct {
	Type         common.MsgType `json:"type"`
	FormatString string  `json:"formatString"`
}

type TestingMsg struct {
	Type common.MsgType `json:"type"`
	Msg  string  `json:"msg"`
}

type DrawResponseMsg struct {
	Type   common.MsgType `json:"type"`
	Status int     `json:"status"`
}

type VerificationFailMsg struct {
	Type   common.MsgType `json:"type"`
	Status int     `json:"status"`
}


type UserLoginResponseMsg struct {
	Type     common.MsgType `json:"type"`
	Status   int     `json:"status"`
	Cooldown int     `json:"cooldown"`
}
