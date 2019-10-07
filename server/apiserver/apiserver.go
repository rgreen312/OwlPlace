package apiserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"html/template"

	"github.com/gorilla/websocket"
	"github.com/rgreen312/owlplace/server/consensus"
)

type ApiServer struct {
	sendc chan consensus.BackendMessage
	recvc chan consensus.ConsensusMessage
}

func NewApiServer(send_channel chan consensus.BackendMessage, recv_channel chan consensus.ConsensusMessage) *ApiServer {
	return &ApiServer{
		sendc: send_channel,
		recvc: recv_channel,
	}
}

func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.HandleFunc("/ws", api.wsEndpoint)
	log.Fatal(http.ListenAndServe(":3010", nil))
}

func (api *ApiServer) GetImage(w http.ResponseWriter, req *http.Request) {
	// Debug message
	fmt.Fprintf(os.Stdout, "Getting Image From Raft\n")
	// Construct the message
	m := consensus.BackendMessage{Type: consensus.GET_IMAGE}
	// Send a message through the channel
	api.sendc <- m
	var ImageTemplate string = `<!DOCTYPE html>
								<html lang="en"><head></head>
								<body><img src="data:image/jpg;base64,{{.Image}}"></body>`
	image_msg := <-api.recvc

	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		fmt.Fprintf(os.Stdout, "Unable to parse image template.\n")
	} else {
		data := map[string]interface{}{"Image": image_msg.Data}
		if err = tmpl.Execute(w, data); err != nil {
			fmt.Fprintf(os.Stdout, "Unable to execute template.\n")
		}
	}
}

func (api *ApiServer) UpdatePixel(w http.ResponseWriter, req *http.Request) {

	// Decode the request
	update := req.URL.Query().Get("update")
	if update != "" {
		fmt.Fprintf(os.Stdout, update)
		// Testing with some dummy data
		m := consensus.BackendMessage{Type: consensus.UPDATE_PIXEL, Data: update}
		api.sendc <- m
		image_msg := <-api.recvc
		fmt.Fprintf(os.Stdout, image_msg.Data)
	}
}

func (api *ApiServer) Headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// Upgrader "upgrades" HTTP endpoint to WS endpoint, this requires a Read and Write buffer size
var upgrader = websocket.Upgrader{} // use default options

// define a reader which will listen for new messages being sent to our WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		// _ (message type) is an int with value websocket.BinaryMessage or websocket.TextMessage
		// p is []byte
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// holds a map of strings to arbitrary data types
		var dat map[string]interface{}

		if err := json.Unmarshal(p, &dat); err != nil {
			log.Printf("error decoding client response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("client response: %q", p)
			panic(err)
		}
		fmt.Println(dat)

		// convert each attribute to appropriate type
		msgType := dat["type"].(float64) // interface {} is float64, not int
		fmt.Println(msgType)

		switch msgType {
		case 1:
			fmt.Println("one")

			// TODO: call the updating function

			// send message back to the client saying it's been updated
			byt := []byte(`Pixel 1 has been updated!`)
			if err := conn.WriteMessage(websocket.TextMessage, byt); err != nil {
				log.Println(err)
				return
			}
		case 2:
			fmt.Println("two")
		case 3:
			fmt.Println("three")
		}
	}
}

func (api *ApiServer) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// checks if incoming request is allowed to connect, otherwise CORS error, currently always true
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// helpful log statement to show connections
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client! We just connected :)")) // sent upon connection to any clients

	reader(ws)

}
