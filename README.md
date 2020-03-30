# Contrail Operator

## References
[E2E test guide](test/env/README.md)  
[Detailed development guide](DEVELOPMENT.md)  

## Requirements
  * Go 1.13
  * Docker
  * Kubernetes client
  * operator-sdk (https://github.com/operator-framework/operator-sdk/)
  * Kubernetes cluster (only one node is supported right now)


# Contrail-Operator Development Quick Start

## Install Go 1.13

* https://golang.org/doc/install#install

## Checkout contrail-operator source code

Contrail-Operator is a Go Module therefore can be downloaded to a folder outside the GOPATH.

    git clone git@github.com:Juniper/contrail-operator.git

## Verify if contrail-operator can be built

    go build cmd/manager/main.go

## Install Kubernetes Client

On Mac OS: https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-macos  
On Linux: https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-linux

## Install IDE

We use Goland and Visual Studio Code. Install your favourite one.

## Install Kind

Kind is used as a lightweight Kubernetes cluster for development purposes

    GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0

Verify if it works (Mac OS):
    
    $ kind version
    kind v0.7.0 go1.13.8 darwin/amd64

Verify if it works (Linux):
    
    $ kind version
    kind v0.7.0 go1.13.8 linux/amd64

If command is not found, then reload `~/.zshrc` (on Mac OS) or `~/.bashrc` (on Linux) and verify if `~/go/bin` is in `$PATH`.

## Install Docker for Desktop (Mac OS only)

* https://hub.docker.com/editions/community/docker-ce-desktop-mac

### Increase memory amount in settings to 8GB:

- click Docker icon
- select Preferences
- go to Resources/Advanced
- increase memory to 8GB
- restart Docker for Desktop

## Install Docker Engine (Linux only)

Instruction for Ubuntu (other distros are available as well): https://docs.docker.com/install/linux/docker-ce/ubuntu/ 

## Log into contrail-nightly docker registry

This is needed for downloading Contrail Command Docker image. 

    docker login hub.juniper.net/contrail-nightly

## Create Kind test environment

Following commands will create Kubernetes cluster.

It also starts Docker registry on port 5000. All pods deployed in the cluster will pull images from this Docker Registry. 

    cd test/env
    ./create_testenv.sh

Verify if it works:

    $ kind get clusters
    kind

## Pull images to locker Docker registry

    cd test/env
    ./update_local_registry.sh

In case when you receive timeouts disable VPN and retry.

## Install operator-sdk

Operator-SDK is a set of tools for developing Kubernates Operators. It is needed for:

- Go code generation
- K8s Custom Resource Definitions generation
- building contrail-operator image
- running e2e tests (aka system tests)

### Operator-sdk installation on Mac OS

    $ curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v0.13.0/operator-sdk-v0.13.0-x86_64-apple-darwin
    $ chmod u+x ./operator-sdk-v0.13.0-x86_64-apple-darwin
    $ mv ./operator-sdk-v0.13.0-x86_64-apple-darwin /usr/local/bin/operator-sdk

Verify if it works:

    $ operator-sdk version
    operator-sdk version: "v0.13.0", commit: "1af9c95bb51420c55a7f7f2b7fabebda24451276", go version: "go1.13.3 darwin/amd64"

### Operator-sdk installation on Linux

    $ curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v0.17.0/operator-sdk-v0.17.0-x86_64-linux-gnu
    $ chmod u+x ./operator-sdk-v0.13.0-x86_64-linux-gnu  
    $ mv ./operator-sdk-v0.13.0-x86_64-linux-gnu /usr/local/bin/operator-sdk


Verify if it works:

    $ operator-sdk version
    operator-sdk version: "v0.13.0", commit: "1af9c95bb51420c55a7f7f2b7fabebda24451276", go version: "go1.13.3 linux/amd64"


## Build Contrail-Operator

In order to run Contrail-Operator in the Kubernetes cluster we have to build Docker Image.

    operator-sdk build localhost:5000/contrail-operator:latest

Verify:

    $ docker images | grep contrail-operator
    contrail-operator   latest   5c0148fdb7e8   4 seconds ago   125MB

After image is created we have to push it into local Docker registry.

    docker push localhost:5000/contrail-operator:latest

## Run Contrail-Operator with sample Contrail configuration

Following command will deploy Contrail-Operator on a working Kubernetes cluster. It will also create a sample Contrail configuration. Note: you can change this configuration by editing `test/env/deploy/cluster.yaml` file.

    cd test/env
    ./apply_contrail_cluster.sh

As soon as contrail-operator is started, it deploys Contrail services. It is a time-consuming process. You can watch the progress using following command:

    watch kubectl get pods -n contrail

Eventually all pods should be Running:

    NAME                                          READY   STATUS      RESTARTS   AGE
    cassandra1-cassandra-statefulset-0            1/1     Running     0          8m15s
    command-command-deployment-77644668cf-dpp6f   1/1     Running     0          7m21s
    config1-config-statefulset-0                  9/9     Running     0          4m47s
    contrail-operator-585f5bd8b5-hfdrz            1/1     Running     0          9m24s
    control1-control-statefulset-0                4/4     Running     0          2m56s
    keystone-keystone-statefulset-0               3/3     Running     0          7m8s
    memcached-deployment-5f5f974bd9-gthzx         1/1     Running     0          8m15s
    postgres-pod                                  1/1     Running     0          8m16s
    provmanager1-provisionmanager-statefulset-0   1/1     Running     0          2m57s
    rabbitmq1-rabbitmq-statefulset-0              1/1     Running     0          8m15s
    swift-proxy-deployment-754f87448b-6l5nc       1/1     Running     0          4m32s
    swift-ring-account-job-rnsxs                  0/1     Completed   0          7s
    swift-ring-container-job-pkb2k                0/1     Completed   0          7s
    swift-ring-object-job-7nn44                   0/1     Completed   0          7s
    swift-storage-statefulset-0                   13/13   Running     0          8m11s
    webui1-webui-statefulset-0                    3/3     Running     0          2m56s
    zookeeper1-zookeeper-statefulset-0            1/1     Running     0          8m16s


### Verify if Contrail Commmand is working

You can access Contrail Command application via web browser. Before that you have to forward the network traffic from localhost to Command's pod.

    kubectl port-forward $(kubectl get pods -l command=command -n contrail -o name) -n contrail 9091:9091

Go to http://localhost:9091

Authenticate using `admin` username and `contrail123` password. 

## Run unit tests

You can run unit test tests on your favourite IDE by executing all tests in `pkg` package.

You can also use command line tool:

    go test ./pkg/...

Eventually you should get results and command should return success:

    ?       github.com/Juniper/contrail-operator/pkg/apis    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/apis/contrail    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/tests    0.943s
    ?       github.com/Juniper/contrail-operator/pkg/client/keystone    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/client/kubeproxy    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/client/swift    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/controller/cassandra    1.414s
    ok      github.com/Juniper/contrail-operator/pkg/controller/command    2.316s
    ?       github.com/Juniper/contrail-operator/pkg/controller/config    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/control    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/enqueue    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/controller/keystone    4.196s
    ?       github.com/Juniper/contrail-operator/pkg/controller/kubemanager    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/controller/manager    1.789s
    ?       github.com/Juniper/contrail-operator/pkg/controller/manager/crs    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/controller/memcached    1.097s
    ok      github.com/Juniper/contrail-operator/pkg/controller/postgres    1.779s
    ?       github.com/Juniper/contrail-operator/pkg/controller/provisionmanager    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/rabbitmq    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/controller/swift    1.002s
    ok      github.com/Juniper/contrail-operator/pkg/controller/swiftproxy    0.870s
    ok      github.com/Juniper/contrail-operator/pkg/controller/swiftstorage    1.147s
    ?       github.com/Juniper/contrail-operator/pkg/controller/utils    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/vrouter    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/webui    [no test files]
    ?       github.com/Juniper/contrail-operator/pkg/controller/zookeeper    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/job    0.389s
    ok      github.com/Juniper/contrail-operator/pkg/k8s    0.558s
    ?       github.com/Juniper/contrail-operator/pkg/randomstring    [no test files]
    ok      github.com/Juniper/contrail-operator/pkg/swift/ring    0.416s
    ok      github.com/Juniper/contrail-operator/pkg/volumeclaims    0.633s
    
## Run e2e tests (aka system tests)

In order to test if the whole system works as expected we have a few plumbing tests. They verify if after deployment all Contrail services can talk to each other and operate as expected.

Before tests can be run you have to have clean the cluster. The fastest way is to delete the cluster:

    kind delete cluster

Then you have to create a new one plus a `contrail` namespace:

    cd test/env
    ./create_testenv.sh
    kubectl create namespace contrail

System tests can be run using operator-sdk tool

    # From contrail-operator root directory
    operator-sdk test local ./test/e2e/ --namespace contrail --go-test-flags "-v -timeout=30m" --up-local

