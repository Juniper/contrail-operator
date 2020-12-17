#!/usr/bin/env sh
set -o errexit

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl apply -f ../../deploy/1-create-operator.yaml

kubectl wait deployment --for=condition=available --timeout=240s -n contrail contrail-operator
kubectl apply -f deploy/secret.yaml
kubectl apply -f deploy/cluster.yaml