# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Image URL to use all building/pushing image targets
IMG := treeship:$(shell git rev-parse --short HEAD)-$(shell date +%s)

# KIND_CLUSTER_NAME refers to the name of the kind cluster to be used for development.
KIND_CLUSTER_NAME ?= treeship

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
KIND ?= kind

.PHONY: run
run: kind-check fmt vet ## Run a treeship from your host.
	CLUSTER_NAME="treeship-1" SERVER_URL="http://localhost:8080/api/v1/conn"  go run ./cmd/agent/main.go --kubeconfig ~/.kube/config -l debug --dev

##@ Safety checks
.PHONY: kind-check
kind-check: ## Check if the kind cluster is running.
	$(KIND) get clusters | grep -q $(KIND_CLUSTER_NAME) || (echo "Kind cluster $(KIND_CLUSTER_NAME) is not running" && exit 1)
## Check if the current context is the kind cluster.
	$(KUBECTL) config current-context | grep -q 'kind-$(KIND_CLUSTER_NAME)' || (echo "Current kubeconfig context is not the $(KIND_CLUSTER_NAME) Kind cluster" && exit 1)

##@ Local Development environment
.PHONY: dev
dev: up  ## Start a Kind cluster and deploy Flux and podinfo to it.
#	$(KUBECTL) wait --for=condition=available --timeout=60s deployment/helmrelease-watcher-controller-manager -n helmrelease-watcher-system

.PHONY: up
up:
	./dev/up.sh
.PHONY: kind-check dev-image-load
image-load: ## Load the docker image into the kind cluster.
	$(KIND) load docker-image ${IMG} --name ${KIND_CLUSTER_NAME}

##@ Tests
.PHONY: test
test: fmt vet
	go test -v ./...
	
.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...	

.PHONY: run_server
run_server: ## Run the Treeship server.
	go run ./cmd/server/main.go

.PHONY: run_agent
run_agent: ## Run the Treeship server.
	CLUSTER_NAME="treeship-1" SERVER_URL="http://localhost:8080/api/v1/conn"  go run ./cmd/agent/main.go --id agent-2 --kubeconfig ~/.kube/config -l debug

# Generate proto files
.PHONY: proto
proto:
	./scripts/generate-proto.sh