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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ConfigurationError = errors.New("invalid apiserver configuration")
)

type ApiServer struct {
	pod_ip  string
	node_id uint64
	pool    *wsutil.Pool
	cons    consensus.IConsensus
}

func NewApiServer(pod_ip string) (*ApiServer, error) {
	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()
	go pool.Run()

	nodeID, err := common.IPToNodeId(pod_ip)
	log.WithFields(log.Fields{
		"pod ip":  pod_ip,
		"node id": nodeID,
	})

	return &ApiServer{
		pod_ip:  pod_ip,
		node_id: nodeID,
		pool:    pool,
	}, err
}

func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/", api.HealthCheck)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", api.ServeWs)
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/consensus_trigger", func(w http.ResponseWriter, req *http.Request) {
		api.httpConsensusTrigger(false, w, req)
	})
	http.HandleFunc("/consensus_join_message", func(w http.ResponseWriter, req *http.Request) {
		api.httpConsensusTrigger(true, w, req)
	})

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", common.ApiPort), nil)
}

// startConsensus attempts to start the consensus module with a list of peers
// collected from k8s.
func (api *ApiServer) startConsensus(join bool) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "retrieving cluster config")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating API handle")
	}

	pods, err := clientset.CoreV1().Pods("dev").List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "listing dev pods")
	}

	log.WithFields(log.Fields{
		"numPods": len(pods.Items),
	}).Debug("finding cluster peers")

	servers := make(map[uint64]string)
	for _, pod := range pods.Items {
		if pod.Status.PodIP != api.pod_ip && !join {
			// Send an http join request to the other nodes
			_, err := http.Get(fmt.Sprintf("http://%s:%d/consensus_join_message", pod.Status.PodIP, common.ApiPort))
			if err != nil {
				return errors.Wrapf(err, "sending join request to server: %s:%d", pod.Status.PodIP, common.ApiPort)
			}
		}
	  nodeid, _ := common.IPToNodeId(pod.Status.PodIP)
		servers[nodeid] = pod.Status.PodIP
	}

	// Start the consensus service in the background
	api.cons, err = consensus.NewConsensusService(servers, api.node_id, api.pool.Broadcast)
	if err != nil {
		return errors.Wrap(err, "creating the ConsensusService")
	}

	err = api.cons.Start(join)
	if err != nil {
		return errors.Wrap(err, "starting ConsensusService")
	}

	return nil
}

func (api *ApiServer) httpConsensusTrigger(join bool, w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"pod ip": api.pod_ip,
		"join":   join,
	}).Debug("ConsensusTrigger request received")

	err := api.startConsensus(join)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
