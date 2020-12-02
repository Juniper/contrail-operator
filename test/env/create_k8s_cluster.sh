#!/usr/bin/env bash

set -o errexit

if [[ $# != 3 ]] ; then
    echo "Usage $(basename "$0") cluster_name insecure_registry_address"
    echo " - cluster_name a name of KIND cluster e.g mycluster"
    echo " - insecure_registry_address address of external insecure registry e.g 192.168.0.2"
    echo " - number of intended cluster nodes"
    exit 1
fi

# desired cluster name; default is "kind"
kind_cluster_name="$1"
insecure_registry_address="$2"
number_of_nodes="$3"

# TODO disable CNI when PV won't require storage provisioner
kindConfig="$(cat << EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry:5000"]
    endpoint = ["http://registry:5000"]
networking:
  disableDefaultCNI: true
  apiServerAddress: 0.0.0.0
  apiServerPort: 6443
nodes:
- role: control-plane
EOF
)"

node="
- role: worker"

for (( i=0; i<(($number_of_nodes-1)); i++ ))
do
  kindConfig+="$node"
done

# create a cluster with the local registry enabled in containerd
echo "$kindConfig" | kind create cluster --name "${kind_cluster_name}" --config=-

# add the registry to /etc/hosts on each node
cmd="echo ${insecure_registry_address} registry >> /etc/hosts"
for node in $(kind get nodes --name "${kind_cluster_name}"); do
  docker exec "${node}" sh -c "${cmd}"
done

# remove default storage class
kubectl delete sc standard

# untaint nodes
kubectl taint nodes --all node-role.kubernetes.io/master- &>/dev/null || true

# label nodes
kubectl label nodes --all node-role.juniper.net/contrail=""
