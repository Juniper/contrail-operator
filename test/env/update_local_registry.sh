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
contrail-operator/engprod-269421/ringcontroller:master.latest
contrail-operator/engprod-269421/contrail-operator:master.latest
contrail-operator/engprod-269421/contrail-statusmonitor:master.latest
contrail-operator/engprod-269421/contrail-operator-provisioner:master.latest
contrail-operator/engprod-269421/crdsloader:master.latest
contrail-nightly/contrail-command:master.1312
contrail-nightly/contrail-controller-config-api:master.1312
contrail-nightly/contrail-controller-config-devicemgr:master.1312
contrail-nightly/contrail-controller-config-schema:master.1312
contrail-nightly/contrail-controller-config-svcmonitor:master.1312
contrail-nightly/contrail-controller-control-control:master.1312
contrail-nightly/contrail-controller-control-dns:master.1312
contrail-nightly/contrail-controller-control-named:master.1312
contrail-nightly/contrail-controller-webui-job:master.1312
contrail-nightly/contrail-controller-webui-web:master.1312
contrail-nightly/contrail-kubernetes-kube-manager:master.1312
contrail-nightly/contrail-kubernetes-cni-init:master.1312
contrail-nightly/contrail-node-init:master.1312
contrail-nightly/contrail-nodemgr:master.1312
contrail-nightly/contrail-controller-config-dnsmasq:master.1312
contrail-nightly/contrail-analytics-api:master.1312
contrail-nightly/contrail-analytics-collector:master.1312
contrail-nightly/contrail-analytics-query-engine:master.1312
contrail-nightly/contrail-controller-config-devicemgr:master.1312
contrail-nightly/contrail-controller-config-api:2008.10
contrail-nightly/contrail-controller-config-devicemgr:2008.10
contrail-nightly/contrail-controller-config-schema:2008.10
contrail-nightly/contrail-controller-config-svcmonitor:2008.10
contrail-nightly/contrail-controller-control-control:2008.10
contrail-nightly/contrail-controller-control-dns:2008.10
contrail-nightly/contrail-controller-control-named:2008.10
contrail-nightly/contrail-controller-webui-job:2008.10
contrail-nightly/contrail-controller-webui-web:2008.10
contrail-nightly/contrail-kubernetes-kube-manager:2008.10
contrail-nightly/contrail-kubernetes-cni-init:2008.10
contrail-nightly/contrail-node-init:2008.10
contrail-nightly/contrail-nodemgr:2008.10
contrail-nightly/contrail-controller-config-dnsmasq:2008.10
contrail-nightly/contrail-analytics-api:2008.10
contrail-nightly/contrail-analytics-collector:2008.10
contrail-nightly/contrail-analytics-query-engine:2008.10
contrail-nightly/contrail-controller-config-devicemgr:2008.10
contrail-nightly/contrail-command:2008.10
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
EOF
