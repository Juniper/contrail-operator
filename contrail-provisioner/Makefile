BUILD_DIR := ../build
GENERATE_DS_REPO := contrail-api-client
GO_API_CLIENT_REPO := contrail-go-api
GO_API_CLIENT_VENDOR := ./vendor/github.com/Juniper/$(GO_API_CLIENT_REPO)
GENERATE_DS_REPO_DIR ?= ""
GENERATE_DS_BRANCH ?= master
GENERATE_DS_REVISION ?= HEAD
GO_API_CLIENT_REPO_DIR ?= ""
GO_API_CLIENT_BRANCH ?= master
GO_API_CLIENT_REVISION ?= HEAD
BASE_IMAGE_REGISTRY ?= opencontrailnightly
BASE_IMAGE_REPOSITORY ?= contrail-base
BASE_IMAGE_TAG ?= latest
GOPATH ?= `go env GOPATH`
SOURCEDIR ?= $(GOPATH)
DOCKER_FILE := $(BUILD_DIR)/docker/Dockerfile
all: vendor generate install test lint ## Perform all checks
vendor: go mod vendor
generate:
	rm -rf $(BUILD_DIR)/$(GENERATE_DS_REPO)
ifeq ($(GENERATE_DS_REPO_DIR),"")
	git clone -b $(GENERATE_DS_BRANCH) https://github.com/Juniper/$(GENERATE_DS_REPO).git $(BUILD_DIR)/$(GENERATE_DS_REPO)
	cd $(BUILD_DIR)/$(GENERATE_DS_REPO) && git checkout $(GENERATE_DS_REVISION)
else
	cp -r $(GENERATE_DS_REPO_DIR) $(BUILD_DIR)/$(GENERATE_DS_REPO)
endif
	$(BUILD_DIR)/$(GENERATE_DS_REPO)/generateds/generateDS.py -f -o $(GO_API_CLIENT_VENDOR)/types -g golang-api $(BUILD_DIR)/$(GENERATE_DS_REPO)/schema/all_cfg.xsd
install:
	go install ./cmd
test:
lint:
fmt:
	gofiles=$(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*gen*")
	goimports -v -w ${gofiles}
	gofmt -s -w ${gofiles}
