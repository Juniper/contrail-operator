# Development environment in KIND (Kubernetes IN Docker)

## Install KIND (Kubernetes IN Docker)

    GO111MODULE="on" go get sigs.k8s.io/kind@v0.6.1

## Create a test env
    export KIND_CLUSTER_NAME=kind
    export EXTERNAL_INSECURE_REGISTRY=172.17.14.127:5000
    export INTERNAL_INSECURE_REGISTRY_PORT=6000
    ./create_testenv.sh

It creates Kubernetes IN Docker cluster with a docker registry. This docker registry is accessible from host at `localhost:6000` and from inside the cluster at `registry:5000`

## Create keystone ssh keys

    ssh-keygen -t rsa -b 1024 -N "" -f deploy/id_rsa

## Apply operator and cluster
    export KIND_CLUSTER_NAME=kind
    ./apply_cluster.sh

## Destroy operator and cluster

    kubectl --context "${KIND_CLUSTER_NAME}"-kind delete -k deploy/

## Delete cluster

    kind delete cluster