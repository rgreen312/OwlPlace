package consensus

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/lni/goutils/syncutil"

	"github.com/rgreen312/owlplace/server/common"
)

type RequestType uint64

const (
	exampleClusterID uint64 = 128
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
	DRAGONBOAT_ERROR     int = 0
	MESSAGE_ERROR        int = 1
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

func NewImageMessage(img image.RGBA) ConsensusMessage {
	// In-memory buffer to store PNG image
	// before we base 64 encode it
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	enc.Encode(img)

	return ConsensusMessage{
		Type: GET_IMAGE,
		Data: encoded_msg,
	}
}

func GetTimestampMessage(timestamp string) ConsensusMessage {
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	enc.Encode(timestamp)

	return ConsensusMessage{
		Type: GET_LAST_USER_UPDATE,
		Data: encoded_msg,
	}
}

func SuccessMessage() ConsensusMessage {
	var empty_buffer bytes.Buffer
	return ConsensusMessage{
		Type: SUCCESS,
		Data: empty_buffer,
	}
}

func FailureMessage(error_type int) ConsensusMessage {
	var encoded_msg bytes.Buffer
	enc := gob.NewEncoder(&encoded_msg)
	enc.Encode(error_type)

	return ConsensusMessage{
		Type: FAILURE,
		Data: encoded_msg,
	}
}

var (
	// initial nodes count is fixed to three, their addresses are also fixed
	addresses = []string{
		"localhost:63001",
		"localhost:63002",
		"localhost:63003",
	}
)

func printUsage() {
	fmt.Fprintf(os.Stdout, "Usage - \n")
	fmt.Fprintf(os.Stdout, "put key value\n")
	fmt.Fprintf(os.Stdout, "get key\n")
}

func MainConsensus(recvc chan BackendMessage, sendc chan ConsensusMessage, servers map[int]*common.ServerConfig, nodeId int) {

	conf := servers[nodeId]
	// For more information on the join parameter, see:
	// https://godoc.org/github.com/lni/dragonboat#NodeHost.StartCluster
	join := false
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

	// Log the node address for debugging purposes.
	fmt.Fprintf(os.Stdout, "node address: %s\n", nodeAddr)
	// Log the node ID
	fmt.Fprintf(os.Stdout, "node id: %d\n", nodeId)
	// Log the node address for debugging purposes.
	fmt.Fprintf(os.Stdout, "peers: %+v\n", peers)

	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	rc := config.Config{
		NodeID:             uint64(nodeId),
		ClusterID:          exampleClusterID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	if err := rc.Validate(); err != nil {
		log.Fatalf("Invalid Dragonboat Configuration: %s\n", err)
	} else {
		log.Printf("Node %d: Valid Dragonboat Configuration\n", nodeId)
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
	if err := nhc.Validate(); err != nil {
		log.Fatalf("Invalid Dragonboat Nodehost Configuration: %s\n", err)
	} else {
		log.Printf("Node %d: Valid Nodehost Dragonboat Configuration\n", nodeId)
	}
	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		log.Fatalf("Failed to create new nodehost: %s\n", err)
	}
	var imgGetter func() image.RGBA
	stateMachineProvider := func(clusterID uint64, nodeID uint64) sm.IOnDiskStateMachine {
		dkv := NewDiskKV(clusterID, nodeID).(*DiskKV)
		imgGetter = dkv.GetInMemoryImage
		return dkv
	}
	if err := nh.StartOnDiskCluster(peers, join, stateMachineProvider, rc); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add cluster, %v\n", err)
		os.Exit(1)
	}
	raftStopper := syncutil.NewStopper()
	raftStopper.RunWorker(func() {
		cs := nh.GetNoOPSession(exampleClusterID)
		for {
			select {
			case backend_msg, ok := <-recvc:
				if !ok {
					return
				}

				// Start background context
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

				// Message decoding depends on the type
				switch backend_msg.Type {
				case GET_IMAGE:
					sendc <- NewImageMessage(imgGetter())
				case UPDATE_PIXEL:
					// Decode the message from the glob
					dec := gob.NewDecoder(&backend_msg.Data)
					var umsg UpdatePixelBackendMessage
					err = dec.Decode(&umsg)
					if err != nil {
						sendc <- FailureMessage(1)
					}

					// Create the kv pair to send to dragonboat
					kv := &KVData{
						Key: fmt.Sprintf("pixel(%s,%s)", umsg.X, umsg.Y),
						Val: fmt.Sprintf("(%s,%s,%s,%s)", umsg.R, umsg.G, umsg.B, umsg.A),
					}
					data, err := json.Marshal(kv)
					if err != nil {
						sendc <- FailureMessage(1)
					}

					// Sync with dragonboat
					_, err = nh.SyncPropose(ctx, cs, data)
					if err != nil {
						fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
					}
					sendc <- SuccessMessage()
				case SET_LAST_USER_UPDATE:
					// Decode the message from the glob
					dec := gob.NewDecoder(&backend_msg.Data)
					var umsg SetUserDataBackendMessage
					err = dec.Decode(&umsg)
					if err != nil {
						sendc <- FailureMessage(1)
					}

					// Create the kv pair to send to dragonboat
					kv := &KVData{
						Key: umsg.UserId,
						Val: umsg.Timestamp,
					}
					data, err := json.Marshal(kv)
					if err != nil {
						sendc <- FailureMessage(1)
					}

					// Sync with dragonboat
					_, err = nh.SyncPropose(ctx, cs, data)
					if err != nil {
						fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
					}
					sendc <- SuccessMessage()

				case GET_LAST_USER_UPDATE:
					// Decode the message from the glob
					dec := gob.NewDecoder(&backend_msg.Data)
					var umsg GetUserDataBackendMessage
					err = dec.Decode(&umsg)
					if err != nil {
						sendc <- FailureMessage(1)
					}

					// Request a ready from dragonboat
					result, err := nh.SyncRead(ctx, exampleClusterID, []byte(umsg.UserId))
					if err != nil {
						sendc <- FailureMessage(0)
					} else {
						sendc <- GetTimestampMessage(result.(string))
					}

				}
				cancel()

			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}
