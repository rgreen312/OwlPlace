package wsutil

import (
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan []byte
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		// TODO: determine if it would be better to have a buffered channel
		// here.
		Broadcast: make(chan []byte),
	}
}

func (p *Pool) forEach(consumer func(*Client)) {
	for client, _ := range p.Clients {
		consumer(client)
	}
}

func (p *Pool) Run() {
	for {
		select {
		case client := <-p.Register:
			// Add the client to our map of clients
			p.Clients[client] = true
		case client := <-p.Unregister:
			// Remove the client who requested to leave
			delete(p.Clients, client)
		case message := <-p.Broadcast:
			log.WithFields(log.Fields{
				"message":    message,
				"numClients": len(p.Clients),
			}).Debug("broadcasting message from ws pool")

			p.forEach(func(c *Client) {
				select {
				case c.Send <- message:
				default:
					log.WithFields(log.Fields{
						"message": message,
						"client":  c,
					}).Debug("failed to send message to client's channel")
					close(c.Send)
					delete(p.Clients, c)
				}
			})
		}
	}
}
