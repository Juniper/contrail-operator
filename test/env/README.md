# Development environment in KIND (Kubernetes IN Docker)

## Install KIND (Kubernetes IN Docker)

    GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0

## Create a test env

### All in one (one-node kind cluster)
    export KIND_CLUSTER_NAME=kind
    export INTERNAL_INSECURE_REGISTRY_PORT=6000
    # number of intended cluster nodes, defaults to 1
    export NODES=1
    ./create_testenv.sh

It creates Kubernetes IN Docker cluster with a docker registry. This docker registry is accessible from host at `localhost:6000` and from inside the cluster at `registry:5000`

### Multi-node kind cluster

To run multi-node cluster you will need sufficient amount of memory (RAM) and CPU to handle multiple replicas of services e.g. `32 GB RAM` and `12 CPU`. Following steps provide an instruction for internal usage how to deploy VM on ESXI with docker installed on and Docker Daemon configured to listen on the `2375` port. If you have a dedicated machine wich fullfills requirements you can start from the step 3.

1. Create a VM which will be your docker host for kind. You can do this by using ready-to-use OVA which is available under `/root/dockerd/dockerd.ova` path on `b5s5-node3` in our lab. You can deploy a VM using following command:

        ovftool --X:logFile=upload.log --X:logLevel=verbose --name=YOUR_VM_NAME /root/dockerd/dockerd.ova  vi://root:c0ntrail123@b5s5-node4.englab.juniper.net
1. Run your machine using default credentials
1. You need to set environemt variable on your local machine with address of Docker Daemon from newly created machine (docker host). To do this type following command on your local terminal:
        
        export DOCKER_HOST=tcp://ip_of_dockerd_vm:2375
1. Kubectl expects k8s APIServer listening on 127.0.0.1 and some random port assigned by Docker. To change k8s `APIServer IP` and `port`, you will have to modify `create_k8s_cluster.sh` script. Change kindConfig by adding those two lines in `networking` section:

        apiServerAddress: "ip_of_dockerd_vm"
        apiServerPort: 6443
1. Create test env using `create_testenv.sh` (remember to set variable NODES to 3). First two variables are optional and have default values.

        export KIND_CLUSTER_NAME=kind
        export INTERNAL_INSECURE_REGISTRY_PORT=6000
         # number of intended cluster nodes, defaults to 1
        export NODES=3
        ./create_testenv.sh

## Download all required images to a local registry
    export INTERNAL_INSECURE_REGISTRY_PORT=6000
    ./update_local_registry.sh

## Apply operator and cluster
    export KIND_CLUSTER_NAME=kind
    ./apply_contrail_cluster.sh

## Destroy operator and cluster

    ./clear_contrail_cluster.sh

## Delete cluster

    kind delete cluster

# E2E tests

## Run test

    # From contrail-operator root directory
    # Use operator-sdk version >= v.0.13
    kubectl create namespace contrail
    operator-sdk test local ./test/e2e/ --namespace contrail --go-test-flags "-v -timeout=30m" --up-local
