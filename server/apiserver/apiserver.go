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
	nodeID   uint64
	nodeAddr string
	pool     *wsutil.Pool
	mp       consensus.MembershipProvider
	cons     consensus.IConsensus
}

func NewApiServer(nodeID uint64, nodeAddr string, mp consensus.MembershipProvider) (*ApiServer, error) {
	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()
	go pool.Run()

	return &ApiServer{
		nodeID:   nodeID,
		nodeAddr: nodeAddr,
		pool:     pool,
		mp:       mp,
	}, nil
}

func (api *ApiServer) ListenAndServe() error {

	log.WithFields(log.Fields{
		"api address": fmt.Sprintf("%s:%d", api.nodeAddr, common.ApiPort),
		"nodeID":      api.nodeID,
	}).Info("owlplace is listening for a trigger to form a dragonboat cluster")

	http.HandleFunc("/", api.HealthCheck)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", api.ServeWs)
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/consensus_trigger", api.startConsensus)
	http.HandleFunc("/consensus_join_message", api.joinConsensus)

	return http.ListenAndServe(fmt.Sprintf(":%d", common.ApiPort), nil)
}

// joinCluster is an HTTP endpoint which indicates that there is an existing
// Dragonboat cluster which can be joined.  This endpoint should not be called
// manually, but rather triggered by the member which receives a consensus
// trigger.
func (api *ApiServer) joinConsensus(w http.ResponseWriter, req *http.Request) {
	// Start the consensus service in the background
	var err error
	api.cons, err = consensus.NewConsensusService(api.mp, api.nodeID, api.pool.Broadcast)
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating the ConsensusService").Error(), http.StatusInternalServerError)
	}

	joinExistingCluster := true
	err = api.cons.Start(joinExistingCluster)
	if err != nil {
		http.Error(w, errors.Wrap(err, "starting the ConsensusService").Error(), http.StatusInternalServerError)
	}

	// Otherwise, indicate a successful join.
	w.WriteHeader(200)
}

// startConsensus attempts to start the consensus module with a list of peers
// collected from k8s.
func (api *ApiServer) startConsensus(w http.ResponseWriter, req *http.Request) {

	// Start the consensus service in the background
	var err error
	api.cons, err = consensus.NewConsensusService(api.mp, api.nodeID, api.pool.Broadcast)
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating the ConsensusService").Error(), http.StatusInternalServerError)
        return
	}

	// This parameter indicates that we're not joining an existing cluster, but
	// forming a new one.
	joinExistingCluster := false
	err = api.cons.Start(joinExistingCluster)
	if err != nil {
		http.Error(w, errors.Wrap(err, "starting the ConsensusService").Error(), http.StatusInternalServerError)
        return
	}

	// Send a join request to the remaining members.
	servers, err := api.mp.GetMembership()

	// servers contains a mapping from nodeID to a full consensus URI, e.g.
	// domain:consensusPort.  We wish to make an HTTP request against the API
	// port.
	for _, conAddress := range servers {
		apiAddress := common.ReplacePort(conAddress, common.ApiPort)
		// Send an HTTP join request to the other nodes
		_, err := http.Get(fmt.Sprintf("http://%s/consensus_join_message", apiAddress))
		if err != nil {
			// TODO: handle error.  I think we should collect all errors here
			// and then report failure using http.Error as demonstrated above.
			continue
		}
	}
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
