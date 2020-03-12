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

push_local_images()
{
  DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
  make -f $DIR/../../Makefile provisioner && \
  docker tag contrail-provisioner:debug $LOCAL_REPO/contrail-provisioner:debug && \
  docker push $LOCAL_REPO/contrail-provisioner:debug
}

push_local_images

while read line; do
	pull_image opencontrailnightly "${line}"
done <<EOF
contrail-controller-config-api:1912-latest
contrail-controller-config-devicemgr:1912-latest
contrail-controller-config-schema:1912-latest
contrail-controller-config-svcmonitor:1912-latest
contrail-controller-control-control:1912-latest
contrail-controller-control-dns:1912-latest
contrail-controller-control-named:1912-latest
contrail-controller-webui-job:1912-latest
contrail-controller-webui-web:1912-latest
contrail-kubernetes-kube-manager:1912-latest
contrail-kubernetes-cni-init:1912-latest
contrail-node-init:1912-latest
contrail-nodemgr:1912-latest
contrail-analytics-api:1912-latest
contrail-analytics-collector:1912-latest
contrail-analytics-query-engine:1912-latest
contrail-external-dnsmasq:1910-latest
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
	pull_image michaelhenkel "${line}"
done <<EOF
contrail-statusmonitor:debug
EOF

while read line; do
	pull_image 10.84.18.17:5000 "${line}"
done <<EOF
contrail-provisioner:1912-20200315
EOF

while read line; do
	pull_image kaweue "${line}"
done <<EOF
contrail-controller-config-devicemgr:1912-reload
contrail-controller-config-dnsmasq:dev
EOF

pull_image hub.juniper.net/contrail-nightly contrail-command:1912-latest
