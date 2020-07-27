#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl apply -f ../../deploy/1-prepare-namespace.yaml
kubectl apply -f ../../deploy/2-create-operator.yaml
kubectl apply -f ../../deploy/kind/secret.yaml
kubectl apply -f ../../deploy/kind/cluster.yaml