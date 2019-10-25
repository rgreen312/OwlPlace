# OwlPlace Backend

Documentation for the go service which handles OwlPlace API requests and
maintains consensus with other replicas.

## API Server

Go Package: [`github.com/rgreen312/owlplace/server/apiserver`](https://github.com/rgreen312/OwlPlace/tree/master/server/apiserver)

### Endpoints & Supported Requests

TODO

## Consensus Module

Go Package: [`github.com/rgreen312/owlplace/server/apiserver`](https://github.com/rgreen312/OwlPlace/tree/master/server/apiserver)

## Development Setup

1. Install Go: https://golang.org/doc/install#install
1. Install RocksDB: https://github.com/facebook/rocksdb/blob/master/INSTALL.md
1. Clone this repo **outside** of `$GOPATH`
1. Building: navigate to this folder and run `go build`!
1. Testing: navigate to this golder and run `go test`!


## Building with Docker

To start a cluster member, first define a cluster configuration file:

```json
[
   "1": {
       "hostname": "backend1",
       "api_port": 3000,
       "consensus_port": 63000
   }
   ...
]
```
where the keys of the JSON object are the node IDs of the cluster members.

## building

The Dockerfile (`Dockerfile.server`) leverages a multi-stage builder pattern that builds our deps, our service, and then copies our final binaries into a runtime image.  To build (from the main folder), run:

```shell
docker build . -f Dockerfile.server -t owlplace
```

For Windows: If Docker doesn't work, try installing [VirtualBox 5.2.6](https://download.virtualbox.org/virtualbox/5.2.6/), other versions may/may not work.

## running

At this point, we have an image that will run our service packaged with its dependencies.  To run several of them at one time, we'll use [docker compose](https://docs.docker.com/compose/) to define a few services.  See `docker-compose.yml` for an example, and run `docker-compose up` to start 3 services.

If localhost:3001 doesn't work, try using the docker IP (docker-machine ip in terminal) insteadl of localhost, like \[docker IP\]:3001

