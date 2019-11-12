package websocket

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// type MsgType uint8

// // Well defined Message types
// const (
// 	Open      MsgType = 0
// 	DrawPixel MsgType = 1
// 	LoginUser MsgType = 2
// 	Image     MsgType = 4
// 	Testing   MsgType = 5
// 	Close     MsgType = 9
// )

// type Client struct {
// 	ID   string
// 	Conn *websocket.Conn
// 	Pool *Pool
// }

// type Msg struct {
// 	Type MsgType `json:"type"`
// }

// type DrawPixelMsg struct {
// 	Type   MsgType `json:"type"`
// 	X      int     `json:"x"`
// 	Y      int     `json:"y"`
// 	R      int     `json:"r"`
// 	G      int     `json:"g"`
// 	B      int     `json:"b"`
// 	UserID string  `json:"userID"`
// }

// type TestingMsg struct {
// 	Type MsgType `json:"type"`
// 	Msg  string  `json:"msg"`
// }

// type LoginUserMsg struct {
// 	Type MsgType `json:"type"`
// 	Id   string  `json:"id"`
// }

// func makeTestingMessage(s string) []byte {
// 	msg := TestingMsg{
// 		Type: Testing,
// 		Msg:  s,
// 	}

// 	// var b []byte
// 	b, err := json.Marshal(msg)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return b
// }

// // reading messages now go in here
// func (c *Client) Read() {
// 	defer func() {
// 		c.Pool.Unregister <- c
// 		c.Conn.Close()
// 	}()

// 	for {
// 		messageType, p, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		message := Message{Type: messageType, Body: string(p)}

// 		// our code vvvv
// 		var dat Msg

// 		if err := json.Unmarshal(p, &dat); err != nil {
// 			log.Printf("error decoding client response: %v", err)
// 			if e, ok := err.(*json.SyntaxError); ok {
// 				log.Printf("syntax error at byte offset %d", e.Offset)
// 			}
// 			log.Printf("client response: %q", p)
// 			panic(err)
// 		}
// 		// fmt.Println(dat)

// 		byt := makeTestingMessage("Default Message")

// 		switch dat.Type {
// 		case DrawPixel:
// 			fmt.Println("DrawPixel message received.")
// 			var dpMsg DrawPixelMsg
// 			if err := json.Unmarshal(p, &dpMsg); err == nil {
// 				fmt.Printf("%+v", dpMsg)
// 				byt = api.CallUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, dpMsg.UserID)
// 			} else {
// 				fmt.Println("JSON decoding error.")
// 			}
// 		case LoginUser:
// 			fmt.Println("CreateUser message received.")
// 			var cuMsg LoginUserMsg
// 			if err := json.Unmarshal(p, &cuMsg); err == nil {
// 				fmt.Printf("%+v", cuMsg)
// 				email := cuMsg.Id
// 				// byt = api.CallUpdateUserList()
// 				byt = []byte("{\"type\": 2, \"Id\": \"" + email + "\"}")
// 			} else {
// 				fmt.Println("JSON decoding error.")
// 			}
// 		default:
// 			fmt.Printf("Message of type: %d received.", dat.Type)
// 		}

// 		if err := c.Conn.WriteMessage(websocket.TextMessage, byt); err != nil {
// 			log.Println(err)
// 		}

// 		// our code ^^^^^
// 		c.Pool.Broadcast <- message
// 		fmt.Printf("Message Received: %+v\n", message)
// 	}
// }

// // OURS define a reader which will listen for new messages being sent to our WebSocket endpoint
// //func Reader(conn *websocket.Conn) {
// //	for {
// //		// read in a message
// //		// _ (message type) is an int with value websocket.BinaryMessage or websocket.TextMessage
// //		// p is []byte
// //		_, p, err := conn.ReadMessage()
// //		if err != nil {
// //			log.Println(err)
// //			return
// //		}
// //
// //		var dat Msg
// //
// //		if err := json.Unmarshal(p, &dat); err != nil {
// //			log.Printf("error decoding client response: %v", err)
// //			if e, ok := err.(*json.SyntaxError); ok {
// //				log.Printf("syntax error at byte offset %d", e.Offset)
// //			}
// //			log.Printf("client response: %q", p)
// //			panic(err)
// //		}
// //		// fmt.Println(dat)
// //
// //		byt := makeTestingMessage("Default Message")
// //
// //		switch dat.Type {
// //		case DrawPixel:
// //			fmt.Println("DrawPixel message received.")
// //			var dp_msg DrawPixelMsg
// //			if err := json.Unmarshal(p, &dp_msg); err == nil {
// //				fmt.Printf("%+v", dp_msg)
// //				byt = api.CallUpdatePixel(dp_msg.X, dp_msg.Y, dp_msg.R, dp_msg.G, dp_msg.B, dp_msg.UserID)
// //			} else {
// //				fmt.Println("JSON decoding error.")
// //			}
// //		case LoginUser:
// //			fmt.Println("CreateUser message received.")
// //			var cu_msg LoginUserMsg
// //			if err := json.Unmarshal(p, &cu_msg); err == nil {
// //				fmt.Printf("%+v", cu_msg)
// //				email := cu_msg.Id
// //				// byt = api.CallUpdateUserList()
// //				byt = []byte("{\"type\": 2, \"Id\": \"" + email + "\"}")
// //			} else {
// //				fmt.Println("JSON decoding error.")
// //			}
// //		default:
// //			fmt.Printf("Message of type: %d received.", dat.Type)
// //		}
// //
// //		if err := conn.WriteMessage(websocket.TextMessage, byt); err != nil {
// //			log.Println(err)
// //		}
// //	}
// //}
