#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl delete -f deploy/secret.yaml
kubectl delete --wait=true -f deploy/cluster.yaml
kubectl delete -f ../../deploy/1-create-operator.yaml
kubectl delete pv $(kubectl get pv -o=jsonpath='{.items[?(@.spec.storageClassName=="local-storage")].metadata.name}')

for p in $(docker ps --filter name=kind-control --filter name=kind-worker -q)
do
    docker exec $p rm -rf /var/lib/cassandra
    docker exec $p rm -rf /var/lib/zookeeper
done