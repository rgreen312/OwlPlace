package consensus

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	// "context"
	// "time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
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
	GET_IMAGE int = 0
	UPDATE_PIXEL int = 1
	ADD_USER int = 2
	GET_LAST_USER_UPDATE int = 3
)

type ConsensusMessage struct {
	Type int
	Data string
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

func MainConsensus(recvc chan BackendMessage, sendc chan ConsensusMessage){

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
	if err := nh.StartOnDiskCluster(peers, *join, NewDiskKV, rc); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add cluster, %v\n", err)
		os.Exit(1)
	}
	raftStopper := syncutil.NewStopper()
	raftStopper.RunWorker(func() {
		// cs := nh.GetNoOPSession(exampleClusterID)
		for {
			select {
			case backend_msg, ok := <-recvc:
				if !ok {
					return
				}

				// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

				switch backend_msg.Type {
				case GET_IMAGE:
					fmt.Fprintf(os.Stderr, "GET_IMAGE Not Implemented\n", err)
				case UPDATE_PIXEL:
					fmt.Fprintf(os.Stderr, "UPDATE_PIXEL Not Implemented\n", err)
				case ADD_USER:
					fmt.Fprintf(os.Stderr, "ADD_USER Not Implemented\n", err)
				case GET_LAST_USER_UPDATE:
					fmt.Fprintf(os.Stderr, "GET_LAST_USER_UPDATE Not Implemented\n", err)
				}
				
				// rt, key, val, ok := parseCommand(msg)
				

				// if rt == PUT {
				// 	kv := &KVData{
				// 		Key: key,
				// 		Val: val,
				// 	}
				// 	data, err := json.Marshal(kv)
				// 	if err != nil {
				// 		panic(err)
				// 	}
				// 	_, err = nh.SyncPropose(ctx, cs, data)
				// 	if err != nil {
				// 		fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
				// 	}
				// } else {
				// 	result, err := nh.SyncRead(ctx, exampleClusterID, []byte(key))
				// 	if err != nil {
				// 		fmt.Fprintf(os.Stderr, "SyncRead returned error %v\n", err)
				// 	} else {
				// 		fmt.Fprintf(os.Stdout, "query key: %s, result: %s\n", key, result)
				// 	}
				// }
				// cancel()
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}