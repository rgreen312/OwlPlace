package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	gwebsocket "github.com/gorilla/websocket"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	"github.com/rgreen312/owlplace/server/websocket"
)

type ApiServer struct {
	sendc chan consensus.BackendMessage
	recvc chan consensus.ConsensusMessage
	port  int
}

type Client struct {
	ID   string
	Conn *gwebsocket.Conn
	Pool *Pool
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

func (api *ApiServer) SetupRoutes() {
	pool := NewPool()
	go pool.Start()

	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.HTTPGetImage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.serveWs(pool, w, r)
	})
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/update_user", api.HTTPUpdateUserList)


	// Although there is nothing wrong with this line, it prevents us from running multiple nodes on a single machine.
	// Therefore, I am making failure non-fatal until we have some way of running locally from the same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.port), nil)
}

func (api *ApiServer) HTTPGetImage(w http.ResponseWriter, req *http.Request) {
	// This is the method that will be removed. Displays the image on a webpage

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

func (api *ApiServer) CallGetImage() string {
	// Debug message
	fmt.Fprintf(os.Stdout, "Getting Image From Raft\n")
	// Construct the message
	m := consensus.BackendMessage{Type: consensus.GET_IMAGE}
	// Send a message through the channel
	api.sendc <- m
	var ImageTemplate string = `<!DOCTYPE html>
								<html lang="en"><head></head>
								<body><img src="data:image/jpg;base64,{{.Image}}"></body>`
	imageMsg := <-api.recvc

	if _, err := template.New("image").Parse(ImageTemplate); err != nil {
		fmt.Fprintf(os.Stdout, "Unable to parse image template.\n")
		return "Unable to parse image template"
	} else {

		// Decode the message from the glob
		dec := gob.NewDecoder(&imageMsg.Data)
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

		// TODO  wrap the encodedString in a message

		return encodedString
	}
}

func (api *ApiServer) HTTPUpdatePixel(w http.ResponseWriter, req *http.Request) {
	// This was previously called UpdatePixel. Its logic has been moved into CallUpdatePixel.
	// This is only a wrapper to allow for testing and all design logic should be flowing
	// through the websocket connection.
	x, _ := strconv.Atoi(req.URL.Query().Get("X"))
	y, _ := strconv.Atoi(req.URL.Query().Get("Y"))
	r, _ := strconv.Atoi(req.URL.Query().Get("R"))
	g, _ := strconv.Atoi(req.URL.Query().Get("G"))
	b, _ := strconv.Atoi(req.URL.Query().Get("B"))
	api.CallUpdatePixel(x, y, r, g, b, "RandomHttpUser")
}

// Call this when telling consensus to updatea pixel.
func (api *ApiServer) CallUpdatePixel(x int, y int, r int, g int, b int, userID string) []byte {
	fmt.Println("\nWithin UpdatePixel")

	// TODO verify that the user is able to update a pizel with the User Data Team
	userVerification := api.validateUser(userID)
	if !userVerification {
		// User verification failed

		log.Println(fmt.Sprintf("USER %s failed authentication", userID))
		// TODO return the appropriate failure message
		imageMsg := "FAILURE. TODO make this properly formatted"

		// send message back to the client saying it's been updated or if it failed
		byt := []byte(imageMsg)
		return byt
	}
	// User has now been verified

	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	msg := consensus.UpdatePixelBackendMessage{
		X: strconv.Itoa(x),
		Y: strconv.Itoa(y),
		R: strconv.Itoa(r),
		G: strconv.Itoa(g),
		B: strconv.Itoa(b),
		A: "255",
	}
	log.Printf("UpdatePixelBackendMessage: %+v\n", msg)
	if err := enc.Encode(msg); err != nil {
		log.Fatalf("Error encoding struct: %s", err)
	}

	// Send the encoded message to the backend
	m := consensus.BackendMessage{Type: consensus.UPDATE_PIXEL, Data: encoded_msg}

	// Send BackendMessage
	api.sendc <- m
	consensus_response := <-api.recvc

	// format message back to the client saying it's been updated or if it failed.
	if consensus_response.Type == consensus.SUCCESS {
		fmt.Fprintf(os.Stdout, "Update Success\n")
		byt := makeStatusMessage(200)
		return byt
	} else {
		fmt.Fprintf(os.Stdout, "Update Failure\n")
		byt := makeStatusMessage(400)
		return byt
	}
}


/*
 * Insert the new user id to the userlist
 */
func (api *ApiServer) HTTPUpdateUserList(w http.ResponseWriter, req *http.Request) {
	// Only for testing
	user_id := req.URL.Query().Get("user_id")
	api.CallUpdateUserList(user_id)
}

/*
 * Insert the new user id to the userlist
 */
func (api *ApiServer) CallUpdateUserList(user_id string) []byte {
	byt := []byte("")
	_, ifErr := api.GetLastUserModification(user_id)
	if (ifErr) {
		err := api.SetLastUserModification(user_id, "0")
		if (err) {
			fmt.Println("Cannot update user list")
			byt = []byte("Test failure")
		} else {
			fmt.Println("Successfully created user")
			byt = []byte("Successfully created user")

		}
		
	} else {
		fmt.Println("User already exists")
		byt = []byte("User already exists")
	}
	return byt
}

func (api *ApiServer) validateUser(user_id string) bool {
	now := time.Now().Unix()
	lastMove, ifErr := api.GetLastUserModification(user_id)
	if (ifErr) {
		return false
	}
	lastMoveInt, err := strconv.Atoi(lastMove)
	if err != nil {
		fmt.Println("SOME ERROR")
	}
	if (int(now) - lastMoveInt > 300) {
		api.SetLastUserModification(user_id, strconv.Itoa(int(now)))
		return true
		//image team update pixel		
	} else {
		return false
	}
}

/*
 * This function takes in a user id and returns a string timestamp for the last time that user updated a pixel
 * If there is an error, this function will return an empty string.
 */
func (api *ApiServer) GetLastUserModification(user_id string) (string, bool) {

	// Encode the GetUserDataBackendMessage struct so we can send it in a BackendMessage
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	err := enc.Encode(consensus.GetUserDataBackendMessage{
		UserId: user_id,
	})

	if err != nil {
		return "", true
	}

	// Testing with some dummy data
	m := consensus.BackendMessage{Type: consensus.GET_LAST_USER_UPDATE, Data: encoded_msg}
	api.sendc <- m
	image_msg := <-api.recvc
	if image_msg.Type == consensus.FAILURE {
		return "", true
	} else {
		dec := gob.NewDecoder(&image_msg.Data)
		var timestamp string
		dec.Decode(&timestamp)
		return timestamp, false
	}
}

/*
 * This function takes in a user id and a string timestamp and replaces the user's current last update timestamp with the given timestamp
 * If there is an error, this function will return true. Otherwise the function will return false.
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
		return true
	}
	// Create the BackendMessage with the encoded data
	m := consensus.BackendMessage{Type: consensus.SET_LAST_USER_UPDATE, Data: encoded_msg}
	// Send BackendMessage
	api.sendc <- m
	image_msg := <-api.recvc
	return image_msg.Type != consensus.SUCCESS
}

func (api *ApiServer) Headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func (api *ApiServer) serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn: conn, // this is the same as websocket instance
		Pool: pool,
	}

	// helpful log statement to show connections
	log.Println("Client Connected")
	if err = client.Conn.WriteMessage(1, makeTestingMessage("{\"Hi Client! We just connected :)\"}")); err != nil { // sent upon connection to any clients
		log.Println(err)
	}
	// send image
	serverString := api.CallGetImage()

	// var img_msg bytes.Buffer
	// img := gob.NewEncoder(&img_msg)
	msg := ImageMsg{
		Type:         Image,
		FormatString: serverString,
	}
	log.Printf("ImageMsg: %+v\n", msg)

	var b []byte
	b, err = json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	if err = client.Conn.WriteMessage(1, b); err != nil {
		log.Println(err)
	}

	pool.Register <- client
	client.Read(api)
}

// reading messages now go in here
func (c *Client) Read(api *ApiServer) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		fmt.Printf("Message Received: %+v\n", message)
		
		var dat Msg

		if err := json.Unmarshal(p, &dat); err != nil {
			log.Printf("error decoding client response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("client response: %q", p)
			panic(err)
		}
		// fmt.Println(dat)
		byt := makeTestingMessage("Default Message")

		switch dat.Type {
		case DrawPixel:
			fmt.Println("DrawPixel message received.")
			var dpMsg DrawPixelMsg
			if err := json.Unmarshal(p, &dpMsg); err == nil {
				fmt.Printf("%+v", dpMsg)
				//api.CallUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, dpMsg.UserID)
				byt = api.CallUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, dpMsg.UserID)
			} else {
				fmt.Println("JSON decoding error.")
			}
		case LoginUser:
			fmt.Println("CreateUser message received.")
			var cu_msg LoginUserMsg
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				fmt.Printf("%+v", cu_msg)
				byt = api.CallUpdateUserList(cu_msg.Id)
			} else {
				fmt.Println("JSON decoding error.")
			}
		default:
			// this is what the case is if the message is recieved from other servers
			fmt.Printf("Message of type: %d received.\n", dat.Type)
		}

		// write message back to the client sent to signal that you received message
		if err := c.Conn.WriteMessage(gwebsocket.TextMessage, byt); err != nil {
			log.Println(err)
		}
	}
}

func makeTestingMessage(s string) []byte {
	msg := TestingMsg{
		Type: DrawResponse,
		Msg:  s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func makeStatusMessage(s int) []byte {
	msg := DrawResponseMsg{
		Type:   DrawResponse,
		Status: s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}