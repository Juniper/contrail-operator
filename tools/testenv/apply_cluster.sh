#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl apply -f ../../deploy/1-create-operator.yaml
kubectl apply -f deploy/cluster.yaml

kubectl --context "${KIND_CLUSTER_NAME}"-kind apply -k deploy/
