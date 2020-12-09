#!/usr/bin/env bash

set -o errexit

export BUILD_SCM_BRANCH=${BUILD_SCM_BRANCH:-"master"}
export BUILD_SCM_REVISION=${BUILD_SCM_REVISION:-"latest"}

if [[ -z "$CEM_RELEASE" ]]; then
    ## If build branch match release branch patern then use it
    if [[ "$BUILD_SCM_BRANCH" =~ ^R2[0-9]{3}$ ]]; then
        export CEM_RELEASE="${BUILD_SCM_BRANCH//R}-latest"
    else
        export CEM_RELEASE="master.1417-ubi"
    fi
fi

E2E_TEST_SUITE=${E2E_TEST_SUITE:-aio}

DIR="$(cd "$(dirname "$0")" && pwd)/../../"
pushd $DIR

cat deploy/1-create-operator.yaml | \
    sed "s/:master.latest/:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}/g" | \
    kubectl apply -f -

if [[ "$LONG_TEST" == "yes" ]]; then
    TEST_CONFIGURATION="-timeout=120m"
else
    TEST_CONFIGURATION="-timeout=45m -test.short"
fi

## Operator-sdk e2e test framework requires namespacedMan and globalMan to be defined, however in our case
## all resources and crds are registered before and automatically. That is way empty.yaml files are provided here
## This is equivalent to "operator-sdk test local --no-setup"
go test -v ./test/e2e/$E2E_TEST_SUITE/... -namespacedMan test/env/deploy/empty.yaml -globalMan test/env/deploy/empty.yaml -root $DIR -parallel=8 -skipCleanupOnError=false $TEST_CONFIGURATION

kubectl delete -f deploy/1-create-operator.yaml
popd
