# Development environment in KIND (Kubernetes IN Docker)

## Install KIND (Kubernetes IN Docker)

    GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0

## Create a test env
    export KIND_CLUSTER_NAME=kind
    export INTERNAL_INSECURE_REGISTRY_PORT=6000
    ./create_testenv.sh

It creates Kubernetes IN Docker cluster with a docker registry. This docker registry is accessible from host at `localhost:6000` and from inside the cluster at `registry:5000`

## Download all required images to a local registry
    export INTERNAL_INSECURE_REGISTRY_PORT=6000
    ./update_local_registry.sh

## Create keystone ssh keys

    ssh-keygen -t rsa -b 1024 -N "" -f deploy/id_rsa

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
    operator-sdk test local ./test/e2e/ --namespace contrail --go-test-flags "-v" --up-local
