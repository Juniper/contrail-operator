#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl delete -f ../../deploy/2-create-operator.yaml
kubectl delete -f ../../deploy/kind/cluster.yaml
kubectl delete -f ../../deploy/kind/secret.yaml
kubectl delete -f ../../deploy/1-prepare-namespace.yaml

