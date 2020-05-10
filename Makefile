ROOT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR := $(ROOT_DIR)/build
GENERATE_DS_REPO := contrail-api-client
GO_API_CLIENT_REPO := contrail-go-api
GENERATE_DS_REPO_DIR ?= ""
GENERATE_DS_BRANCH ?= master
GENERATE_DS_REVISION ?= d964cff846242468ff9ba968d73a42a60b7058cd # `Merge "generateDS.py broken on MacOS"`

all: generate provisioner statusmonitor

generate: ## Generate go api client types
	rm -rf $(BUILD_DIR)/$(GENERATE_DS_REPO) $(BUILD_DIR)/$(GO_API_CLIENT_REPO)
ifeq ($(GENERATE_DS_REPO_DIR),"")
	git clone -b $(GENERATE_DS_BRANCH) https://github.com/Juniper/$(GO_API_CLIENT_REPO).git $(BUILD_DIR)/$(GO_API_CLIENT_REPO)
	git clone -b $(GENERATE_DS_BRANCH) https://github.com/Juniper/$(GENERATE_DS_REPO).git $(BUILD_DIR)/$(GENERATE_DS_REPO)
	cd $(BUILD_DIR)/$(GENERATE_DS_REPO) && git checkout $(GENERATE_DS_REVISION)
	cd $(BUILD_DIR)/$(GO_API_CLIENT_REPO) && go mod init github.com/contrail-operator/build/contrail-go-api
	go mod edit -replace github.com/Juniper/contrail-go-api=./build/$(GO_API_CLIENT_REPO)
else
	cp -r $(GENERATE_DS_REPO_DIR) $(BUILD_DIR)/$(GENERATE_DS_REPO)
endif
	$(BUILD_DIR)/$(GENERATE_DS_REPO)/generateds/generateDS.py -f -o $(BUILD_DIR)/$(GO_API_CLIENT_REPO)/types -g golang-api $(BUILD_DIR)/$(GENERATE_DS_REPO)/schema/all_cfg.xsd

statusmonitor:
	cd $(ROOT_DIR)/statusmonitor && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o statusmonitor && docker build . -f Dockerfile.debug -t contrail-statusmonitor:latest

provisioner:
	cd $(ROOT_DIR)/contrail-provisioner && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o contrail-provisioner && docker build . -t contrail-provisioner:master.1175

.PHONY: statusmonitor provisioner generate
