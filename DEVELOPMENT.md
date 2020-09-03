# Development

This file contains more detailed description of the development process. For quickstart guide see [README.md](README.md). For description of E2E test environment see [E2E test guide](test/env/README.md).

## Repository structure overview
During development the following files are most often edited:  
* `pkg/apis/contrail/v1alpha1/*` - definitions of custom resources written in Go (based on these files yaml files are generated, see section  [Generate k8s files](#generate-k8s-files))
* `pkg/controller/*` - code of controllers
* `test/*` - code of e2e tests
* `test/env` - scripts for deploying kind cluster
* `contrail-provisioner/*` - contrail-provisioner source code



## Add new API and controller
Replace Memcached with the new resource name.

    operator-sdk add api --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached
    operator-sdk add controller --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached


## Generate k8s files
After custom resource specification (in `pkg/apis/contrail/v1alpha1/*_types.go`) is changed, code needs to be re-generated.

    operator-sdk generate k8s
    operator-sdk generate crds
    cd deploy
    ./create_manifest.sh

The last line `./create_manifest.sh` generates file `deploy/1-create-operator.yaml` which contains all custom resource definitions and `contrail-operator` deployment specification. You can also edit this file by hand, for example to change image of `contrail-operator` that is used (by default it is `registry:5000/contrail-operator:latest`).

## Troubleshooting

* Problem: unsupported type invalid type for invalid type
  Solution: export GOROOT
* Problem: on running operator container `/usr/local/bin/entrypoint: Permission denied`
  Solution:

      sudo chown -R `id -u`:`id -g` ./**/*
      chmod -R 755 build/bin build/_output
      <rebuild operator>

### Debug Contrail operator in VSCode

Contrail operator can be run outside of K8s/OpenShift cluster what makes it easier to debug. First of all the configuration of the cluster needs to be available on the host where operator will be run. This can be achieved by creating local K8s cluster with script `test/env/create_testenv.sh` or by putting configuration of external cluster to standard directory for example `~/.kube/config`.

Next the [Delve](https://github.com/go-delve/delve) go debugger needs to be installed and available in `PATH`. Please check [installation notes](https://github.com/go-delve/delve/tree/master/Documentation/installation).

Prepare env for start of local operator by running `test/env/debug_operator.sh` script. Then go to VSCode set some breakpoints in source files by clicking on red dots before line numbers. Switch to `Run` view and on the top left choose `Run operator locally` and hit `Start Debugging`. Operator's log will appear in `DEBUG CONSOLE` and panels on the left will give insight to variables values, call stack etc. at defined breakpoints.

To run operator locally with debugger enabled outside of VSCode this command can be used: `operator-sdk run --local --watch-namespace contrail --enable-delve`.

## Updating Contrail operator
    operator-sdk build contrail-operator
    # optionally: docker tag and docker push

## Building contrail-operator-provisioner and contrail-statusmonitor

Containers contrail-operator-provisioner and contrail-statusmonitor are services which have source code in this repository (directories contrail-provisioner and statusmonitor). The only officially supported way to build them now is with bazel.

### Install bazel
On Mac:

    brew install bazel

On Linux:

    wget https://github.com/bazelbuild/bazel/releases/download/0.29.1/bazel-0.29.1-installer-linux-x86_64.sh
    sudo ./bazel-0.29.1-installer-linux-x86_64.sh


### Build containers and push them to local registry
In order to change parameters, for example registry to push to, edit file:

    contrail-provisioner/BUILD.bazel, rule contrail-provisioner-push-local
    statusmonitor/BUILD.bazel, rule contrail-statusmonitor-push-local

Make sure that the registry (by default localhost:5000) is up and you have write access to it and then run:

    bazel run //contrail-provisioner:contrail-provisioner-push-local
    bazel run //statusmonitor:contrail-statusmonitor-push-local


### Change images that are used in Contrail cluster
Edit file `test/env/deploy/cluster.yaml`, find proper images and change them.  
  
For example for contrail-operator-provisioner - find `provisionManager` section and change `provisioner` image (by default it is `registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:master.latest`).  

    provisionManager:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: provmanager1
      spec:
        commonConfiguration:
          nodeSelector:
            node-role.kubernetes.io/master: ""
          replicas: 1
        serviceConfiguration:
          keystoneInstance: keystone
          globalVrouterConfiguration:
            ecmpHashingIncludeFields:
              destinationIp: true
              destinationPort: true
              hashingConfigured: true
              ipProtocol: true
              sourceIp: true
              sourcePort: true
            encapPriority: VXLAN,MPLSoGRE,MPLSoUDP
            vxlanNetworkIdentifierMode: automatic
          containers:
          - name: init
            image: registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine
          - name: provisioner
            image: registry:5000/contrail-operator-provisioner:latest

In order to change contrail-statusmonitor - change `registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:master.latest`. Similarly, you can also change configuration of other services in file `test/env/deploy/cluster.yaml`.  

After making changes apply changes to k8s cluster:

    cd test/env
    ./apply_contrail_cluster.sh

## Measuring Unit Tests code coverage

### Measure code coverage by package

    go test -coverprofile=cov.out ./pkg/...

### Calculate total code coverage 

    go tool cover -func cov.out | awk '/total/{print $3}'

### Display code coverage in the browser

    go tool cover -html=cov.out
