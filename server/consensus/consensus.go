package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/rgreen312/owlplace/server/common"
)

const (
	ClusterID                    uint64 = 128
	SyncOpTimeout                       = 3 * time.Second
	logLevel                            = logger.ERROR
	intervalDiscoveryServiceScan        = 60 * time.Second
)

var (
	DragonboatConfigurationError = errors.New("dragonboat configuration")
	NoSuchUser                   = errors.New("no such user")
)

type IConsensus interface {
	Start(join bool) error
	SyncGetImage() (*image.RGBA, error)
	SyncUpdatePixel(x, y, r, g, b, a int) error
	SyncGetLastUserModification(userId string) (*time.Time, error)
	SyncSetLastUserModification(userId string, timestamp time.Time) error
}

type ConsensusService struct {
	nh              *dragonboat.NodeHost
	raftConfig      config.Config
	dkv             *DiskKV
	nodeId          uint64
	raftAddress     string
	apiAddress      string
	clusterId       uint64
	Broadcast       chan []byte
	mp              MembershipProvider
	discoveryCloser chan bool
}

func NewConsensusService(mp MembershipProvider, nodeId uint64, address string, broadcast chan []byte) (*ConsensusService, error) {

	raftAddress := common.SetAddressPort(address, common.ConsensusPort)
	apiAddress := common.SetAddressPort(address, common.ApiPort)

	// https://github.com/golang/go/issues/17393
	if runtime.GOOS == "darwin" {
		signal.Ignore(syscall.Signal(0xd))
	}

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
		RaftAddress:    raftAddress,
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
		nh:          nh,
		dkv:         NewDiskKV(ClusterID, uint64(nodeId), broadcast),
		nodeId:      nodeId,
		raftAddress: raftAddress,
		apiAddress:  apiAddress,
		clusterId:   ClusterID,
		raftConfig:  rc,
		Broadcast:   broadcast,
		mp:          mp,
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

func (cs *ConsensusService) scanDiscoveryService() {
	ticker := time.NewTicker(intervalDiscoveryServiceScan)
	cs.discoveryCloser = make(chan bool)
	go func() {
		for {
			select {
			case <-cs.discoveryCloser:
				return
			case <-ticker.C:
				log.Debug("scanning discovery service")

				desiredMembers, err := cs.mp.GetMembership()
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("retrieving desired membership")
					continue
				}

				ctx, _ := context.WithTimeout(context.TODO(), 3000*time.Millisecond)
				currentMembership, err := cs.nh.SyncGetClusterMembership(ctx, ClusterID)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("retrieving current membership from dragonboat")
					continue
				}
				currentMembers := currentMembership.Nodes

				for nodeID, address := range desiredMembers {
					if _, ok := currentMembers[nodeID]; ok {
						// TODO: should we check to make sure the address is
						// the same?  otherwise, we can simply continue here as
						// we only want to add new nodes here.
						continue
					}

					err := cs.requestAddNode(nodeID, address)
					if err != nil {
						log.WithFields(log.Fields{
							"err": err,
						}).Error("adding node to dragonboat cluster")
						continue
					}

				}
			}
		}
	}()
}

// requestAddNodes requests to add a node to the consensus service.
func (cs *ConsensusService) requestAddNode(nodeID uint64, address string) error {

	// First, signal the node to indicate it should start its consensus
	// service.
	_, err := http.Get(fmt.Sprintf("http://%s:%d/consensus_join_message", strings.Split(address, ":")[0], common.ApiPort))
	if err != nil {
		return errors.Wrap(err, "sending consensus join message")
	}

	requestData, err := cs.nh.RequestAddNode(ClusterID, uint64(nodeID), address, 0, 10*time.Second)
	if err != nil {
		log.WithFields(log.Fields{
			"nodeID":  nodeID,
			"address": address,
			"err":     err,
		}).Debug("adding node to cluster")
		return errors.Wrapf(err, "adding node (ID=%d) to cluster", nodeID)
	}

	// Wait for response or timeout
	results := <-requestData.CompletedC

	if !results.Completed() {
		log.WithFields(log.Fields{
			"nodeID":  nodeID,
			"address": address,
			"result":  results.GetResult(),
		}).Error("failed to add node to cluster")

		return errors.Wrapf(errors.New("failed to add node to cluster"), "request results: %+v", results.GetResult())
	}

	log.WithFields(log.Fields{
		"nodeID":  nodeID,
		"address": address,
	}).Debug("successfully added node to cluster")

	return nil
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

	initialMembers := map[uint64]string{
		cs.nodeId: cs.raftAddress,
	}
	if join {
		initialMembers = make(map[uint64]string)
	}

	log.WithFields(log.Fields{
		"initialMembers":       initialMembers,
		"join":                 join,
		"stateMachineProvider": stateMachineProvider,
		"raft config":          cs.raftConfig,
	}).Debug("starting on disk cluster")

	err := cs.nh.StartOnDiskCluster(initialMembers, join, stateMachineProvider, cs.raftConfig)
	if err != nil {
		return err
	}

	// if we successfully started the cluster, begin scanning for new peers
	cs.scanDiscoveryService()

	return nil
}

// TODO(gabe): figure out how to shutdown a dragonboat node
func (cs *ConsensusService) Stop() error {
	// shut down the goroutine responsible for scanning k8s for new pods
	cs.discoveryCloser <- true
	return nil
}
