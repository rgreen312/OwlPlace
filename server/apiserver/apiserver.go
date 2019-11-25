package apiserver

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/rgreen312/owlplace/server/common"
	"github.com/rgreen312/owlplace/server/consensus"
	"github.com/rgreen312/owlplace/server/wsutil"
)

var (
	ConfigurationError = errors.New("invalid apiserver configuration")
)

type ApiServer struct {
	config   *common.ServerConfig
	upgrader *websocket.Upgrader
	pool     *wsutil.Pool
	cons     consensus.IConsensus
}

func NewApiServer(servers map[int]*common.ServerConfig, nodeId int) (*ApiServer, error) {

	conf, ok := servers[nodeId]
	if !ok {
		return nil, errors.Wrapf(ConfigurationError, "missing entry for node: %d", nodeId)
	}

	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()

	cons, err := consensus.NewConsensusService(servers, nodeId, pool.Broadcast)
	if err != nil {
		return nil, errors.Wrap(err, "creating ConsensusService")
	}

	err = cons.Start()
	if err != nil {
		return nil, errors.Wrap(err, "starting ConsensusService")
	}

	log.WithFields(log.Fields{
		"server config":     conf,
		"consensus service": cons,
	}).Debug()

	return &ApiServer{
		config: conf,
		pool:   pool,
		cons:   cons,
	}, nil
}

func (api *ApiServer) ListenAndServe() {
	go api.pool.Run()

	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsutil.ServeWs(api.pool, api.cons, w, r)
	})
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.config.ApiPort), nil)
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

	img, err := api.cons.SyncGetImage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodedString := common.Base64Encode(img)
	msg := common.ImageMsg{
		Type:         common.Image,
		FormatString: encodedString,
	}

	log.WithFields(log.Fields{
		"ImageMsg": msg,
	}).Debug("constructed websocket message")

    js, err := json.Marshal(msg)
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

	err = api.cons.SyncUpdatePixel(msg.X, msg.Y, msg.R, msg.G, msg.B, msg.A)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


