package apiserver

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rgreen312/owlplace/server/common"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan common.ChangeClientPixelMsg
}

func NewPool(c chan common.ChangeClientPixelMsg) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  c,
	}
}

func (pool *Pool) Start(Mux sync.Mutex) {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool after Register: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				Mux.Lock()
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."}) // This line seems to be causing a panic when you refresh it multiple times quickly
				Mux.Unlock()
			}
			fmt.Println("That was all the clients.")
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool after Unregister: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				Mux.Lock()
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
				Mux.Unlock()
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool, which is ", len(pool.Clients))
			for client, _ := range pool.Clients {
				msg, _ := json.Marshal(message)
				fmt.Printf("new_msg: " + string(msg))
				Mux.Lock()
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
				Mux.Unlock()

			}
		}
	}
}
