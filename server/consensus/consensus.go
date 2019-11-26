package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"os"
	"net/http"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/rgreen312/owlplace/server/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	ClusterID     uint64 = 128
	SyncOpTimeout        = 3 * time.Second
	logLevel             = logger.ERROR
)

var (
	DragonboatConfigurationError = errors.New("dragonboat configuration")
	NoSuchUser                   = errors.New("no such user")
)

type IConsensus interface {
	SyncGetImage() (*image.RGBA, error)
	SyncUpdatePixel(x, y, r, g, b, a int) error
	SyncGetLastUserModification(userId string) (*time.Time, error)
	SyncSetLastUserModification(userId string, timestamp time.Time) error
}

type ConsensusService struct {
	nh         *dragonboat.NodeHost
	config     string
	raftConfig config.Config
	dkv        *DiskKV
	nodeId     uint64
	clusterId  uint64
	Broadcast  chan []byte
	// TODO: pull this out when we start using the kubernetes discovery
	// service.
	peers map[uint64]string
}

func NewConsensusService(servers map[uint64]string, nodeId uint64, broadcast chan []byte) (*ConsensusService, error) {

	conf, ok := servers[nodeId]
	if !ok {
		return nil, errors.Wrapf(DragonboatConfigurationError, "NodeID provided (%d) not present in server map.", nodeId)
	}

	nodeAddr := fmt.Sprintf("%s:%d", conf, common.ConsensusPort)

	// https://github.com/golang/go/issues/17393
	if runtime.GOOS == "darwin" {
		signal.Ignore(syscall.Signal(0xd))
	}

	peers := make(map[uint64]string)
	if len(servers) > 1 {
		for id, srv := range servers {
			peers[uint64(id)] = fmt.Sprintf("%s:%d", srv, common.ConsensusPort)
		}
	}

	log.WithFields(log.Fields{
		"node address": nodeAddr,
		"node id":      nodeId,
		"peers":        peers,
	}).Debug()

	// dragonboat provides it's own logging utilities.
	logger.GetLogger("raft").SetLevel(logLevel)
	logger.GetLogger("rsm").SetLevel(logLevel)
	logger.GetLogger("transport").SetLevel(logLevel)
	logger.GetLogger("grpc").SetLevel(logLevel)
	rc := config.Config{
		NodeID:             uint64(nodeId),
		ClusterID:          ClusterID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	log.WithFields(log.Fields{
		"raft config": rc,
	}).Debug("Dragonboat Configuration")
	if err := rc.Validate(); err != nil {
		return nil, err
	}

	datadir := filepath.Join(
		"example-data",
		"helloworld-data",
		fmt.Sprintf("node%d", nodeId))

	nhc := config.NodeHostConfig{
		DeploymentID:   1,
		WALDir:         datadir,
		NodeHostDir:    datadir,
		RTTMillisecond: 200,
		RaftAddress:    nodeAddr,
	}
	log.WithFields(log.Fields{
		"node host config": nhc,
	}).Debug("Dragonboat Configuration")
	if err := nhc.Validate(); err != nil {
		return nil, err
	}

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		return nil, errors.Wrap(err, "creating dragonboat nodehost")
	}

	return &ConsensusService{
		nh:         nh,
		dkv:        NewDiskKV(ClusterID, uint64(nodeId), broadcast),
		config:     conf,
		nodeId:     nodeId,
		clusterId:  ClusterID,
		peers:      peers,
		raftConfig: rc,
		Broadcast:  broadcast,
	}, nil
}

func (cs *ConsensusService) SyncGetImage() (*image.RGBA, error) {
	img := cs.dkv.GetInMemoryImage()
	return &img, nil
}

func (cs *ConsensusService) SyncUpdatePixel(x, y, r, g, b, a int) error {
	// Create the kv pair to send to dragonboat
	kv := &KVData{
		Key: fmt.Sprintf("pixel(%d,%d)", x, y),
		Val: fmt.Sprintf("(%d,%d,%d,%d)", r, g, b, a),
	}

	data, err := json.Marshal(kv)
	if err != nil {
		return errors.Wrap(err, "marshalling update pixel kv data")
	}

	// sync with dragonboat
	// TODO(gabe): determine if we should validate / check the result
	// (currently not using it.)
	session := cs.nh.GetNoOPSession(cs.clusterId)
	ctx, _ := context.WithTimeout(context.Background(), SyncOpTimeout)
	_, err = cs.nh.SyncPropose(ctx, session, data)
	if err != nil {
		return errors.Wrap(err, "syncing with dragonboat")
	}

	return nil
}

