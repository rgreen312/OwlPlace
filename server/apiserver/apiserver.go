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
	"strings"
	"strconv"
	"github.com/gorilla/websocket"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ApiServer struct {
	sendc 				chan consensus.BackendMessage
	recvc 				chan consensus.ConsensusMessage
	port  				int
	pod_ip 				string
}

func NewApiServer(pod_ip string) *ApiServer {

	apiPort := 3001

	// Make the channels that for api server and consensus module communication
	sendc := make(chan consensus.BackendMessage)
	recvc := make(chan consensus.ConsensusMessage)

	return &ApiServer{
		sendc: sendc,
		recvc: recvc,
		port:  apiPort,
		pod_ip: pod_ip,
	}
}

func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.HandleFunc("/consensus_trigger", api.ConsensusTrigger)
	// http.HandleFunc("/consensus_join_message", api.ConsensusJoinMessage)
	http.HandleFunc("/ws", api.wsEndpoint)
	http.ListenAndServe(fmt.Sprintf(":%d", api.port), nil)
}

func (api *ApiServer) IPToNodeId(ip_address string) int {
	// This function maps ip addresses to node-ids
	components := strings.Split(ip_address, ".")
	combined := fmt.Sprintf("%03s%03s", components[2], components[3])
	node_id, _ := strconv.Atoi(combined)
	return node_id
}

func (api *ApiServer) ConsensusTrigger(w http.ResponseWriter, req *http.Request) {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("dev").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	servers :=  make(map[int]*common.ServerConfig)
	for _, pod := range pods.Items {
		if(pod.Status.PodIP != api.pod_ip){
			servers[api.IPToNodeId(pod.Status.PodIP)] = &common.ServerConfig{
				Hostname: pod.Status.PodIP,
				ApiPort: 3001,
				ConsensusPort: 3010,
			}
		}
    }
	//At first, just print something so that we know http requests are working inside kubernetes
	fmt.Fprintf(os.Stdout, "Pod Trigger Called\n")

	// Start the consensus service in the background
	go consensus.MainConsensus(api.sendc, api.recvc, servers, api.IPToNodeId(api.pod_ip) )
}

// func (api *ApiServer) ConsensusJoinMessage(){
// 		// Start the consensus service in the background
// 	go consensus.MainConsensus(sendc, recvc, servers, nodeId)
// }

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
	imageMsg := <-api.recvc

	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		fmt.Fprintf(os.Stdout, "Unable to parse image template.\n")
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
	if(image_msg.Type == consensus.FAILURE){
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

// Upgrader "upgrades" HTTP endpoint to WS endpoint, this requires a Read and Write buffer size
var upgrader = websocket.Upgrader{} // use default options

// define a reader which will listen for new messages being sent to our WebSocket endpoint
func (api *ApiServer) reader(conn *websocket.Conn) {
	for {
		// read in a message
		// _ (message type) is an int with value websocket.BinaryMessage or websocket.TextMessage
		// p is []byte
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var dat Msg

		if err := json.Unmarshal(p, &dat); err != nil {
			log.Printf("error decoding client response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("client response: %q", p)
			panic(err)
		}
		fmt.Println(dat)

		byt := []byte("Default message")

		switch dat.Type {
		case DrawPixel:
			fmt.Println("DrawPixel message received.")
			var dp_msg DrawPixelMsg
			if err := json.Unmarshal(p, &dp_msg); err == nil {
				fmt.Printf("%+v", dp_msg)

				byt = api.CallUpdatePixel(dp_msg.X, dp_msg.Y, dp_msg.R, dp_msg.G, dp_msg.B, dp_msg.UserID)
			} else {
				fmt.Println("JSON decoding error.")
			}
		case CreateUser:
			fmt.Println("CreateUser message received.")
			var cu_msg CreateUserMsg
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				fmt.Printf("%+v", cu_msg)
			} else {
				fmt.Println("JSON decoding error.")
			}
		default:
			fmt.Printf("Message of type: %d received.", dat.Type)
		}

		if err := conn.WriteMessage(websocket.TextMessage, byt); err != nil {
			log.Println(err)
		}
	}
}

// Call this when telling consensus to updatea pixel.
func (api *ApiServer) CallUpdatePixel(x int, y int, r int, g int, b int, userID string) []byte {
	fmt.Println("\nWithin CallUpdatePixel")

	// TODO verify that the user is able to update a pizel with the User Data Team
	userVerification := true
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
		X: string(x),
		Y: string(y),
		R: string(r),
		G: string(g),
		B: string(b),
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
	if consensus_response.Type == consensus.SUCCESS {
		fmt.Fprintf(os.Stdout, "Update Success")
	}
	// format message back to the client saying it's been updated or if it failed.
	// byt := []byte(consensus_response.Data)
	byt := []byte("test success")
	return byt
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

	api.reader(ws)

}
