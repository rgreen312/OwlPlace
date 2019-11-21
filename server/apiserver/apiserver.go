package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"time"
	"os"
	"sync"


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
	pool       *Pool
	Mux        sync.Mutex
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


var cooldown, _ = time.ParseDuration("3s")

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


	// Now that consensus is active, we can start listening for websocket connections
	pool := NewPool(api.conService.Broadcast)
	api.pool = pool
	go pool.Start(api.Mux)


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

  	http.HandleFunc("/consensus_trigger", api.ConsensusTrigger)
	http.HandleFunc("/consensus_join_message", api.ConsensusJoinMessage)

	http.HandleFunc("/ws", api.serveWs)
	// http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	// http.HandleFunc("/update_user", api.HTTPUserLogin)

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

//func (api *ApiServer) HTTPGetImageJson(w http.ResponseWriter, req *http.Request) {
//	log.WithFields(log.Fields{
//		"request": req,
//	})
//
//	img, err := api.conService.SyncGetImage()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	js, err := json.Marshal(map[string]string{
//		"data": base64Encode(img),
//	})
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(js)
//}
//
//func (api *ApiServer) HTTPGetImage(w http.ResponseWriter, req *http.Request) {
//	log.WithFields(log.Fields{
//		"request": req,
//	})
//
//	img, err := api.conService.SyncGetImage()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	encodedString := base64Encode(img)
//
//	tmpl, err := template.New("image").Parse(ImageTemplate)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	data := map[string]interface{}{"Image": encodedString}
//	if err = tmpl.Execute(w, data); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func (api *ApiServer) HTTPUpdatePixel(w http.ResponseWriter, req *http.Request) {
//	msg, err := NewDrawPixelMsg(req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	err = api.conService.SyncUpdatePixel(msg.X, msg.Y, msg.R, msg.G, msg.B, msg.A)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func (api *ApiServer) HTTPUserLogin(w http.ResponseWriter, req *http.Request) []byte {
//	userID := req.URL.Query().Get("id")
//	//if userID == "" {
//	//	http.Error(w, errors.New("empty param: userID").Error(), http.StatusInternalServerError)
//	//	return
//	//}
//	byt := makeUserLoginResponseMsg(400, -1)
//	lastMove, getErr := api.conService.SyncGetLastUserModification(userID)
//	if getErr == consensus.NoSuchUser {
//		setErr := api.conService.SyncSetLastUserModification(userID, time.Unix(0,0)) // Default Timestamp for New Users
//		if setErr != nil {
//			byt = makeUserLoginResponseMsg(200, 0)
//		}
//	} else if lastMove != nil {
//		timeSinceLastMove := time.Since(*lastMove)
//		if timeSinceLastMove.Milliseconds() >= cooldown.Milliseconds() {
//			byt = makeUserLoginResponseMsg(200, 0)
//		} else {
//			byt = makeUserLoginResponseMsg(200, int(cooldown.Milliseconds() - timeSinceLastMove.Milliseconds()))
//		}
//	}
//	return byt
//	//if err != nil {
//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	//	return
//	//}
//}

func (api *ApiServer) serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn: conn, // this is the same as websocket instance
		Pool: api.pool,
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
		Type:         common.Image,
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

	api.pool.Register <- client
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
		case common.DrawPixel:
			fmt.Println("DrawPixel message received.")
			var dpMsg DrawPixelMsg
			if err := json.Unmarshal(p, &dpMsg); err == nil {
				log.WithFields(log.Fields{
					"message": dpMsg,
				}).Debug("received ws message")

				fmt.Println("<Start Validate User>")
				lastMove, getErr := api.conService.SyncGetLastUserModification(dpMsg.UserID)
				var userVerification int
				if getErr != nil {
					// Cannot get this user's last modification
					userVerification = 401
				}
				timeSinceLastMove := time.Since(*lastMove)

				if timeSinceLastMove.Milliseconds() >= cooldown.Milliseconds() {
					err := api.conService.SyncSetLastUserModification(dpMsg.UserID, time.Now())
					if err == nil {
						// Successfully updated the user's last modification
						userVerification = 200	
					} else {
						// Error from SetLastUserModification call
						userVerification = 403
					}
				} else {
					// User cannot make a move yet.
					userVerification = 429
				}

				if userVerification != 200 {
					// User verification failed
					fmt.Println(fmt.Sprintf("USER %s failed authentication", dpMsg.UserID))
					// send message back to the client indicating verification failure
					byt = makeVerificationFailMessage(userVerification)
					break
				}

				// lastMove, getErr := api.conService.SyncGetLastUserModification()
				fmt.Println("<End Validate User>")

				err := api.conService.SyncUpdatePixel(dpMsg.X, dpMsg.Y, dpMsg.R, dpMsg.G, dpMsg.B, AlphaMask)
				if err != nil {
					// TODO(backend team): handle error response
				}

			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("unmarshalling JSON")
			}
      
		case common.LoginUser:
			fmt.Println("CreateUser message received.")

			var cu_msg LoginUserMsg
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				log.WithFields(log.Fields{
					"message": cu_msg,
				}).Debug("received ws message")
				userID := cu_msg.Email
				//if userID == "" {
				//	http.Error(w, errors.New("empty param: userID").Error(), http.StatusInternalServerError)
				//	return
				//}
				byt = makeUserLoginResponseMsg(400, -1)
				lastMove, getErr := api.conService.SyncGetLastUserModification(userID)
				if getErr == consensus.NoSuchUser {
					setErr := api.conService.SyncSetLastUserModification(userID, time.Unix(0,0)) // Default Timestamp for New Users
					if setErr == nil {
						byt = makeUserLoginResponseMsg(200, 0)
					}
				} else if lastMove != nil {
					timeSinceLastMove := time.Since(*lastMove)
					if timeSinceLastMove.Milliseconds() >= cooldown.Milliseconds() {
						byt = makeUserLoginResponseMsg(200, 0)
					} else {
						byt = makeUserLoginResponseMsg(200, int(cooldown.Milliseconds() - timeSinceLastMove.Milliseconds()))
					}
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
		// Check the response message
		fmt.Println("byt:", string(byt))
		api.Mux.Lock()
		if err := c.Conn.WriteMessage(gwebsocket.TextMessage, byt); err != nil {
			log.Println(err)
		}
		api.Mux.Unlock()
	}
}
func makeChangeClientMessage(x int, y int, r int, g int, b int, userID string) []byte {
	msg := common.ChangeClientPixelMsg{
		Type:   common.ChangeClientPixel,
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
		Type: common.DrawResponse,
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
		Type:   common.DrawResponse,
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
		Type:   common.VerificationFail,
		Status: s,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}


func makeUserLoginResponseMsg(s int, c int) []byte {
	msg := UserLoginResponseMsg{
		Type:     common.UserLoginResponse,
		Status:   s,
		Cooldown: c,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}
