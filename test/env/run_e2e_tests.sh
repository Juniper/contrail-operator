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

## Use fixed working version of CEM
export CEM_RELEASE="master.1312"

E2E_TEST_SUITE=${E2E_TEST_SUITE:-aio}

DIR="$(cd "$(dirname "$0")" && pwd)/../../"
pushd $DIR

cat deploy/1-create-operator.yaml | \
    sed "s/:master.latest/:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}/g" | \
    kubectl apply -f -

go test -v ./test/e2e/$E2E_TEST_SUITE/... -namespacedMan test/env/deploy/empty.yaml -globalMan test/env/deploy/empty.yaml -root $DIR -timeout=30m -parallel=8 -skipCleanupOnError=false

kubectl delete -f deploy/1-create-operator.yaml
popd
