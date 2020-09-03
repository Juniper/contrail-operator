#!/usr/bin/env bash
# This script will run Contrail operator locally with delve debugger enabled.
# See DEVELOPMENT.md for more details.
set -e

DIR="$(dirname "${BASH_SOURCE[0]}")"

${DIR}/../../deploy/create_manifest.sh -l
kubectl apply -f ${DIR}/../../deploy/local_operator.yaml
kubectl apply -f ${DIR}/deploy/secret.yaml
kubectl apply -f ${DIR}/deploy/cluster.yaml
