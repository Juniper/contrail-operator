#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

kubectl --context "${KIND_CLUSTER_NAME}"-kind apply -k deploy/
kubectl delete -f ../../deploy/cluster.yaml
kubectl delete -f ../../deploy/1-create-operator.yaml
kubectl delete -f ../../deploy/0-create-persistent-volumes.yaml

# create directories for persistent volumes
docker exec "${KIND_CLUSTER_NAME}"-control-plane bash -c "for directory in \$(seq 5); do
  mkdir -p /mnt/volumes/$directory
  rm -rf /mnt/volumes/$directory/*
done"