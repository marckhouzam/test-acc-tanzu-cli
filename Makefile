ROOT_DIR_RELATIVE := .

include $(ROOT_DIR_RELATIVE)/common.mk
include $(ROOT_DIR_RELATIVE)/plugin-tooling.mk

TOOLS_DIR := tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
GOLANGCI_LINT_VERSION := 1.49.0

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the plugin
	$(GOLANGCI_LINT) run -v

.PHONY: gomod
gomod: ## Update go module dependencies
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile cover.out ./...

$(TOOLS_BIN_DIR):
	-mkdir -p $@

$(GOLANGCI_LINT): $(TOOLS_BIN_DIR) ## Install golangci-lint
	curl -L https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(GOHOSTOS)-$(GOHOSTARCH).tar.gz | tar -xz -C /tmp/
	mv /tmp/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(GOHOSTOS)-$(GOHOSTARCH)/golangci-lint $(@)

.PHONY: create-kind-cluster
create-kind-cluster:
	kind delete clusters e2e-acc-cluster
	kind create cluster --name e2e-acc-cluster --config ./e2e/assets/acc-kind-node.yml
	kubectl config use-context kind-e2e-acc-cluster

.PHONY: install-prereqs
install-prereqs:
	kubectl apply -f ./e2e/assets/fluxcd-flux2-install.yaml
	kubectl apply -f https://gist.githubusercontent.com/trisberg/0fc3ed74e3673af72c63e867ffcf8972/raw/9d5834bd81fa407665c00a238d6dd0ade135ca10/acc-source.yaml

.PHONY: install-bundle
install-bundle:
	kubectl create namespace accelerator-system	
	./e2e/scripts/deploy-app.sh

.PHONY: add-test-accelerators
add-test-accelerators:
	kubectl create -f ./e2e/assets/test-accelerators.yml
	sleep 10

.PHONY: create-context
create-context: create-kind-cluster install-prereqs install-bundle add-test-accelerators
