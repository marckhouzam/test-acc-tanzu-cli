BUILD_VERSION ?= $$(cat BUILD_VERSION)
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")

LD_FLAGS = -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildDate=$(BUILD_DATE)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildSHA=$(BUILD_SHA)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)'


build-local:
	@echo BUILD_VERSION: $(BUILD_VERSION)
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin --target local

build:
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin

test:
	go test ./...