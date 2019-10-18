package consensus

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	// "context"
	// "time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/lni/goutils/syncutil"
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

type ConsensusMessage struct {
	Type int
	Data string
}

func NewImageMessage(img image.RGBA) ConsensusMessage {
	// In-memory buffer to store PNG image
	// before we base 64 encode it
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, &img)

	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	return ConsensusMessage{
		Type: GET_IMAGE,
		Data: encodedString,
	}
}

func GetTimestampMessage(timestamp string) ConsensusMessage {

	return ConsensusMessage{
		Type: GET_LAST_USER_UPDATE,
		Data: timestamp,
	}
}

func SuccessMessage() ConsensusMessage {
	return ConsensusMessage{
		Type: SUCCESS,
		Data: "Success\n",
	}
}

func FailureMessage() ConsensusMessage {
	return ConsensusMessage{
		Type: FAILURE,
		Data: "Failure\n",
	}
}

type BackendMessage struct {
	Type int
	Data string
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

func parseCommand(msg string) (RequestType, string, string, bool) {
	parts := strings.Split(strings.TrimSpace(msg), " ")
	if len(parts) == 0 || (parts[0] != "put" && parts[0] != "get") {
		return PUT, "", "", false
	}
	if parts[0] == "put" {
		if len(parts) != 3 {
			return PUT, "", "", false
		}
		return PUT, parts[1], parts[2], true
	}
	if len(parts) != 2 {
		return GET, "", "", false
	}
	return GET, parts[1], "", true
}

func MainConsensus(recvc chan BackendMessage, sendc chan ConsensusMessage) {

	nodeID := flag.Int("nodeid", 1, "NodeID to use")
	addr := flag.String("addr", "", "Nodehost address")
	join := flag.Bool("join", false, "Joining a new node")
	flag.Parse()
	if len(*addr) == 0 && *nodeID != 1 && *nodeID != 2 && *nodeID != 3 {
		fmt.Fprintf(os.Stderr, "node id must be 1, 2 or 3 when address is not specified\n")
		os.Exit(1)
	}
	// https://github.com/golang/go/issues/17393
	if runtime.GOOS == "darwin" {
		signal.Ignore(syscall.Signal(0xd))
	}
	peers := make(map[uint64]string)
	if !*join {
		for idx, v := range addresses {
			peers[uint64(idx+1)] = v
		}
	}
	var nodeAddr string
	if len(*addr) != 0 {
		nodeAddr = *addr
	} else {
		nodeAddr = peers[uint64(*nodeID)]
	}
	fmt.Fprintf(os.Stdout, "node address: %s\n", nodeAddr)
	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	rc := config.Config{
		NodeID:             uint64(*nodeID),
		ClusterID:          exampleClusterID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	datadir := filepath.Join(
		"example-data",
		"helloworld-data",
		fmt.Sprintf("node%d", *nodeID))
	nhc := config.NodeHostConfig{
		WALDir:         datadir,
		NodeHostDir:    datadir,
		RTTMillisecond: 200,
		RaftAddress:    nodeAddr,
	}
	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		panic(err)
	}
	var imgGetter func() image.RGBA
	stateMachineProvider := func(clusterID uint64, nodeID uint64) sm.IOnDiskStateMachine {
		dkv := NewDiskKV(clusterID, nodeID).(*DiskKV)
		imgGetter = dkv.GetInMemoryImage
		return dkv
	}
	if err := nh.StartOnDiskCluster(peers, *join, stateMachineProvider, rc); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add cluster, %v\n", err)
		os.Exit(1)
	}
	raftStopper := syncutil.NewStopper()
	raftStopper.RunWorker(func() {
		cs := nh.GetNoOPSession(exampleClusterID)
		for {
			select {
			case backend_msg, ok := <-recvc:
				fmt.Println("\n Consensus received the backend_msg")
				if !ok {
					return
				}
				switch backend_msg.Type {
				case GET_IMAGE:
					sendc <- NewImageMessage(imgGetter())
				case UPDATE_PIXEL, SET_LAST_USER_UPDATE:
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					_, key, val, ok := parseCommand(backend_msg.Data)
					if !ok {
						sendc <- FailureMessage()
					} else {
						kv := &KVData{
							Key: key,
							Val: val,
						}
						data, err := json.Marshal(kv)
						if err != nil {
							panic(err)
						}
						_, err = nh.SyncPropose(ctx, cs, data)
						if err != nil {
							fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
						}
						sendc <- SuccessMessage()
					}
					cancel()

				case GET_LAST_USER_UPDATE:
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					_, key, _, ok := parseCommand(backend_msg.Data)
					if !ok {
						sendc <- FailureMessage()
					}
					result, err := nh.SyncRead(ctx, exampleClusterID, []byte(key))
					if err != nil {
						sendc <- FailureMessage()
					} else {
						sendc <- GetTimestampMessage(result.(string))
					}
					cancel()
				}

			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}
