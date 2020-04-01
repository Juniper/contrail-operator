#!/bin/bash
# This script pulls contrail, postgres and other images and pushes them to local repo registry
set -o errexit

LOCAL_REPO=localhost:"${INTERNAL_INSECURE_REGISTRY_PORT:-5000}"

pull_image() 
{
    docker pull $1/$2 
    docker tag $1/$2 $LOCAL_REPO/$2
    docker push $LOCAL_REPO/$2
}

while read line; do
	pull_image opencontrailnightly "${line}"
done <<EOF
contrail-controller-config-api:master.1115
contrail-controller-config-devicemgr:master.1115
contrail-controller-config-schema:master.1115
contrail-controller-config-svcmonitor:master.1115
contrail-controller-control-control:master.1115
contrail-controller-control-dns:master.1115
contrail-controller-control-named:master.1115
contrail-controller-webui-job:master.1115
contrail-controller-webui-web:master.1115
contrail-kubernetes-kube-manager:master.1115
contrail-kubernetes-cni-init:master.1115
contrail-node-init:master.1115
contrail-nodemgr:master.1115
contrail-analytics-api:master.1115
contrail-analytics-collector:master.1115
contrail-analytics-query-engine:master.1115
contrail-controller-config-devicemgr:master.1115
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
	pull_image kaweue "${line}"
done <<EOF
contrail-statusmonitor:debug
contrail-controller-config-dnsmasq:dev
contrail-provisioner:master.1115
EOF

pull_image hub.juniper.net/contrail-nightly contrail-command:master.1115
