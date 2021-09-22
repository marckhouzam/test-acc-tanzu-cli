BUILD_VERSION ?= $$(cat BUILD_VERSION)
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")

LD_FLAGS = -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildDate=$(BUILD_DATE)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildSHA=$(BUILD_SHA)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)'

GO_SOURCES = $(shell find ./cmd ./pkg -type f -name '*.go')



build-local:
	@echo BUILD_VERSION: $(BUILD_VERSION)
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin --target local

build:
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin

test:
	go test -coverprofile cover.out ./...

create-artifact: build
	tar -zcvf tanzu-accelerator-plugin.tar.gz artifacts

docs: $(GO_SOURCES)
	@rm -rf docs
	go run --ldflags "$(LD_FLAGS)" ./cmd/plugin/accelerator docs -d docs

create-kind-cluster:
	kind delete clusters e2e-acc-cluster
	kind create cluster --name e2e-acc-cluster --config ./e2e/assets/acc-kind-node.yml
	kubectl config use-context kind-e2e-acc-cluster

install-flux2:
	kubectl apply -f https://gist.githubusercontent.com/trisberg/f53bbaa0b8aacba0ec64372a6fb6acdf/raw/45259afd682caa2f6270f4b8c07c995aa8487a12/acc-flux2.yaml	

install-bundle:
	kubectl create namespace accelerator-system	
	./e2e/scripts/deploy-app.sh

add-test-accelerators:
	kubectl create -f ./e2e/assets/test-accelerators.yml

create-context: create-kind-cluster install-flux2 install-bundle add-test-accelerators	