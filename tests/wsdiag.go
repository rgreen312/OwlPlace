// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const (
	msgType = 1
	r       = 100
	g       = 100
	b       = 100
	a       = 255
	userid  = "user@example.com"
)

var addr = flag.String("addr", "localhost:3001", "service address")

type DrawPixelMsg struct {
	Type   int    `json:"type"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	R      int    `json:"r"`
	G      int    `json:"g"`
	B      int    `json:"b"`
	A      int    `json:"a"`
	UserID string `json:"userID"`
}

func newMessage(x, y int) *DrawPixelMsg {
	return &DrawPixelMsg{
		Type:   msgType,
		X:      x,
		Y:      y,
		R:      r,
		G:      g,
		B:      b,
		A:      a,
		UserID: userid,
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	i := 0

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if i > 10 {
				interrupt <- os.Interrupt
			}
			msg := newMessage(i, i)
			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("marshal:", err)
				continue
			}
			log.Println("sent:", string(b))
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
