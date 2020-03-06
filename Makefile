BUILD_DIR := build
GENERATE_DS_REPO := contrail-api-client
GO_API_CLIENT_REPO := contrail-go-api
GO_API_CLIENT_VENDOR := ./vendor/github.com/Juniper/$(GO_API_CLIENT_REPO)
GENERATE_DS_REPO_DIR ?= ""
GENERATE_DS_BRANCH ?= master
GENERATE_DS_REVISION ?= HEAD
GO_API_CLIENT_REPO_DIR ?= ""
GO_API_CLIENT_BRANCH ?= master
GO_API_CLIENT_REVISION ?= HEAD
GOPATH ?= `go env GOPATH`
SOURCEDIR ?= $(GOPATH)
all: vendor generate
vendor: ## Ensure vendor dependencies
	go mod vendor
generate: ## Generate go api client types
	rm -rf $(BUILD_DIR)/$(GENERATE_DS_REPO) $(BUILD_DIR)/$(GO_API_CLIENT_REPO)
ifeq ($(GENERATE_DS_REPO_DIR),"")
	git clone -b $(GENERATE_DS_BRANCH) https://github.com/michaelhenkel/$(GO_API_CLIENT_REPO).git $(BUILD_DIR)/$(GO_API_CLIENT_REPO)
	git clone -b $(GENERATE_DS_BRANCH) https://github.com/Juniper/$(GENERATE_DS_REPO).git $(BUILD_DIR)/$(GENERATE_DS_REPO)
	cd $(BUILD_DIR)/$(GENERATE_DS_REPO) && git checkout $(GENERATE_DS_REVISION) && git fetch https://github.com/Juniper/$(GENERATE_DS_REPO).git refs/changes/19/56219/1 && git cherry-pick FETCH_HEAD
	cd  $(BUILD_DIR)/$(GO_API_CLIENT_REPO) && go mod init
else
	cp -r $(GENERATE_DS_REPO_DIR) $(BUILD_DIR)/$(GENERATE_DS_REPO)
endif
	$(BUILD_DIR)/$(GENERATE_DS_REPO)/generateds/generateDS.py -f -o $(BUILD_DIR)/$(GO_API_CLIENT_REPO)/types -g golang-api $(BUILD_DIR)/$(GENERATE_DS_REPO)/schema/all_cfg.xsd
monitor: ## make statusmonitor
	cd statusmonitor && go build -o statusmonitor
provisioner: ## make provisioner
	cd contrail-provisioner && go build -o contrail-provisioner
