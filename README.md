# owlplace

This repository contains an implementation of [Reddit
Place](https://www.reddit.com/r/place/) created by a team in Rice's COMP413:
Distributed Program Construction during the Fall of 2019.  

## building

The backend is written in Go, but has a external dependency on RocksDB.  If you
have RocksDB installed you should be able to simply build the backend using a
typical `go build`.  Alternatively, we have a make target available that builds a
docker image: `make image`. 


pushing and pulling image:
https://cloud.google.com/container-registry/docs/pushing-and-pulling

## running

There are a few methods of running a set of owlplace servers together.

### locally with docker compose

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

### locally with minikube

You can test Kubernetes locally using
[minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

1.) Install minikube
2.) Start minikube
3.) Run the following commands

**Create Container**

```
make image
```

**Initial Configuration**
```
make minikube-setup
```

**Run the Deployment**
```
make minikube-deploy
kubectl get services
```

The last command will give you output that looks similar to the following
```
NAME               TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
owlplace-backend   LoadBalancer   10.99.205.94   <pending>     3001:32596/TCP   16m
```

The only thing that matters here is the `32596` because that is the local port
the forwards to our application.

**minikube IP address**
```
minikube ip
```

You will see the external IP address is pending for a while.

**Start the cluster**

```
curl -k <minikube ip>:<forward port>/consensus_trigger
```

**Destroy the deployment**
```
kubectl delete service,deployment owlplace-backend
```

## deploying to a test cluster on gke

Since our application is packaged with kubernetes, we can create a test cluster
using GKE as well.

```
make gke-create
```

```
make gke-setup
```

```
make gke-deploy
```
