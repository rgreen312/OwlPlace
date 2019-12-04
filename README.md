# owlplace

This repository contains an implementation of [Reddit
Place](https://www.reddit.com/r/place/) created by a team in Rice's COMP413:
Distributed Program Construction during the Fall of 2019.  

## building

The backend is written in Go, but has a external dependency on RocksDB.  If you
have RocksDB installed you should be able to simply build the backend using a
typical `go build`.  Alternatively, we have a make target available that builds a
docker image: `make image`. 

## running

There are two methods of running a set of owlplace servers together locally.

### docker compose

Build the service image.
```
make image
```

Start 3 owlplace services, see `docker-compose.yml` for details.  The container
API ports are exposed at host (or `docker-machine`) ports `3001-3003`.
```
docker-compose up
```

Send a trigger to one of the backend services.
```
curl localhost:3001/consensus_trigger
```

### kubernetes

**TODO**
