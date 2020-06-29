#!/usr/bin/env bash

set -o errexit

BUILD_SCM_BRANCH=${BUILD_SCM_BRANCH:-"master"}
BUILD_SCM_REVISION=${BUILD_SCM_REVISION:-"latest"}

operator_image="registry:5000/contrail-operator/engprod-269421/contrail-operator:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}"

kubectl apply -f deploy/cluster_role.yaml
kubectl apply -f deploy/cluster_role_binding.yaml
operator-sdk test local ./test/e2e/ --image "$operator_image" --namespace contrail --go-test-flags "-v -timeout=30m"

kubectl delete -f deploy/cluster_role.yaml
kubectl delete -f deploy/cluster_role_binding.yaml