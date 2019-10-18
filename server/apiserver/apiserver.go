package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
)

type ApiServer struct {
	sendc chan consensus.BackendMessage
	recvc chan consensus.ConsensusMessage
	port  int
}

func NewApiServer(servers map[int]*common.ServerConfig, nodeId int) *ApiServer {

	apiPort := servers[nodeId].ApiPort

	// Make the channels that for api server and consensus module communication
	sendc := make(chan consensus.BackendMessage)
	recvc := make(chan consensus.ConsensusMessage)

	// Start the consensus service in the background
	go consensus.MainConsensus(sendc, recvc, servers, nodeId)

	return &ApiServer{
		sendc: sendc,
		recvc: recvc,
		port:  apiPort,
	}
}

func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.HandleFunc("/ws", api.wsEndpoint)

	// Although there is nothing wrong with this line, it prevents us from running multiple nodes on a single machine.
	// Therefore, I am making failure non-fatal until we have some way of running locally from the same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.port), nil)
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

		// Decode the message from the glob
		dec := gob.NewDecoder(&image_msg.Data)
		var img image.RGBA
		dec.Decode(&img)

		// In-memory buffer to store PNG image
		// before we base 64 encode it
		var buff bytes.Buffer

		// The Buffer satisfies the Writer interface so we can use it with Encode
		// In previous example we encoded to a file, this time to a temp buffer
		png.Encode(&buff, &img)

		// Encode the bytes in the buffer to a base64 string
		encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

		data := map[string]interface{}{"Image": encodedString}
		if err = tmpl.Execute(w, data); err != nil {
			fmt.Fprintf(os.Stdout, "Unable to execute template.\n")
		}
	}
}

func (api *ApiServer) UpdatePixel(w http.ResponseWriter, req *http.Request) {

	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	msg := consensus.UpdatePixelBackendMessage{
		X: req.URL.Query().Get("X"),
		Y: req.URL.Query().Get("Y"),
		R: req.URL.Query().Get("R"),
		G: req.URL.Query().Get("G"),
		B: req.URL.Query().Get("B"),
		A: "255",
	}
	log.Printf("UpdatePixelBackendMessage: %+v\n", msg)
	if err := enc.Encode(msg); err != nil {
		log.Fatalf("Error encoding struct: %s", err)
	}

	// Send the encoded message to the backend
	m := consensus.BackendMessage{Type: consensus.UPDATE_PIXEL, Data: encoded_msg}
	api.sendc <- m
	consensus_response := <-api.recvc
	if consensus_response.Type == consensus.SUCCESS {
		fmt.Fprintf(os.Stdout, "Update Success")
	}

}

/*
 * This function takes in a user id and returns a string timestamp for the last time that user updated a pixel
 * If there is an error, this function will return an empty string.
 */
func (api *ApiServer) GetLastUserModification(user_id string) string {

	// Encode the GetUserDataBackendMessage struct so we can send it in a BackendMessage
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	err := enc.Encode(consensus.GetUserDataBackendMessage{
		UserId: user_id,
	})

	if err != nil {
		return ""
	}

	// Testing with some dummy data
	m := consensus.BackendMessage{Type: consensus.GET_LAST_USER_UPDATE, Data: encoded_msg}
	api.sendc <- m
	image_msg := <-api.recvc

	dec := gob.NewDecoder(&image_msg.Data)
	var timestamp string
	dec.Decode(&timestamp)

	return timestamp
}

/*
 * This function takes in a user id and a string timestamp and replaces the user's current last update timestamp with the given timestamp
 * If there is an error, this function will return false. Otherwise the function will return true.
 */
func (api *ApiServer) SetLastUserModification(user_id string, last_modification string) bool {

	// Encode the SetUserDataBackendMessage struct so we can send it in a BackendMessage
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	err := enc.Encode(consensus.SetUserDataBackendMessage{
		UserId:    user_id,
		Timestamp: last_modification,
	})

	if err != nil {
		return false
	}
	// Create the BackendMessage with the encoded data
	m := consensus.BackendMessage{Type: consensus.SET_LAST_USER_UPDATE, Data: encoded_msg}

	// Send BackendMessage
	api.sendc <- m
	image_msg := <-api.recvc
	return image_msg.Type == consensus.SUCCESS
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
