package apiserver

import (
	"fmt"
	"sync"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan ChangeClientPixelMsg
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan ChangeClientPixelMsg),
	}
}

func (pool *Pool) Start(Mux sync.Mutex) {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				Mux.Lock()
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
				Mux.Unlock()
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				Mux.Lock()
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
				Mux.Unlock()
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}

			}
		}
	}
}
