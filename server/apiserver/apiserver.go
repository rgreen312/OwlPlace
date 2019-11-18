package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"math"
	"net/http"
	"strconv"
	"time"

	gwebsocket "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
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
	config     *common.ServerConfig
	conService *consensus.ConsensusService
	Mux        sync.Mutex
}

type Client struct {
	ID   string
	Conn *gwebsocket.Conn
	Pool *Pool
}

var cooldown, _ = time.ParseDuration("5m")

func NewApiServer(servers map[int]*common.ServerConfig, nodeId int) (*ApiServer, error) {

	conf, ok := servers[nodeId]
	if !ok {
		return nil, errors.Wrapf(configError, "missing entry for node: %d", nodeId)
	}

	conService, err := consensus.NewConsensusService(servers, nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "creating ConsensusService")
	}

	err = conService.Start()
	if err != nil {
		return nil, errors.Wrap(err, "starting ConsensusService")
	}

	log.WithFields(log.Fields{
		"server config":     conf,
		"consensus service": conService,
	}).Debug()

	return &ApiServer{
		config:     conf,
		conService: conService,
	}, nil
}

func (api *ApiServer) ListenAndServe() {
	pool := NewPool()
	go pool.Start(api.Mux)
	http.HandleFunc("/get_image", api.HTTPGetImage)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.serveWs(pool, w, r)
	})
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/update_user", api.HTTPUserLogin)

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.config.ApiPort), nil)
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

func (api *ApiServer) HTTPUserLogin(w http.ResponseWriter, req *http.Request) []byte {
	userID := req.URL.Query().Get("id")
	//if userID == "" {
	//	http.Error(w, errors.New("empty param: userID").Error(), http.StatusInternalServerError)
	//	return
	//}
	byt := makeUserLoginResponseMsg(400, -1)
	lastMove, getErr := api.conService.SyncGetLastUserModification(userID)
	if getErr == consensus.NoSuchUser {
		setErr := api.conService.SyncSetLastUserModification(userID, time.Unix(0,0)) // Default Timestamp for New Users
		if setErr != nil {
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
	return byt
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
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
				// <Start Validate User>
				lastMove, getErr := api.conService.SyncGetLastUserModification()
				// <End Validate User>

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
				api.Mux.Lock()
				c.Pool.Broadcast <- ccpMsg
				api.Mux.Unlock()
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
			fmt.Println(cu_msg)
			if err := json.Unmarshal(p, &cu_msg); err == nil {
				log.WithFields(log.Fields{
					"message": cu_msg,
				}).Debug("received ws message")

				userID := cu_msg.Id
				//if userID == "" {
				//	http.Error(w, errors.New("empty param: userID").Error(), http.StatusInternalServerError)
				//	return
				//}
				byt := makeUserLoginResponseMsg(400, -1)
				lastMove, getErr := api.conService.SyncGetLastUserModification(userID)
				if getErr == consensus.NoSuchUser {
					setErr := api.conService.SyncSetLastUserModification(userID, time.Unix(0,0)) // Default Timestamp for New Users
					if setErr != nil {
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
		log.Println(byt)

		api.Mux.Lock()
		if err := c.Conn.WriteMessage(gwebsocket.TextMessage, byt); err != nil {
			log.Println(err)
		}
		api.Mux.Unlock()

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

func makeUserLoginResponseMsg(s int, c int) []byte {
	msg := UserLoginResponseMsg{
		Type:     UserLoginResponse,
		Status:   s,
		Cooldown: c,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}
