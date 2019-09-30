package apiserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "ERICA\n")
}

func Headers(w http.ResponseWriter, req *http.Request) {

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
