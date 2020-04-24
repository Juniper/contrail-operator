#!/bin/bash
# This script pulls contrail, postgres and other images and pushes them to local repo registry
set -o errexit

LOCAL_REPO=localhost:"${INTERNAL_INSECURE_REGISTRY_PORT:-5000}"

CONTRAIL_TAG="${CONTRAIL_TAG:-master.latest}"
CONTRAIL_REGISTRY="${CONTRAIL_REGISTRY:-svl-artifactory.juniper.net/contrail-nightly}"
OPERATOR_TAG="${OPERATOR_TAG:-master.latest}"
OPERATOR_REGISTRY="${OPERATOR_REGISTRY:-svl-artifactory.juniper.net/contrail-operator.gcr.io/eng-prod-237922}"

pull_image() 
{
    docker pull $1/$2 
    docker tag $1/$2 $LOCAL_REPO/$2
    docker push $LOCAL_REPO/$2
}

while read line; do
	pull_image ${CONTRAIL_REGISTRY} "${line}"
done <<EOF
contrail-controller-config-api:${CONTRAIL_TAG}
contrail-controller-config-devicemgr:${CONTRAIL_TAG}
contrail-controller-config-schema:${CONTRAIL_TAG}
contrail-controller-config-svcmonitor:${CONTRAIL_TAG}
contrail-controller-control-control:${CONTRAIL_TAG}
contrail-controller-control-dns:${CONTRAIL_TAG}
contrail-controller-control-named:${CONTRAIL_TAG}
contrail-controller-webui-job:${CONTRAIL_TAG}
contrail-controller-webui-web:${CONTRAIL_TAG}
contrail-kubernetes-kube-manager:${CONTRAIL_TAG}
contrail-kubernetes-cni-init:${CONTRAIL_TAG}
contrail-node-init:${CONTRAIL_TAG}
contrail-nodemgr:${CONTRAIL_TAG}
contrail-analytics-api:${CONTRAIL_TAG}
contrail-analytics-collector:${CONTRAIL_TAG}
contrail-analytics-query-engine:${CONTRAIL_TAG}
contrail-controller-config-devicemgr:${CONTRAIL_TAG}
contrail-command:${CONTRAIL_TAG}
contrail-controller-config-dnsmasq:${CONTRAIL_TAG}
EOF

while read line; do
	pull_image kolla "${line}"
done <<EOF
centos-binary-keystone-fernet:train
centos-binary-keystone:train
centos-binary-keystone-ssh:train
centos-binary-swift-account:train
centos-binary-swift-container:train
centos-binary-swift-object-expirer:train
centos-binary-swift-object:train
centos-binary-swift-proxy-server:train
centos-binary-swift-rsyncd:train
centos-binary-kolla-toolbox:train
centos-source-swift-base:train
centos-binary-memcached:train
EOF

while read line; do
	pull_image docker.io "${line}"
done <<EOF
busybox:latest
cassandra:3.11.4
cassandra:3.11.3
zookeeper:3.5.5
zookeeper:3.5.4-beta
postgres
python:alpine
redis:4.0.2
rabbitmq:3.7
rabbitmq:3.7.16
EOF

pull_image tmaier postgresql-client

while read line; do
	pull_image ${OPERATOR_REGISTRY} "${line}"
done <<EOF
contrail-statusmonitor:${OPERATOR_TAG}
contrail-provisioner:${OPERATOR_TAG}
EOF
