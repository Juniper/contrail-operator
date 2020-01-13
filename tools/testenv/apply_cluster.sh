#!/usr/bin/env sh

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"

cp ../../deploy/0-create-persistent-volumes.yaml deploy/0-create-persistent-volumes.yaml
cp ../../deploy/1-create-operator.yaml deploy/1-create-operator.yaml

kubectl --context "${KIND_CLUSTER_NAME}"-kind apply -k deploy/
