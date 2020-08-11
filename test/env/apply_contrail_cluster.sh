#!/usr/bin/env sh
set -o errexit

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl apply -f ../../deploy/1-create-operator.yaml

kubectl wait pod --for=condition=Ready --timeout=30s -n contrail -l name=contrail-operator
kubectl apply -f deploy/secret.yaml
kubectl apply -f deploy/cluster.yaml