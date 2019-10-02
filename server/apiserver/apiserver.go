package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rgreen312/owlplace/server/consensus"
)

type ApiServer struct {
	c chan consensus.BackendMessage
}

func NewApiServer(channel chan consensus.BackendMessage) *ApiServer {
	return &ApiServer{
		c: channel,
	}
}

func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/hello", api.PutHello)
	http.HandleFunc("/headers", api.Headers)
	http.ListenAndServe(":3000", nil)
}

func (api *ApiServer) PutHello(w http.ResponseWriter, req *http.Request) {
	// Debug message
	fmt.Fprintf(os.Stdout, "Sending Hello Message\n")
	// Construct the message
	m := consensus.BackendMessage{Dummy: "put server hello2!"}
	// Send a message through the channel
	api.c <- m
}

func (api *ApiServer) Headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

var upgrader = websocket.Upgrader{} // use default options

func Echo(w http.ResponseWriter, r *http.Request) {

	// fmt.Fprintf(w, "ECHO\n")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("server read:", err)
			break
		}
		log.Printf("server recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("server write:", err)
			break
		}
	}
}
