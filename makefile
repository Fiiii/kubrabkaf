# ==============================================================================
# Variables
GOLANG          := golang:1.21.4
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.3
KIND_CLUSTER    := kubrabkaf-cluster
NAMESPACE       := kubrabkaf-infra
APP             := kubrabkaf
BASE_IMAGE_NAME := fiiii/kubrabkaf
SERVICE_NAME    := kubrabkaf-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)


# ==============================================================================
# Local development
docker-build:
	docker build \
		-f zarf/docker/dockerfile \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

docker-run: docker-build
	docker run -it \
	  -p 3000:3000 \
	  $(SERVICE_IMAGE)

# ==============================================================================
# Infra k8s dependencies

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize

# ==============================================================================
# Running k8s cluster

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	kubectl create -f https://download.elastic.co/downloads/eck/2.10.0/crds.yaml
	kubectl apply -f https://download.elastic.co/downloads/eck/2.10.0/operator.yaml


dev-load:
	cd zarf/k8s/dev/kubrabkaf/; kustomize edit set image service-image=$(SERVICE_IMAGE)
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/kafka | kubectl apply -f -
	kustomize build zarf/k8s/dev/elk/elasticsearch | kubectl apply -f -
	kustomize build zarf/k8s/dev/kubrabkaf | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6

make dev-update: docker-build dev-load dev-apply dev-restart

# ==============================================================================
# Core modules commands

dev-kafka-srv:
	kubectl get services -n kafka


# ==============================================================================
# Geth commands
geth-init:
	geth --datadir ./zarf/network/test-net-blockchain  init ./zarf/network/seed/genesis.json

geth-console:
	geth --networkid 15 console

geth-account:
	geth account new --datadir ./zarf/network/test-net-blockchain

geth-cleanup:
	geth removedb --datadir ./zarf/network/test-net-blockchain

geth-run:
	geth --identity "MyTestNetNode"

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache