# OwlPlace Backend

Documentation for the go service which handles OwlPlace API requests and
maintains consensus with other replicas.

## Developing with Docker

Ensure you have a working version of Docker and GNU Make installed.  The
Makefile provided stores some convenience directives for building and running
our service.

**Note for Windows: If Docker doesn't work, try installing [VirtualBox
5.2.6](https://download.virtualbox.org/virtualbox/5.2.6/), other versions
may/may not work.**

### building and packaging


#### building the dependency image

```
$ make builder
```

#### building the code

```
$ make build
```

#### building a runtime image

```
$ make image
```

### running a cluster

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
This is the `docker-owlplace.json` file.

To deploy 3 services as defined in `docker-compose.yml`, run:
```
$ docker-compose up
```

This should run 3 backend services that communicate together over their own
docker network.  Note that in our two config files (`docker-compose.yml` and
`docker-owlplace.json`) we set our API ports to `3000` and bind those container
ports to host ports `3001`, `3002`, and `3003`.  So, you should be able to make
API requests against our API at `localhost:3001`.  If you're using an older
version of Windows, it's possible Docker failed to correctly map these ports,
so you should replace `localhost` with your docker machine's IP address, which
you can find using `docker-machine ip` in the Docker Quickstart Terminal.
