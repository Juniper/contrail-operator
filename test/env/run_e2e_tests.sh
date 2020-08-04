#!/usr/bin/env bash

set -o errexit

export BUILD_SCM_BRANCH=${BUILD_SCM_BRANCH:-"master"}
export BUILD_SCM_REVISION=${BUILD_SCM_REVISION:-"latest"}

if [[ -z "$CEM_RELEASE" ]]; then
    ## If build branch match release branch patern then use it
    if [[ "$BUILD_SCM_BRANCH" =~ ^R2[0-9]{3}$ ]]; then
        export CEM_RELEASE="${BUILD_SCM_BRANCH//R}-latest"
    else
        export CEM_RELEASE="master-latest"
    fi
fi
    

E2E_TEST_SUITE=${E2E_TEST_SUITE:-aio}

DIR="$(cd "$(dirname "$0")" && pwd)/../../"
pushd $DIR

operator_image="registry:5000/contrail-operator/engprod-269421/contrail-operator:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}"

kubectl apply -f deploy/cluster_role.yaml
kubectl apply -f deploy/cluster_role_binding.yaml

operator-sdk test local --verbose ./test/e2e/$E2E_TEST_SUITE --image "$operator_image" --go-test-flags "-timeout=30m -parallel=8"

kubectl delete -f deploy/cluster_role.yaml
kubectl delete -f deploy/cluster_role_binding.yaml

popd
