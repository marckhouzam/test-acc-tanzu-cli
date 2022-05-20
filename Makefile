BUILD_VERSION ?= $$(cat BUILD_VERSION)
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")

LD_FLAGS = -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildDate=$(BUILD_DATE)' \
           -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildSHA=$(BUILD_SHA)$(BUILD_DIRTY)' \
           -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)'

GO_SOURCES = $(shell find ./cmd ./pkg -type f -name '*.go')

build-local:
	@echo BUILD_VERSION: $(BUILD_VERSION)
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin --target local

build:
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin

.PHONY: build-%
build-%:
	$(eval ARCH = $(word 2,$(subst -, ,$*)))
	$(eval OS = $(word 1,$(subst -, ,$*)))
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin --artifacts artifacts/${OS}/${ARCH}/cli --target ${OS}_${ARCH}

.PHONY: publish-%
publish-%:
	$(eval ARCH = $(word 2,$(subst -, ,$*)))
	$(eval OS = $(word 1,$(subst -, ,$*)))
	tanzu builder publish --type local --plugins "accelerator" --version $(BUILD_VERSION) --local-output-discovery-dir standalone/${OS}-${ARCH}/discovery/standalone --local-output-distribution-dir standalone/${OS}-${ARCH}/distribution --input-artifact-dir artifacts --os-arch "${OS}-${ARCH}"

test:
	go test -coverprofile cover.out ./...

create-artifact: build-darwin-amd64 build-linux-amd64 build-windows-amd64 publish-darwin-amd64 publish-linux-amd64 publish-windows-amd64
	tar -zcvf tanzu-accelerator-plugin.tar.gz standalone

docs: $(GO_SOURCES)
	@rm -rf docs
	ACC_SERVER_URL="" go run --ldflags "$(LD_FLAGS)" ./cmd/plugin/accelerator docs -d docs

create-kind-cluster:
	kind delete clusters e2e-acc-cluster
	kind create cluster --name e2e-acc-cluster --config ./e2e/assets/acc-kind-node.yml
	kubectl config use-context kind-e2e-acc-cluster

install-prereqs:
	kubectl apply -f https://gist.githubusercontent.com/trisberg/f53bbaa0b8aacba0ec64372a6fb6acdf/raw/45259afd682caa2f6270f4b8c07c995aa8487a12/acc-flux2.yaml	
	kubectl apply -f https://gist.githubusercontent.com/trisberg/0fc3ed74e3673af72c63e867ffcf8972/raw/9d5834bd81fa407665c00a238d6dd0ade135ca10/acc-source.yaml

install-bundle:
	kubectl create namespace accelerator-system	
	./e2e/scripts/deploy-app.sh

add-test-accelerators:
	kubectl create -f ./e2e/assets/test-accelerators.yml

create-context: create-kind-cluster install-prereqs install-bundle add-test-accelerators	