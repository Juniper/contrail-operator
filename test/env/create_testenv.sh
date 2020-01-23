#!/usr/bin/env bash

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

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
${DIR}/create_k8s_cluster.sh $KIND_CLUSTER_NAME $(docker inspect -f "{{.NetworkSettings.IPAddress}}" "${reg_name}")

docker cp -L /etc/localtime "${KIND_CLUSTER_NAME}"-control-plane:/etc/localtime
