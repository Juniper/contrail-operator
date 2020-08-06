#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl delete -f deploy/secret.yaml
kubectl delete -f deploy/cluster.yaml
kubectl delete -f ../../deploy/1-create-operator.yaml
kubectl delete pv $(kubectl get pv -o=jsonpath='{.items[?(@.spec.storageClassName=="local-storage")].metadata.name}')
