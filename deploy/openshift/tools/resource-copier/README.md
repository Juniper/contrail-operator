# RESOURCE COPIER

## What is it?
This script allows to copy one resource from one Kubernetes cluster to another (or from one namespace to another namespace in scope of the same cluster).

## What do I need to use it?
CLuster selection is based on kubeconfig file. So in order to access clusters, prepare kubeconfig file for both cluster from which you want to copy resource and cluster where you want to create copy of resource.

## Usage
"""
./resource-copier.sh --from-cluster <kubeconfig file of cluster with existing resource> --to-cluster <kubeconfig file of cluster where to copy resource> --from-namespace <target namespace where is original resource> --to-namespace <target namespace where to copy resource> --resource <resource type> --name <name of the resource>
"""

If one of kubeconfigs is not passed then script will automatically perform copy in scope of the same cluster.
If one of namespaces is not passed then script will perform copy between same namespace names.

