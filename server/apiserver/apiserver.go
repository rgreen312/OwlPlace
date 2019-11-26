package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
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
	pod_ip     string
	node_id     uint64
	upgrader *websocket.Upgrader
	pool     *wsutil.Pool
	cons     consensus.IConsensus
}

func NewApiServer(pod_ip string) (*ApiServer, error) {
	return &ApiServer{
		pod_ip: pod_ip,
		node_id: common.IPToNodeId(pod_ip),
	}, nil
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

	servers :=  make(map[uint64]string)
	for _, pod := range pods.Items {
		if(pod.Status.PodIP != api.pod_ip && !join){
			// Send an http join request to the other nodes
			_, err := http.Get(fmt.Sprintf("http://%s:%d/consensus_join_message", pod.Status.PodIP, common.ApiPort))
			if(err != nil){
				panic(err)
			}
		}
		servers[common.IPToNodeId(pod.Status.PodIP)] = pod.Status.PodIP
			
    }
	//At first, just print something so that we know http requests are working inside kubernetes
	fmt.Fprintf(os.Stdout, "Pod Trigger Called\n")


	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()
	api.pool = pool
	go api.pool.Run()

	// Start the consensus service in the background
	conService, err := consensus.NewConsensusService(servers, api.node_id, pool.Broadcast)
	api.cons = conService

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

func (api *ApiServer) HealthCheck(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(os.Stdout, "HealthCheck\n")
	w.WriteHeader(200)
}


func (api *ApiServer) ListenAndServe() {
	http.HandleFunc("/", api.HealthCheck)
	http.HandleFunc("/json/image", api.HTTPGetImageJson)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsutil.ServeWs(api.pool, api.cons, w, r)
	})
	http.HandleFunc("/update_pixel", api.HTTPUpdatePixel)
	http.HandleFunc("/consensus_trigger", api.ConsensusTrigger)
	http.HandleFunc("/consensus_join_message", api.ConsensusJoinMessage)

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", common.ApiPort), nil)
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

	encodedString := base64Encode(img)
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
