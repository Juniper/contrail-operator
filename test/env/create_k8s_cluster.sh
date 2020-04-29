#!/usr/bin/env bash

set -o errexit

if [[ $# != 2 ]] ; then
    echo "Usage $(basename "$0") cluster_name insecure_registry_address"
    echo " - cluster_name a name of KIND cluster e.g mycluster"
    echo " - insecure_registry_address address of external insecure registry e.g 192.168.0.2"
    exit 1
fi

# desired cluster name; default is "kind"
kind_cluster_name="$1"
insecure_registry_address="$2"

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name "${kind_cluster_name}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry:5000"]
    endpoint = ["http://registry:5000"]
EOF

# add the registry to /etc/hosts on each node
cmd="echo ${insecure_registry_address} registry >> /etc/hosts"
for node in $(kind get nodes --name "${kind_cluster_name}"); do
  docker exec "${node}" sh -c "${cmd}"
done
