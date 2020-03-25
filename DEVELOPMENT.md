# Development

This file contains more detailed description of the development process. For quickstart guide see: [MAC OS quickstart](QUICKSTART-MACOS.md) or
[Linux quickstart](QUICKSTART-LINUX.md). For description of E2E test environment see [E2E test guide](test/env/README.md).


Instructions below assume that you are running operator-sdk in docker. If it is already installed locally, then the following steps can be omitted:
```
$ docker run --rm -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator
...
$ exit
```

## Repository structure overview
During development the following files are most often edited:  
`pkg/apis/contrail/v1alpha1/*` contains definitions of custom resources written in Go (based on these files yaml files are generates, see section Generate k8s files).  
`pkg/controller/*` contains code of controllers  
`test/*` contains code of e2e tests  
`test/env` contains scripts for deploying kind cluster  
`contrail-provisioner/*` contains contrail-provisioner source code.



## Add new API and controller
Replace Memcached with the new resource name.

```
$ cd github.com/Juniper/contrail-operator
$ docker run --rm -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator
$ operator-sdk add api --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached
$ operator-sdk add controller --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached
$ exit

$ sudo chown -R `id -u`:`id -g` ./**/*
```


## Generate k8s files
After custom resource specification (in `pkg/apis/contrail/v1alpha1/*_types.go`) has changed, code needs to be re-generated.

```
$ cd github.com/Juniper/contrail-operator
$ docker run --rm -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator
$ operator-sdk generate k8s
$ operator-sdk generate openapi
$ exit

$ cd deploy
$ ./create_manifest.sh
```

## Troubleshooting

* Problem: unsupported type invalid type for invalid type
  Solution: export GOROOT
* Problem: on running operator container `/usr/local/bin/entrypoint: Permission denied`
  Solution:
  ```
  sudo chown -R `id -u`:`id -g` build
  chmod -R 755 build/bin build/_output
  <rebuild operator>
  ```


## Updating Contrail operator
```
$ cd github.com/Juniper/contrail-operator
$ docker run --rm -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator
$ operator-sdk build contrail-operator
$ exit
```

## Building contrail-provisioner

### Change go.mod file
In go.mod file uncomment the following line:
    
    github.com/Juniper/contrail-go-api => ./build/contrail-go-api
    
In repository this line is commented out, because contrail-operator CI does not perform 'make generate' and building operator would fail with unresolved import `build/contrail-go-api`


### Install additional dependancies
In order to generate `contrail-go-api`, which is used by contrail-provisioner, Python2 is required and the following Python libraries need to be installed:

    pip install future lxml

### Generate contrail-go-api

    make generate

### Build contrail-provisioner

    make provisioner

### Push contrail-provisioner to local registry
Assuming you have kind-registry running on port 5000 (on how to setup this, see
[E2E test guide](test/env/README.md)):
```
docker tag contrail-provisioner:latest localhost:5000/contrail-provisioner:latest
docker push localhost:5000/contrail-provisioner:latest
```
  
### Change contrail-provisioner image that is used in Contrail cluster
If you add tag othen than `latest` to contrail-provisioner, image that is deployed in kind cluster needs to be changed. Edit file `test/env/deploy/cluster.yaml`, find `provisionManager:` and change `provisioner` image. Similarly, you can also change configuration of other services in file `test/env/deploy/cluster.yaml`.
```
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
          image: registry:5000/python:alpine
        provisioner:
          image: registry:5000/contrail-provisioner:latest
```
After this change apply changes to k8s cluster:
```
cd test/env
./apply_contrail_cluster.sh
```