package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"net/http"
	"time"
	"os"

	gwebsocket "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
  "github.com/rgreen312/owlplace/server/websocket"
  )

const (
	AlphaMask = 255
)


var (
	configError   = errors.New("invalid configuration error")
	ImageTemplate = `
    <!DOCTYPE html> 
    <html lang="en">
    <head></head>
    <body>
        <img src="data:image/jpg;base64,{{.Image}}">
    </body>
    `
)

type ApiServer struct {
	pod_ip     string
	node_id     int
	conService *consensus.ConsensusService
}

type Client struct {
	ID   string
	Conn *gwebsocket.Conn
	Pool *Pool
}


const (
	API_PORT int = 3001
)

const (
	CONSENSUS_PORT int = 3010
)

var cooldown = 300000

func NewApiServer(pod_ip string) *ApiServer {

	return &ApiServer{
		conService: nil,
    	pod_ip: pod_ip,
    	node_id: common.IPToNodeId(pod_ip),
	}
}


func (api *ApiServer) StartConsensus(join bool){
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
		if(pod.Status.PodIP != api.pod_ip && !join){
			// Send an http join request to the other nodes
			_, err := http.Get(fmt.Sprintf("http://%s:%d/consensus_join_message", pod.Status.PodIP, API_PORT))
			if(err != nil){
				panic(err)
			}
		}
		servers[common.IPToNodeId(pod.Status.PodIP)] = &common.ServerConfig{
			Hostname: pod.Status.PodIP,
			ApiPort: API_PORT,
			ConsensusPort: CONSENSUS_PORT,
		}
    }
	//At first, just print something so that we know http requests are working inside kubernetes
	fmt.Fprintf(os.Stdout, "Pod Trigger Called\n")


	// Start the consensus service in the background
	conService, err := consensus.NewConsensusService(servers, api.node_id)
	api.conService = conService
	// if err != nil {
	// 	return nil, errors.Wrap(err, "creating ConsensusService")
	// }

	err = conService.Start(join)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "starting ConsensusService")
	// }
}

func (api *ApiServer) ConsensusTrigger(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "ConsensusTrigger\n")
	api.StartConsensus(false)
}

func (api *ApiServer) ConsensusJoinMessage(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(os.Stdout, "ConsensusJoin\n")
	api.StartConsensus(true)
}



func (api *ApiServer) ListenAndServe() {
	pool := NewPool()
	go pool.Start()
	http.HandleFunc("/get_image", api.HTTPGetImage)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/consensus_trigger", api.ConsensusTrigger)
	http.HandleFunc("/consensus_join_message", api.ConsensusJoinMessage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.serveWs(pool, w, r)
	})
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/update_user", api.HTTPUpdateUserList)

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", common.ApiPort), nil)
}

func base64Encode(img *image.RGBA) string {
	// In-memory buffer to store PNG image
	// before we base 64 encode it
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, img)

	// Encode the bytes in the buffer to a base64 string
	return base64.StdEncoding.EncodeToString(buff.Bytes())
}

