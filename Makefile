DOCKERFILE = $(PWD)/Dockerfile.server
OWLPLACE_CACHE = $(PWD)/.go-cache
BUILD_IMAGE = owlplace-builder
COMPILE_IMAGE = owlplace-fat
RUNTIME_IMAGE = owlplace
PROJECT_DIR = $(PWD)/server
GCP_PROJECT_ID = owlplace
GKE_DEV_CLUSTER = owlplace-dev
GCR_IMAGE = gcr.io/$(GCP_PROJECT_ID)/$(RUNTIME_IMAGE)

all: image

builder: $(DOCKERFILE)
	docker build . -f $(DOCKERFILE) --target builder -t $(BUILD_IMAGE)

build:
	docker run --rm -it -v $(PROJECT_DIR):/app -v $(OWLPLACE_CACHE):/go/pkg/mod --workdir /app $(BUILD_IMAGE) go build

docker-build:
	docker build . -f $(DOCKERFILE) --target compile-image -t $(COMPILE_IMAGE)

image: $(DOCKERFILE)
	docker build . -f $(DOCKERFILE) --target runtime-image -t $(RUNTIME_IMAGE)
	docker tag $(RUNTIME_IMAGE) $(GCR_IMAGE)

lint:
	goimports -d $(PROJECT_DIR)

fix:
	goimports -w $(PROJECT_DIR)

gcp-auth:
	gcloud auth configure-docker

minikube-setup:
	kubectl config set-context dev --namespace=dev --cluster=minikube --user=minikube
	kubectl config use-context dev
	kubectl apply -f kubernetes/namespaces/dev.yaml
	kubectl apply -f kubernetes/roles/pods-reader.yaml

minikube-deploy:
	kubectl apply -f kubernetes/deployments/owlplace-backend.yaml
	kubectl expose deployment owlplace-backend --type=LoadBalancer

gke-create:
	gcloud container clusters create $(GKE_DEV_CLUSTER) \
		--enable-cloud-monitoring \
		--machine-type n1-standard-1 \
		--zone us-central1-a

gke-setup:
	kubectl config use-context "gke_$(GCP_PROJECT_ID)_us-central1-a_$(GKE_DEV_CLUSTER)"

gke-deploy:
	kubectl apply -f kubernetes/namespaces/dev.yaml
	kubectl apply -f kubernetes/roles/pods-reader.yaml
	kubectl apply -f kubernetes/deployments/owlplace-backend-dev.yaml
	kubectl expose deployment owlplace-backend --type=LoadBalancer

gke-destroy:
	gcloud container clusters delete $(GKE_DEV_CLUSTER)

gcr-push:
	docker push $(GCR_IMAGE)

.PHONY: builder image build docker-build lint fix
