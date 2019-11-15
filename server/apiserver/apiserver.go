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

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
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
}

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
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/get/image", api.GetImageJson)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.HandleFunc("/ws", api.wsEndpoint)

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

func (api *ApiServer) GetImageJson(w http.ResponseWriter, req *http.Request) {
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

func (api *ApiServer) GetImage(w http.ResponseWriter, req *http.Request) {
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

func (api *ApiServer) UpdatePixel(w http.ResponseWriter, req *http.Request) {
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

				// TODO(user team): add user verification here
				err := api.conService.SyncUpdatePixel(dp_msg.X, dp_msg.Y, dp_msg.R, dp_msg.G, dp_msg.B, AlphaMask)
				if err != nil {
					// TODO(backend team): handle error response
				}
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
