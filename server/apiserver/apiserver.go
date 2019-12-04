package apiserver

import (
	"encoding/json"
	"fmt"

	"net/http"

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
	nodeID  uint64
	address string
	pool    *wsutil.Pool
	mp      consensus.MembershipProvider
	cons    consensus.IConsensus
}

func NewApiServer(nodeID uint64, address string, mp consensus.MembershipProvider) (*ApiServer, error) {
	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()
	go pool.Run()

	return &ApiServer{
		nodeID:  nodeID,
		address: address,
		pool:    pool,
		mp:      mp,
	}, nil
}

func (api *ApiServer) ListenAndServe() error {

	log.WithFields(log.Fields{
		"api address": fmt.Sprintf("%s:%d", api.address, common.ApiPort),
		"nodeID":      api.nodeID,
	}).Info("owlplace is listening for a trigger to form a dragonboat cluster")

	http.HandleFunc("/", api.HealthCheck)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", api.ServeWs)
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/consensus_trigger", func(w http.ResponseWriter, req *http.Request) {
		api.startConsensus(false, w, req)
	})
	http.HandleFunc("/consensus_join_message", func(w http.ResponseWriter, req *http.Request) {
		api.startConsensus(true, w, req)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", common.ApiPort), nil)
}

// startConsensus is an HTTP endpoint which indicates this node should start
// it's ConsensusService.  Additionally, this wrapper takes in a flag
// indicating whether this node is being invited to join an existing cluster,
// or whether it should start it's own.
func (api *ApiServer) startConsensus(joinExistingCluster bool, w http.ResponseWriter, req *http.Request) {
	var err error
	api.cons, err = consensus.NewConsensusService(api.mp, api.nodeID, api.address, api.pool.Broadcast)
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating the ConsensusService").Error(), http.StatusInternalServerError)
		return
	}

	err = api.cons.Start(joinExistingCluster)
	if err != nil {
		http.Error(w, errors.Wrap(err, "starting the ConsensusService").Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (api *ApiServer) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func (api *ApiServer) ServeWs(w http.ResponseWriter, req *http.Request) {
	wsutil.ServeWs(api.pool, api.cons, w, req)
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
	msg := map[string]string{
		"data": encodedString,
	}

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
