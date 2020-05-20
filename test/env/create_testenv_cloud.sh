#!/usr/bin/env bash

# Creates enviroment (cluster) in cloud node and exposes it to outside world
# After it's done just copy ./kube/config on your local machine replace adress in config
# to public node address

set -o errexit

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"
INTERNAL_INSECURE_REGISTRY_PORT="${INTERNAL_INSECURE_REGISTRY_PORT:-5000}"

# create registry container unless it already exists
reg_name='kind-registry'
reg_port=${INTERNAL_INSECURE_REGISTRY_PORT}
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

insecure_registry_address=$(docker inspect -f "{{.NetworkSettings.IPAddress}}" "${reg_name}")

# Get the addresses
paublic_ip=$(curl -s ifconfig.me)
private_ip=$(ip route get 1 | awk '{print $NF;exit}')

#create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name "${KIND_CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: $private_ip
  apiServerPort: 6443
nodes:
- role: control-plane
- role: control-plane
- role: control-plane
- role: worker
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry:5000"]
    endpoint = ["http://registry:5000"]
EOF

# connect the registry to the cluster network
# docker network connect "kind" "${reg_name}"

# tell https://tilt.dev to use the registry
# https://docs.tilt.dev/choosing_clusters.html#discovering-the-registry
#for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
#  kubectl annotate node "${node}" "kind.x-k8s.io/registry=registry:5000";
#done

# Get kubeadmin config and add public address
#kubectl -n kube-system get configmap kubeadm-config -o jsonpath='{.data.ClusterConfiguration}' > kubeadm.yaml
#sed  -i "/^\s\scertSANs:.*/a\ \ - $paublic_ip" kubeadm.yaml
#
## On each control plane node regenerate cerificate with added public address
#control_plane_node_pattern="${KIND_CLUSTER_NAME}-control-plane*"
#for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
#  if [[ ${node} =~ ${control_plane_node_pattern} ]]; then
#    echo "configuring node $node"
#    docker cp kubeadm.yaml "${node}:/"
#    docker exec "${node}" bash -c "mv /etc/kubernetes/pki/apiserver.{crt,key} /etc/kubernetes/"
#    docker exec "${node}" bash -c "kubeadm init phase certs apiserver --config /kubeadm.yaml"
#  fi
#done

# add the registry to /etc/hosts on each node
cmd="echo ${insecure_registry_address} registry >> /etc/hosts"
for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
  docker exec "${node}" sh -c "${cmd}"
done