func ScanDiscoveryService(servers map[uint64]string, nh *dragonboat.NodeHost){
	for {

		fmt.Fprintf(os.Stdout, "Scanning Discovery Service\n")
		
		//Actually scan discovery service
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

		for _, pod := range pods.Items {
			if _, ok := servers[common.IPToNodeId(pod.Status.PodIP)]; !ok {

				fmt.Fprintf(os.Stdout, "Found pod that's not in cluster\n")

				// Adding pod to server map
				nodeId := common.IPToNodeId(pod.Status.PodIP)
				servers[nodeId] = pod.Status.PodIP

				// Adding pod to cluster 
				request_data, request_err := nh.RequestAddNode(ClusterID, uint64(nodeId), fmt.Sprintf("%s:%d", pod.Status.PodIP, common.ConsensusPort), 0, 1000*time.Millisecond)
				if(request_err != nil){
					panic(request_err)
				}

				// Wait for response or timeout
				results := <-request_data.CompletedC

				if(results.Completed()){
					fmt.Fprintf(os.Stdout, "Pod join success\n")
					// Send an http join request to the other nodes
					_, err := http.Get(fmt.Sprintf("http://%s:%d/consensus_join_message", pod.Status.PodIP , common.ApiPort))
					if(err != nil){
						panic(err)
					}
				} else {
					fmt.Fprintf(os.Stdout, "Pod join failure\n")
				}

	
			}
			
	    }

		// Wait for 10 seconds before scanning again
		time.Sleep(10000 * time.Millisecond)

	}
}


func (cs *ConsensusService) SyncGetLastUserModification(userId string) (*time.Time, error) {

	// Request a ready from dragonboat
	ctx, _ := context.WithTimeout(context.Background(), SyncOpTimeout)
	result, err := cs.nh.SyncRead(ctx, cs.clusterId, []byte(userId))
	if err != nil {
		return nil, errors.Wrap(err, "reading from dragonboat")
	}

	resultString := string(result.([]byte))
	if resultString == "" {
		return nil, NoSuchUser
	}

	t, err := time.Parse(common.TimeFormat, resultString)
	if err != nil {
		return nil, errors.Wrap(err, "parsing time returned from dragonboat")
	}

	return &t, nil
}

func (cs *ConsensusService) SyncSetLastUserModification(userId string, timestamp time.Time) error {

	// Create the kv pair to send to dragonboat
	kv := &KVData{
		Key: userId,
		Val: timestamp.Format(common.TimeFormat),
	}

	data, err := json.Marshal(kv)
	if err != nil {
		return errors.Wrap(err, "marshalling update pixel kv data")
	}

	// Sync with dragonboat
	ctx, _ := context.WithTimeout(context.Background(), SyncOpTimeout)
	session := cs.nh.GetNoOPSession(cs.clusterId)
	_, err = cs.nh.SyncPropose(ctx, session, data)
	if err != nil {
		return errors.Wrap(err, "syncing with dragonboat")
	}

	return nil
}

func (cs *ConsensusService) Start(join bool) error {
	// Function to provide a state-machine reference to Raft
	stateMachineProvider := func(clusterId uint64, nodeId uint64) sm.IOnDiskStateMachine {
		return cs.dkv
	}

	peers := make(map[uint64]string)
	if(!join){
		peers = cs.peers
	}

	log.WithFields(log.Fields{
		"peers":                cs.peers,
		"join":                 join,
		"stateMachineProvider": stateMachineProvider,
		"raft config":          cs.raftConfig,
	}).Debug("starting on disk cluster")
	err := cs.nh.StartOnDiskCluster(peers, join, stateMachineProvider, cs.raftConfig)
	go ScanDiscoveryService(cs.peers, cs.nh)
	return err
}

// TODO(gabe): figure out how to shutdown a dragonboat node
func (cs *ConsensusService) Stop() error {
	return nil
}
