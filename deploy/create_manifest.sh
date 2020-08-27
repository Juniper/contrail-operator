#!/usr/bin/env bash

MANIFEST_FILE="1-create-operator.yaml"
PREREQUISITES="namespace.yaml \
role.yaml \
cluster_role.yaml \
service_account.yaml \
role_binding.yaml \
cluster_role_binding.yaml"

while getopts ":l" opt; do
  case ${opt} in
    l ) LOCAL_OPERATOR_RUN="yes"; MANIFEST_FILE="local_operator.yaml"
      ;;
  esac
done

function combine_manifests {
  > "$MANIFEST_FILE"
  for doc in $@; do
     echo "---" >> "$MANIFEST_FILE"
     cat "$doc" >> "$MANIFEST_FILE"
  done
}

pushd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null

if [ -z "${LOCAL_OPERATOR_RUN}" ]; then
  combine_manifests $PREREQUISITES "operator.yaml"
  cat crds/contrail_v1alpha1_manager_cr.yaml > 2-start-operator-1node.yaml
  echo "---" > 2-start-operator-3node.yaml
  sed 's/replicas: 1/replicas: 3/g' 2-start-operator-1node.yaml > 2-start-operator-3node.yaml
else
  combine_manifests $PREREQUISITES crds/*_crd.yaml
fi

popd > /dev/null