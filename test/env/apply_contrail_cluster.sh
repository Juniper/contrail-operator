#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl apply -f ../../deploy/1-create-operator.yaml

kubectl wait --for=condition=Ready --timeout=30s -n contrail -l name=contrail-operator
kubectl apply -f deploy/secret.yaml
kubectl apply -f deploy/cluster.yaml