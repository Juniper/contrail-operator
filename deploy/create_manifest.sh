#!/usr/bin/env bash

cat namespace.yaml > 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
cat role.yaml >> 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
cat cluster_role.yaml >> 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
cat service_account.yaml >> 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
cat role_binding.yaml >> 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
cat cluster_role_binding.yaml >> 1-prepare-namespace.yaml
echo "---" >> 1-prepare-namespace.yaml
for i in $(ls crds/*_crd.yaml); do
  cat $i >> 1-prepare-namespace.yaml
  echo "---" >> 1-prepare-namespace.yaml
done
cat operator.yaml > 2-create-operator.yaml

cat crds/contrail_v1alpha1_manager_cr.yaml > 3-start-operator-1node.yaml
echo "---" > 3-start-operator-3node.yaml
sed 's/replicas: 1/replicas: 3/g' 3-start-operator-1node.yaml > 3-start-operator-3node.yaml
