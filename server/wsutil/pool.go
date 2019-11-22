package wsutil

import (
	log "github.com/sirupsen/logrus"

	"github.com/rgreen312/owlplace/server/common"
)

const (
	defaultPoolSize = 100
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan common.ChangeClientPixelMsg
}

func NewPool() *Pool {
	// We use buffered channels for performance reasons.  An unbuffered
	// channel is effectively coordinated function invocations, while a
	// buffered channel allows us to queue several without blocking the calling
	// thread.
	return &Pool{
		Register:   make(chan *Client, defaultPoolSize),
		Unregister: make(chan *Client, defaultPoolSize),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan common.ChangeClientPixelMsg, defaultPoolSize),
	}
}

func (p *Pool) forEach(consumer func(*Client) error) error {
	for client, _ := range p.Clients {
		if err := consumer(client); err != nil {
			return err
		}
	}
	return nil
}

// WriteJson writes a JSON message to all connected websockets.
func (p *Pool) WriteJSON(v interface{}) error {
	return nil
}

// WriteMessage writes a byte array to all connected websockets.
func (p *Pool) WriteMessage(messageType int, data []byte) error {
	return nil
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			// Add the client to our map of clients
			p.Clients[client] = true
			break
		case client := <-p.Unregister:
			// Remove the client who requested to leave
			delete(p.Clients, client)
			break
		case message := <-p.Broadcast:

			log.WithFields(log.Fields{
				"message":    message,
				"numClients": len(p.Clients),
			}).Debug("broadcasting message from websocket pool")

			p.forEach(func(c *Client) error {
				// TODO: responsibly handle errors here.  realistically we need
				// to determine what the issue is, and if the issue is a stale
				// connection we should remove this client from our pool.
				if err := c.WriteJSON(message); err != nil {
					return err
				}

				return nil
			})
		}
	}
}
