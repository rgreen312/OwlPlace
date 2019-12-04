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

### minikube

You can test Kubernetes locally using
[minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

1.) Install minikube
2.) Start minikube
3.) Run the following commands

**Create Container**
```
docker build . -f Dockerfile.server -t <username>/owlplace
docker push <username>/owlplace
```

**Initial Configuration**
```
kubectl create -f kubernetes/namespaces/dev.yaml
kubectl config set-context dev --namespace=dev --cluster=minikube --user=minikube
kubectl config use-context dev
kubectl create -f kubernetes/roles/pods-reader.yaml
kubectl create clusterrolebinding service-reader-pod-dev-2 --clusterrole=service-reader-dev-2 --serviceaccount=dev:default
```

**Run the Deployment**
```
kubectl create -f kubernetes/deployments/owlplace-backend.yaml
kubectl expose deployment owlplace-backend --type=LoadBalancer
kubectl get services
```

The last command will give you output that looks similar to the following
```
NAME               TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
owlplace-backend   LoadBalancer   10.99.205.94   <pending>     3001:32596/TCP   16m
```
The only thing that matters here is the `32596` because that is the local port the forwards to our application.

**Get your minikube ip address**
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

## deploying to k8s

TODO: add information on deploying to the production cluster
