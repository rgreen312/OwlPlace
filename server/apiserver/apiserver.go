package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	"github.com/rgreen312/owlplace/server/wsutil"
)

const (
	wsUpgraderReadBufferSize  = 1024
	wsUpgraderWriteBufferSize = 1024
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
	upgrader   *websocket.Upgrader
	pool       *wsutil.Pool
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
		config: conf,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  wsUpgraderReadBufferSize,
			WriteBufferSize: wsUpgraderWriteBufferSize,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		pool:       wsutil.NewPool(),
		conService: conService,
	}, nil
}

func (api *ApiServer) ListenAndServe() {
	go api.pool.Start()

	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", api.Upgrade)
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.config.ApiPort), nil)
}

func (api *ApiServer) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	client := wsutil.NewClient(conn, api.pool)

	// Register the upgraded connection with our pool.
	api.pool.Register <- client

	// Start communicating with the frontend on the current goroutine.  Note
	// that these HTTP handlers are in fact spawned on their own routines,
	// rather than the main goroutine handling new requests.
	client.Start()
}

// HTTPGetImageJson provides a synchronous method with which to request the
// canvas.  It returns a JSON object structured as:
//
//      {
//          "data": "... base64 encoded rgba png ..."
//      }
//
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

func (api *ApiServer) HTTPUpdatePixel(w http.ResponseWriter, req *http.Request) {
	msg, err := common.NewDrawPixelMsg(req)
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

// base64Encode returns a base64 string representation of an RGBA image.
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
