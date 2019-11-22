package wsutil

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Read / Write Buffer
	wsUpgraderReadBufferSize  = 1024
	wsUpgraderWriteBufferSize = 1024
)
const ()

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  wsUpgraderReadBufferSize,
		WriteBufferSize: wsUpgraderWriteBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// Client is a middleman between the websocket connection and the pool.
type Client struct {
	pool *Pool

	// The websocket connection.
	conn *websocket.Conn

	// Consensus Provider.
	cons consensus.IConsensus

	// Buffered channel of outbound messages.
	Send chan []byte
}

// TODO: Remove this method.  This is a temporary fix so we don't chance
// everything at once.  See the comment on handleMessage for more information.
func (c *Client) handleDrawPixel(p []byte) {
	var dpMsg common.DrawPixelMsg
	err := json.Unmarshal(p, &dpMsg)
	if err != nil {
		log.WithFields(log.Fields{
			"ws packet": p,
			"err":       err,
		}).Debug("unmarshalling JSON")
		return
	}

	lastMove, err := c.cons.SyncGetLastUserModification(dpMsg.UserID)
	if err != nil && err != consensus.NoSuchUser {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("retrieving user's last modification")
		return
	}

	timeSinceLastMove := time.Since(*lastMove)
	if timeSinceLastMove.Milliseconds() < common.Cooldown.Milliseconds() {

		// Here we'd like to send a message to the client indicating that they
		// need to wait a bit longer before making another change to the
		// canvas.
		message := common.MakeVerificationFailMessage(0)
		c.Send <- message
		return
	}

	err = c.cons.SyncUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, common.AlphaMask)
	if err != nil {
		// TODO: handle error response
		return
	}

	// This should not be done until after the user is able to successfully
	// update the canvas.
	err = c.cons.SyncSetLastUserModification(dpMsg.UserID, time.Now())
	if err != nil {
		// TODO: handle error response, potentially shut down the websocket
		// connection? the issue here is that we realistically need to get this
		// set after allowing the user to update the pixel.  however, it's that
		// big of a deal, as we don't expect very many errors to appear here.
		return
	}

	// Tell all clients to update their board
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
	c.pool.Broadcast <- msg
}

func (c *Client) handleLoginUser(p []byte) {

	var cu_msg common.LoginUserMsg
	err := json.Unmarshal(p, &cu_msg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("unmarshalling JSON")
	}

	log.WithFields(log.Fields{
		"message": cu_msg,
	}).Debug("received ws message")
	userID := cu_msg.Email

	lastMove, err := c.cons.SyncGetLastUserModification(userID)
	if err != nil && err != consensus.NoSuchUser {
		// TODO: determine how we want to handle this error
		byt := common.MakeUserLoginResponseMsg(501, 0)
		c.Send <- byt
		return
	}

	if err == consensus.NoSuchUser {
		// Set the default timestamp for a new user to the current time less
		// the cooldown period.  This allows immediate editing.
		stamp := time.Now().Add(-common.Cooldown)
		err := c.cons.SyncSetLastUserModification(userID, stamp) // Default Timestamp for New Users
		if err != nil {
			byt := common.MakeUserLoginResponseMsg(501, 0)
			c.Send <- byt
			return
		}
		lastMove = &stamp
	}

	timeSinceLastMove := time.Since(*lastMove)
	userCooldown := int(common.Cooldown.Milliseconds() - timeSinceLastMove.Milliseconds())
	userCooldown = int(math.Max(0, float64(userCooldown)))
	byt := common.MakeUserLoginResponseMsg(200, userCooldown)
	c.Send <- byt
}

// TODO: Refactor this method.  We should not be using a switch here and
// stuffing it with untestable code.  Ideally, we'd refactor all messages to
// support visitor / double visitor pattern.
func (c *Client) handleMessage(p []byte) {
	var dat common.Msg
	if err := json.Unmarshal(p, &dat); err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"client": c,
		}).Debug("Unmarshalling JSON Message from Client")
		return
	}

	switch dat.Type {
	case common.DrawPixel:
		c.handleDrawPixel(p)
	case common.LoginUser:
		c.handleLoginUser(p)
	default:
		log.WithFields(log.Fields{
			"message type": dat.Type,
		}).Debug("client received unknown message")
	}
}

// readPump pumps messages from the websocket connection to the pool.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.pool.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			// If the error we receive is NOT in the list of expected codes
			// (CloseGoingAway/CloseAbnormalClosure), report the error.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithFields(log.Fields{
					"err": err,
				}).Debug("reading ws message")
			}
			break
		}
		c.handleMessage(p)
	}
}

// writePump pumps messages from the pool to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The pool closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			//w, err := c.conn.NextWriter(websocket.TextMessage)
			//if err != nil {
			//return
			//}
			//w.Write(message)

			//// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//w.Write(newline)
			//w.Write(<-c.send)
			//}

			//if err := w.Close(); err != nil {
			//return
			//}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("upgrading conn to ws")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &Client{
		pool: pool,
		conn: conn,
		Send: make(chan []byte, 256),
	}
	pool.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
