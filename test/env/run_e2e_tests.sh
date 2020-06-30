#!/usr/bin/env bash

set -o errexit

BUILD_SCM_BRANCH=${BUILD_SCM_BRANCH:-"master"}
BUILD_SCM_REVISION=${BUILD_SCM_REVISION:-"latest"}
SKIP_TEST=${SKIP_TEST:-no}

if [[ $BUILD_SCM_REVISION != "latest" ]] ; then
    while : ; do
        check_suites=$(curl -s -H "Accept: application/vnd.github.antiope-preview+json" https://api.github.com/repos/Juniper/contrail-operator/commits/$BUILD_SCM_REVISION/check-suites)
        status=$(echo "$check_suites" | jq -cr '.check_suites[]|select(.app.slug == "google-cloud-build").status')
        [[ $status != "completed" ]] && echo "Waiting for upstream job. Current status: $status" && sleep 120 && continue
        conclusion=$(echo "$check_suites" | jq -cr '.check_suites[]|select(.app.slug == "google-cloud-build").conclusion')
        [[ $conclusion == "success" ]] && break
        echo "Upstream job failed with conclusion: $conclusion"
        exit 1
    done
fi


DIR="$(cd "$(dirname "$0")" && pwd)/../../"
pushd $DIR

operator_image="registry:5000/contrail-operator/engprod-269421/contrail-operator:${BUILD_SCM_BRANCH}.${BUILD_SCM_REVISION}"

kubectl apply -f deploy/cluster_role.yaml
kubectl apply -f deploy/cluster_role_binding.yaml

if [[ $SKIP_TEST != "yes" ]]; then
    operator-sdk test local ./test/e2e/ --image "$operator_image" --namespace contrail --go-test-flags "-v -timeout=30m"
fi

kubectl delete -f deploy/cluster_role.yaml
kubectl delete -f deploy/cluster_role_binding.yaml

popd