func (api *ApiServer) HTTPGetImageJson(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"request": req,
	})

	img, err := api.conService.SyncGetImage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(map[string]string{
		"data": base64Encode(img),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (api *ApiServer) HTTPGetImage(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"request": req,
	})

	img, err := api.conService.SyncGetImage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodedString := base64Encode(img)

	tmpl, err := template.New("image").Parse(ImageTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Image": encodedString}
	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *ApiServer) HTTPUpdatePixel(w http.ResponseWriter, req *http.Request) {
	msg, err := NewDrawPixelMsg(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.conService.SyncUpdatePixel(msg.X, msg.Y, msg.R, msg.G, msg.B, msg.A)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * Insert the new user id to the userlist
 */
func (api *ApiServer) HTTPUpdateUserList(w http.ResponseWriter, req *http.Request) {
	// Only for testing
	user_id := req.URL.Query().Get("user_id")
	if user_id == "" {
		http.Error(w, errors.New("empty param: user_id").Error(), http.StatusInternalServerError)
		return
	}

	timestamp := time.Now()
	err := api.conService.SyncSetLastUserModification(user_id, timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	img, err := api.conService.SyncGetImage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodedString := base64Encode(img)
	msg := ImageMsg{
		Type:         Image,
		FormatString: encodedString,
	}

	log.WithFields(log.Fields{
		"ImageMsg": msg,
	}).Debug("constructed websocket message")

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
		_, p, err := c.Conn.ReadMessage()
		fmt.Printf("p: " + string(p) + "\n") // we want the ccpmsg we send to be like this
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
		byt := makeTestingMessage("Default Message")

		switch dat.Type {
		case DrawPixel:
			fmt.Println("DrawPixel message received.")
			var dpMsg DrawPixelMsg
			if err := json.Unmarshal(p, &dpMsg); err == nil {
				log.WithFields(log.Fields{
					"message": dpMsg,
				}).Debug("received ws message")

				// TODO(user team): add user verification here
				err := api.conService.SyncUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, AlphaMask)
				if err != nil {
					// TODO(backend team): handle error response
				}

				// tell all clients to update their board
				ccpMsg := ChangeClientPixelMsg{
					Type:   ChangeClientPixel,
					X:      dpMsg.X,
					Y:      dpMsg.Y,
					R:      dpMsg.R,
					G:      dpMsg.G,
					B:      dpMsg.B,
					UserID: dpMsg.UserID,
				}

				msg, _ := json.Marshal(ccpMsg)
				fmt.Printf("msg: " + string(msg))
				c.Pool.Broadcast <- ccpMsg
			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("unmarshalling JSON")
			}
			// pretty sure this is not going to be received
		// case ChangeClientPixel:
		// 	fmt.Println("ChangeClientPixel message received.")
		// 	var ccpMsg ChangeClientPixelMsg
		// 	if err := json.Unmarshal(p, &ccpMsg); err == nil {

		// 		fmt.Printf("%+v", ccpMsg)
		// 		// send a message to front end to update this pixel
		// 		if err := c.Conn.WriteMessage(gwebsocket.TextMessage,
		// 			makeChangeClientMessage(ccpMsg.X, ccpMsg.Y, ccpMsg.R, ccpMsg.G, ccpMsg.B, ccpMsg.UserID)); err != nil {
		// 			log.Println(err)
		// 		}
		// 	}

		case LoginUser:
			fmt.Println("CreateUser message received.")
			var cu_msg LoginUserMsg
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				log.WithFields(log.Fields{
					"message": cu_msg,
				}).Debug("received ws message")

				timestamp := time.Now()
				err := api.conService.SyncSetLastUserModification(cu_msg.Id, timestamp)
				if err != nil {
					// TODO(backend team): handle error response
				}

			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("unmarshalling JSON")
			}
		default:
			// this is what the case is if the message is recieved from other servers
			fmt.Printf("Message of type: %d received.\n", dat.Type)
		}
		log.Println(byt)
		// // COMMENTED OUT BC CONCURRENT WRITES write message back to the client sent to signal that you received message
		// if err := c.Conn.WriteMessage(gwebsocket.TextMessage, byt); err != nil {
		// 	log.Println(err)
		// }

	}
}
func makeChangeClientMessage(x int, y int, r int, g int, b int, userID string) []byte {
	msg := ChangeClientPixelMsg{
		Type:   ChangeClientPixel,
		X:      x,
		Y:      y,
		R:      r,
		G:      g,
		B:      b,
		UserID: userID,
	}
	bt, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return bt
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

func makeVerificationFailMessage(s int) []byte {
	msg := VerificationFailMsg{
		Type:   VerificationFail,
		Status: s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func makeCreateUserMessage(s int, c int) []byte {
	msg := CreateUserMsg{
		Type:     CreateUser,
		Status:   s,
		Cooldown: c,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}
