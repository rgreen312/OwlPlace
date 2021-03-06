package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

type MsgType int8

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
	A      int     `json:"a"`
	UserID string  `json:"userID"`
}

// base64Encode returns a base64 string representation of an RGBA image.
func Base64Encode(img *image.RGBA) string {
	// In-memory buffer to store PNG image
	// before we base 64 encode it
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, img)

	// Encode the bytes in the buffer to a base64 string
	return base64.StdEncoding.EncodeToString(buff.Bytes())
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

	// TODO(backend team): determine if you'd like to have the alpha mask
	// required
	a, err := intExtractor("A")
	if err != nil {
		a = 255
	}

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

type LoginUserMsg struct {
	Type  MsgType `json:"type"`
	Email string  `json:"email"`
}

type ImageMsg struct {
	Type         MsgType `json:"type"`
	FormatString string  `json:"formatString"`
}

type TestingMsg struct {
	Type MsgType `json:"type"`
	Msg  string  `json:"msg"`
}

type DrawResponseMsg struct {
	Type          MsgType `json:"type"`
	Status        int     `json:"status"`
	RemainingTime int     `json:"remainingTime"`
}

type VerificationFailMsg struct {
	Type   MsgType `json:"type"`
	Status int     `json:"status"`
}

type UserLoginResponseMsg struct {
	Type     MsgType `json:"type"`
	Status   int     `json:"status"`
	Cooldown int     `json:"cooldown"`
}

func MakeTestingMessage(s string) []byte {
	msg := TestingMsg{
		Type: Testing,
		Msg:  s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func MakeStatusMessage(s int, remaining int) []byte {
	msg := DrawResponseMsg{
		Type:          DrawResponse,
		Status:        s,
		RemainingTime: remaining,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func MakeCooldownMessage(time int) []byte {
	msg := DrawResponseMsg{
		Type:          DrawResponse,
		Status:        429,
		RemainingTime: time,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func MakeVerificationFailMessage(s int) []byte {
	msg := VerificationFailMsg{
		Type:   VerificationFail,
		Status: s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func MakeUserLoginResponseMsg(s int, c int) []byte {
	msg := UserLoginResponseMsg{
		Type:     UserLoginResponse,
		Status:   s,
		Cooldown: c,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func MakeChangeClientMessage(x int, y int, r int, g int, b int, userID string) []byte {
	msg := ChangeClientPixelMsg{
		Type:   ChangeClientPixel,
		X:      x,
		Y:      y,
		R:      r,
		G:      g,
		B:      b,
		UserID: userID,
	}
	bt, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return bt
}
