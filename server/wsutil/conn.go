package wsutil

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	log "github.com/sirupsen/logrus"
)

type Writer interface {
	WriteJSON(v interface{}) error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

type Client struct {
	conn *websocket.Conn
	pool *Pool
	cons consensus.IConsensus
	mux  sync.Mutex
}

func NewClient(c *websocket.Conn, p *Pool) *Client {
	return &Client{
		conn: c,
		pool: p,
	}
}

func (c *Client) Start() {

	defer func() {
		c.pool.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, p, err := c.conn.ReadMessage()
		fmt.Printf("p: " + string(p) + "\n") // we want the ccpmsg we send to be like this
		if err != nil {
			log.Println(err)
			return
		}

		var dat common.Msg
		if err := json.Unmarshal(p, &dat); err != nil {
			log.Printf("error decoding client response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("client response: %q", p)
			panic(err)
		}
		byt := common.MakeTestingMessage("Default Message")

		switch dat.Type {
		case common.DrawPixel:
			fmt.Println("DrawPixel message received.")
			var dpMsg common.DrawPixelMsg
			if err := json.Unmarshal(p, &dpMsg); err == nil {
				log.WithFields(log.Fields{
					"message": dpMsg,
				}).Debug("received ws message")

				fmt.Println("<Start Validate User>")
				lastMove, getErr := c.cons.SyncGetLastUserModification(dpMsg.UserID)

				var userVerification int
				if getErr != nil {
					// Cannot get this user's last modification
					userVerification = 401
				}
				timeSinceLastMove := time.Since(*lastMove)
				if timeSinceLastMove.Milliseconds() >= common.Cooldown.Milliseconds() {
					err := c.cons.SyncSetLastUserModification(dpMsg.UserID, time.Now())
					if err == nil {
						// Successfully updated the user's last modification
						userVerification = 200
					} else {
						// Error from SetLastUserModification call
						userVerification = 403
					}
				} else {
					// User cannot make a move yet.
					userVerification = 429
				}

				if userVerification != 200 {
					// User verification failed
					fmt.Println(fmt.Sprintf("USER %s failed authentication", dpMsg.UserID))
					// send message back to the client indicating verification failure
					byt = common.MakeVerificationFailMessage(userVerification)
					break
				}

				// lastMove, getErr := c.cons.SyncGetLastUserModification()
				fmt.Println("<End Validate User>")

				err := c.cons.SyncUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, common.AlphaMask)
				if err != nil {
					// TODO(backend team): handle error response
				}

				// tell all clients to update their board
				ccpMsg := common.ChangeClientPixelMsg{
					Type:   common.ChangeClientPixel,
					X:      dpMsg.X,
					Y:      dpMsg.Y,
					R:      dpMsg.R,
					G:      dpMsg.G,
					B:      dpMsg.B,
					UserID: dpMsg.UserID,
				}

				msg, _ := json.Marshal(ccpMsg)
				fmt.Printf("msg: " + string(msg))
				c.pool.Broadcast <- ccpMsg
			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("unmarshalling JSON")
			}
			// pretty sure this is not going to be received
		// case ChangeClientPixel:
		// 	fmt.Println("ChangeClientPixel message received.")
		// 	var ccpMsg ChangeClientPixelMsg
		// 	if err := json.Unmarshal(p, &ccpMsg); err == nil {

		// 		fmt.Printf("%+v", ccpMsg)
		// 		// send a message to front end to update this pixel
		// 		if err := c.Conn.WriteMessage(websocket.TextMessage,
		// 			makeChangeClientMessage(ccpMsg.X, ccpMsg.Y, ccpMsg.R, ccpMsg.G, ccpMsg.B, ccpMsg.UserID)); err != nil {
		// 			log.Println(err)
		// 		}
		// 	}

		case common.LoginUser:
			fmt.Println("LoginUser message received.")
			var cu_msg common.LoginUserMsg
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				log.WithFields(log.Fields{
					"message": cu_msg,
				}).Debug("received ws message")
				userID := cu_msg.Email
				//if userID == "" {
				//	http.Error(w, errors.New("empty param: userID").Error(), http.StatusInternalServerError)
				//	return
				//}
				byt = common.MakeUserLoginResponseMsg(400, -1)
				lastMove, getErr := c.cons.SyncGetLastUserModification(userID)
				if getErr == consensus.NoSuchUser {
					setErr := c.cons.SyncSetLastUserModification(userID, time.Unix(0, 0)) // Default Timestamp for New Users
					if setErr == nil {
						byt = common.MakeUserLoginResponseMsg(200, 0)
					}
				} else if lastMove != nil {
					timeSinceLastMove := time.Since(*lastMove)
					if timeSinceLastMove.Milliseconds() >= common.Cooldown.Milliseconds() {
						byt = common.MakeUserLoginResponseMsg(200, 0)
					} else {
						byt = common.MakeUserLoginResponseMsg(200, int(common.Cooldown.Milliseconds()-timeSinceLastMove.Milliseconds()))
					}
				}

			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("unmarshalling JSON")
			}
		default:
			// this is what the case is if the message is recieved from other servers
			fmt.Printf("Message of type: %d received.\n", dat.Type)
		}

		if err := c.WriteMessage(websocket.TextMessage, byt); err != nil {
			log.Println(err)
		}
	}

}

func (c *Client) WriteJSON(v interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Note that in go one cannot simply write:
	//  return c.conn.WriteJSON(v)
	// For more information, refer to:
	err := c.conn.WriteJSON(v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) WriteMessage(messageType int, data []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Note that in go one cannot simply write:
	//  return c.conn.WriteMessage(messageType, data)
	// For more information, refer to:
	err := c.conn.WriteMessage(messageType, data)
	return err
}

func (c *Client) ReadMessage() (messageType int, p []byte, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Note that in go one cannot simply write:
	//  return c.conn.WriteMessage(messageType, data)
	// For more information, refer to:
	messageType, p, err = c.conn.ReadMessage()
	return messageType, p, err
}
