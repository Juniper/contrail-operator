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
	pull_image svl-artifactory.juniper.net "${line}"
done <<EOF
contrail-operator.gcr.io/eng-prod-237922/contrail-statusmonitor:configure-encap-priority-provisioner.bee1bc8
contrail-operator.gcr.io/eng-prod-237922/contrail-provisioner:configure-encap-priority-provisioner.bee1bc8
contrail-nightly/contrail-controller-config-api:master.1175
contrail-nightly/contrail-controller-config-devicemgr:master.1175
contrail-nightly/contrail-controller-config-schema:master.1175
contrail-nightly/contrail-controller-config-svcmonitor:master.1175
contrail-nightly/contrail-controller-control-control:master.1175
contrail-nightly/contrail-controller-control-dns:master.1175
contrail-nightly/contrail-controller-control-named:master.1175
contrail-nightly/contrail-controller-webui-job:master.1175
contrail-nightly/contrail-controller-webui-web:master.1175
contrail-nightly/contrail-kubernetes-kube-manager:master.1175
contrail-nightly/contrail-kubernetes-cni-init:master.1175
contrail-nightly/contrail-node-init:master.1175
contrail-nightly/contrail-nodemgr:master.1175
contrail-nightly/contrail-controller-config-dnsmasq:master.1175
contrail-nightly/contrail-analytics-api:master.1175
contrail-nightly/contrail-analytics-collector:master.1175
contrail-nightly/contrail-analytics-query-engine:master.1175
contrail-nightly/contrail-controller-config-devicemgr:master.1175
contrail-nightly/contrail-command:master.1175
common-docker-third-party/contrail/postgresql-client:1.0
common-docker-third-party/contrail/busybox:1.31
common-docker-third-party/contrail/cassandra:3.11.4
common-docker-third-party/contrail/cassandra:3.11.3
common-docker-third-party/contrail/zookeeper:3.5.5
common-docker-third-party/contrail/zookeeper:3.5.4-beta
common-docker-third-party/contrail/postgres:12.2
common-docker-third-party/contrail/python:3.8.2-alpine
common-docker-third-party/contrail/redis:4.0.2
common-docker-third-party/contrail/rabbitmq:3.7
common-docker-third-party/contrail/rabbitmq:3.7.16
common-docker-third-party/contrail/centos-binary-keystone-fernet:train-2005
common-docker-third-party/contrail/centos-binary-keystone:train-2005
common-docker-third-party/contrail/centos-binary-keystone-ssh:train-2005
common-docker-third-party/contrail/centos-binary-swift-account:train-2005
common-docker-third-party/contrail/centos-binary-swift-container:train-2005
common-docker-third-party/contrail/centos-binary-swift-object-expirer:train-2005
common-docker-third-party/contrail/centos-binary-swift-object:train-2005
common-docker-third-party/contrail/centos-binary-swift-proxy-server:train-2005
common-docker-third-party/contrail/centos-binary-swift-rsyncd:train-2005
common-docker-third-party/contrail/centos-binary-kolla-toolbox:train-2005
common-docker-third-party/contrail/centos-source-swift-base:train-2005
common-docker-third-party/contrail/centos-binary-memcached:train-2005
EOF
