DOCKERFILE = $(PWD)/Dockerfile.server
OWLPLACE_CACHE = $(PWD)/.go-cache
BUILD_IMAGE = owlplace-builder
COMPILE_IMAGE = owlplace-fat
RUNTIME_IMAGE = owlplace
PROJECT_DIR = $(PWD)/server

all: image

builder: $(DOCKERFILE)
	docker build . -f $(DOCKERFILE) --target builder -t $(BUILD_IMAGE)

build:
	docker run --rm -it -v $(PROJECT_DIR):/app -v $(OWLPLACE_CACHE):/go/pkg/mod --workdir /app $(BUILD_IMAGE) go build

docker-build:
	docker build . -f $(DOCKERFILE) --target compile-image -t $(COMPILE_IMAGE)

image: $(DOCKERFILE)
	docker build . -f $(DOCKERFILE) --target runtime-image -t $(RUNTIME_IMAGE)

.PHONY: builder image build docker-build
