package apiserver

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type MsgType uint8

// Well defined Message types
const (
	Open       MsgType = 0
	DrawPixel  MsgType = 1
	CreateUser MsgType = 2
	Close      MsgType = 9
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
	A      int     `json:"a"`
	UserID string  `json:"userID"`
}

func NewDrawPixelMsg(req *http.Request) (*DrawPixelMsg, error) {

	intExtractor := func(queryParam string) (int, error) {
		val, err := strconv.ParseInt(req.URL.Query().Get(queryParam), 10, 8)
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
		Type: DrawPixel,
		X:    x,
		Y:    y,
		R:    r,
		G:    g,
		B:    b,
		A:    a,
	}, nil
}

type CreateUserMsg struct {
	Type MsgType `json:"type"`
	Id   string  `json:"id"`
}
