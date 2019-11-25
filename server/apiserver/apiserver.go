package apiserver

import (
	
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
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
	config   *common.ServerConfig
	upgrader *websocket.Upgrader
	pool     *wsutil.Pool
	cons     consensus.IConsensus
}

func NewApiServer(servers map[int]*common.ServerConfig, nodeId int) (*ApiServer, error) {

	conf, ok := servers[nodeId]
	if !ok {
		return nil, errors.Wrapf(ConfigurationError, "missing entry for node: %d", nodeId)
	}

	// First we create the pool because we're going to share it's broadcast
	// channel with the consensus service.
	pool := wsutil.NewPool()

	cons, err := consensus.NewConsensusService(servers, nodeId, pool.Broadcast)
	if err != nil {
		return nil, errors.Wrap(err, "creating ConsensusService")
	}

	err = cons.Start()
	if err != nil {
		return nil, errors.Wrap(err, "starting ConsensusService")
	}

	log.WithFields(log.Fields{
		"server config":     conf,
		"consensus service": cons,
	}).Debug()

	return &ApiServer{
		config: conf,
		pool:   pool,
		cons:   cons,
	}, nil
}

func (api *ApiServer) ListenAndServe() {
	go api.pool.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsutil.ServeWs(api.pool, api.cons, w, r)
	})

	// Although there is nothing wrong with this line, it prevents us from
	// running multiple nodes on a single machine.  Therefore, I am making
	// failure non-fatal until we have some way of running locally from the
	// same port (i.e. docker)
	// log.Fatal(http.ListenAndServe(":3010", nil))
	http.ListenAndServe(fmt.Sprintf(":%d", api.config.ApiPort), nil)
}

