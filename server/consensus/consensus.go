package consensus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"os/signal"
	"path/filepath"
	"runtime"
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

type RequestType uint64

const (
	ClusterId uint64 = 128
)

const (
	PUT RequestType = iota
	GET
)

const (
	GET_IMAGE            int = 0
	UPDATE_PIXEL         int = 1
	ADD_USER             int = 2
	GET_LAST_USER_UPDATE int = 3
	SET_LAST_USER_UPDATE int = 4
	SUCCESS              int = 5
	FAILURE              int = 6
)

const (
	DRAGONBOAT_ERROR int = 0
	MESSAGE_ERROR    int = 1
)

type ConsensusMessage struct {
	Type int
	Data bytes.Buffer
}

type BackendMessage struct {
	Type int
	Data bytes.Buffer
}

type GetUserDataBackendMessage struct {
	UserId string
}

type SetUserDataBackendMessage struct {
	UserId, Timestamp string
}

type UpdatePixelBackendMessage struct {
	X, Y, R, G, B, A string
}

const (
	SyncOpTimeout = 3 * time.Second
)

var (
	DragonboatConfigurationError = errors.New("dragonboat configuration")
	NoSuchUser                   = errors.New("no such user")
)

type IConsensus interface {
	SyncGetImage() (*image.RGBA, error)
	SyncGetLastUserModification(userId string) (time.Time, error)
	SyncSetLastUserModification(userId string, timestamp time.Time) error
}

type ConsensusService struct {
	nh         *dragonboat.NodeHost
	config     *common.ServerConfig
	raftConfig config.Config
	dkv        *DiskKV
	nodeId     int
	clusterId  uint64
	// TODO: pull this out when we start using the kubernetes discovery
	// service.
	peers map[uint64]string
}

func NewConsensusService(servers map[int]*common.ServerConfig, nodeId int) (*ConsensusService, error) {

	conf, ok := servers[nodeId]
	if !ok {
		return nil, errors.Wrapf(DragonboatConfigurationError, "NodeID provided (%d) not present in server map.", nodeId)
	}

	nodeAddr := fmt.Sprintf("%s:%d", conf.Hostname, conf.ConsensusPort)

	// https://github.com/golang/go/issues/17393
	if runtime.GOOS == "darwin" {
		signal.Ignore(syscall.Signal(0xd))
	}

	peers := make(map[uint64]string)
	if len(servers) > 1 {
		for id, srv := range servers {
			peers[uint64(id)] = fmt.Sprintf("%s:%d", srv.Hostname, srv.ConsensusPort)
		}
	}

	log.WithFields(log.Fields{
		"node address": nodeAddr,
		"node id":      nodeId,
		"peers":        peers,
	}).Debug()

	// dragonboat provides it's own logging utilities.
	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	rc := config.Config{
		NodeID:             uint64(nodeId),
		ClusterID:          ClusterId,
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
		dkv:        NewDiskKV(ClusterId, uint64(nodeId)),
		config:     conf,
		nodeId:     nodeId,
		clusterId:  ClusterId,
		peers:      peers,
		raftConfig: rc,
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

func (cs *ConsensusService) Start() error {
	// For more information on the join parameter, see:
	// https://godoc.org/github.com/lni/dragonboat#NodeHost.StartCluster
	join := false

	// Function to provide a state-machine reference to Raft
	stateMachineProvider := func(clusterId uint64, nodeId uint64) sm.IOnDiskStateMachine {
		return cs.dkv
	}
	log.WithFields(log.Fields{
		"peers":                cs.peers,
		"join":                 join,
		"stateMachineProvider": stateMachineProvider,
		"raft config":          cs.raftConfig,
	}).Debug("starting on disk cluster")
	return cs.nh.StartOnDiskCluster(cs.peers, join, stateMachineProvider, cs.raftConfig)
}

// TODO(gabe): figure out how to shutdown a dragonboat node
func (cs *ConsensusService) Stop() error {
	return nil
}
