#!/usr/bin/env sh

set -o errexit

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"
EXTERNAL_INSECURE_REGISTRY="${EXTERNAL_INSECURE_REGISTRY:-127.0.0.1:5000}"
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

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name "${KIND_CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches: 
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry:5000"]
    endpoint = ["http://registry:5000"]
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."${EXTERNAL_INSECURE_REGISTRY}"]
    endpoint = ["http://${EXTERNAL_INSECURE_REGISTRY}"]
EOF

# add the registry to /etc/hosts on each node
ip_fmt='{{.NetworkSettings.IPAddress}}'
cmd="echo $(docker inspect -f "${ip_fmt}" "${reg_name}") registry >> /etc/hosts"
for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
  docker exec "${node}" sh -c "${cmd}"
done

docker cp -L /etc/localtime "${KIND_CLUSTER_NAME}"-control-plane:/etc/localtime
