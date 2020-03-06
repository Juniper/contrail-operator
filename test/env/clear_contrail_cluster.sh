#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl --context "${KIND_CLUSTER_NAME}"-kind delete -k deploy/
kubectl delete -f deploy/secret.yaml
kubectl delete -f deploy/cluster.yaml
kubectl delete -f ../../deploy/1-create-operator.yaml

