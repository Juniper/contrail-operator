#!/usr/bin/env bash

set -o errexit

BUILD_SCM_BRANCH=${BUILD_SCM_BRANCH:-"master"}
BUILD_SCM_REVISION=${BUILD_SCM_REVISION:-"latest"}
CEM_RELEASE=${CEM_RELEASE:-"master.1320"}
SKIP_TEST=${SKIP_TEST:-no}
E2E_TEST_SUITE=${E2E_TEST_SUITE:-aio}

DIR="$(cd "$(dirname "$0")" && pwd)/../../"
pushd $DIR

operator_image="registry:5000/contrail-operator/engprod-269421/contrail-operator:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}"

kubectl apply -f deploy/cluster_role.yaml
kubectl apply -f deploy/cluster_role_binding.yaml

operator-sdk test local --verbose ./test/e2e/$E2E_TEST_SUITE --image "$operator_image" --go-test-flags "-singleNamespace -timeout=30m -parallel=8"

kubectl delete -f deploy/cluster_role.yaml
kubectl delete -f deploy/cluster_role_binding.yaml

popd
