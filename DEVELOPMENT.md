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


## Updating Contrail operator
    operator-sdk build contrail-operator
    # optionally: docker tag and docker push

## Building contrail-provisioner

### Generate contrail-api-client
First, install Python2 and the following Python libraries:

    pip install future lxml

Then, generate `contrail-api-client` files:

    make generate

`make generate` will also change go.mod file, so that locally generated `contrail-api-client` can be imported. This needs to be run only once.

### Build contrail-provisioner

    make provisioner

### Push contrail-provisioner to local registry
Assuming you have kind-registry running on port 5000 (on how to setup this, see
[E2E test guide](test/env/README.md)):

    docker tag contrail-provisioner:latest localhost:5000/contrail-provisioner:latest
    docker push localhost:5000/contrail-provisioner:latest


### Change contrail-provisioner image that is used in Contrail cluster
If you add tag other than `latest` to contrail-provisioner, image that is deployed in kind cluster needs to be changed. Edit file `test/env/deploy/cluster.yaml`, find `provisionManager` and change `provisioner` image. Similarly, you can also change configuration of other services in file `test/env/deploy/cluster.yaml`.

    provisionManager:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: provmanager1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""
          replicas: 1
        serviceConfiguration:
          containers:
            init:
              image: registry:5000/python:3.8.2-alpine
            provisioner:
              image: registry:5000/contrail-provisioner:latest

After this change apply changes to k8s cluster:

    cd test/env
    ./apply_contrail_cluster.sh

## Measuring Unit Tests code coverage

### Measure code coverage by package

    go test -coverprofile=cov.out ./pkg/...

### Calculate total code coverage 

    go tool cover -func cov.out | grep total | awk '{print $3}'

### Display code coverage in the browser

    go tool cover -html=cov.out
