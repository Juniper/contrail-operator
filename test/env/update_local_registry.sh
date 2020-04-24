#!/bin/bash
# This script pulls contrail, postgres and other images and pushes them to local repo registry
set -o errexit

# Overwrite this variable if you want to push images to external registry i. ex.
# LOCAL_REPO=svl-artifactory.juniper.net/common-docker-third-party/contrail
LOCAL_REPO=localhost:"${INTERNAL_INSECURE_REGISTRY_PORT:-5000}"
KOLLA_NEW_TAG="train-2005"
CONTRAIL_TAG="master.1175"

# pull_image pulls docker images from external registries, tags and push to registry
# Arguments:
# 1) external registry i. ex. opencontrailnightly
# 2) image's name with tag
# 3) image with new tag (optional)
pull_image()
{
    docker pull $1/$2

    new_tagged_image=$2
    if [ ! -z "$3" ]
    then
      new_tagged_image=$3
    fi

    docker tag $1/$2 $LOCAL_REPO/$new_tagged_image
    docker push $LOCAL_REPO/$new_tagged_image
}


while read line; do
	pull_image opencontrailnightly ${line}
done <<EOF
contrail-controller-config-api:master.1175
contrail-controller-config-devicemgr:master.1175
contrail-controller-config-schema:master.1175
contrail-controller-config-svcmonitor:master.1175
contrail-controller-control-control:master.1175
contrail-controller-control-dns:master.1175
contrail-controller-control-named:master.1175
contrail-controller-webui-job:master.1175
contrail-controller-webui-web:master.1175
contrail-kubernetes-kube-manager:master.1175
contrail-kubernetes-cni-init:master.1175
contrail-node-init:master.1175
contrail-nodemgr:master.1175
contrail-analytics-api:master.1175
contrail-analytics-collector:master.1175
contrail-analytics-query-engine:master.1175
contrail-controller-config-devicemgr:master.1175
EOF

# pinned kolla images (22/Apr/2020) tagged to "train-2005"
while read line; do
	pull_image kolla ${line}
done <<EOF
centos-binary-keystone-fernet@sha256:053c16581e30112ac5e17d3821143f8873c4d2ea47201a57fcac47ce2e2e8a33           centos-binary-keystone-fernet:${KOLLA_NEW_TAG}
centos-binary-keystone@sha256:c16e38faafb4bbbae2e2fc3422afc6f54f3ff457e25e96845192f56b0fefa22e                  centos-binary-keystone:${KOLLA_NEW_TAG}
centos-binary-keystone-ssh@sha256:202a5f3c13a7807c6f6a42ff49488052e5592644f4ed7106252038791e93502e              centos-binary-keystone-ssh:${KOLLA_NEW_TAG}
centos-binary-swift-account@sha256:22a2e3fa407bd303bbf63cd212ab682e99332ef9d2cb7a928f9614334dc07c2a             centos-binary-swift-account:${KOLLA_NEW_TAG}
centos-binary-swift-container@sha256:3d21baea74623fa8f53dba07701a0ef2b394c9af3354932a7c29be3b54a9079f           centos-binary-swift-container:${KOLLA_NEW_TAG}
centos-binary-swift-object-expirer@sha256:9485f886d015d4ce5f56f823a09293527b0cfe486162bd19f9ef9cd1382289bc      centos-binary-swift-object-expirer:${KOLLA_NEW_TAG}
centos-binary-swift-object@sha256:f1d0f5b988537755fc98f559ff3c692993f5defcf5882d0bce890c1e2c53d50b              centos-binary-swift-object:${KOLLA_NEW_TAG}
centos-binary-swift-proxy-server@sha256:2eddb9260c04993a42b7042cb3becac6e152eb1658ad12fa37e6984c85170ff1        centos-binary-swift-proxy-server:${KOLLA_NEW_TAG}
centos-binary-swift-rsyncd@sha256:f77883c4df839fec7b250a3f732f0ab9dd82c95649fb9d9b2bb0a7e320866c66              centos-binary-swift-rsyncd:${KOLLA_NEW_TAG}
centos-binary-kolla-toolbox@sha256:0bf3d4746f4aa0b735d274f362f390ceccdfea4448409c560e4a0d96ab220288             centos-binary-kolla-toolbox:${KOLLA_NEW_TAG}
centos-source-swift-base@sha256:75e37961ef80e6d41019627255f1c2c1a8fb6631b39e79690281b467976ca992                centos-source-swift-base:${KOLLA_NEW_TAG}
centos-binary-memcached@sha256:6348d4b541b749702f989fbf7ec589d323839edfd6ba8c2227b81e06dab4fa0c                 centos-binary-memcached:${KOLLA_NEW_TAG}
EOF

while read line; do
	pull_image docker.io "${line}"
done <<EOF
busybox:1.31
cassandra:3.11.4
zookeeper:3.5.5
postgres:12.2
python:3.8.2-alpine
redis:4.0.2
rabbitmq:3.7.16
EOF

pull_image tmaier postgresql-client@sha256:de4887ec672e3bc00d604b8d1cac09825ceaf5fc465cb1ca85ad76eee383e001 postgresql-client:1.0

while read line; do
	pull_image kaweue ${line}
done <<EOF
contrail-statusmonitor:master-180ab9
contrail-provisioner:master.1175
EOF

pull_image hub.juniper.net/contrail-nightly contrail-command:master.1175
pull_image hub.juniper.net/contrail-nightly contrail-controller-config-dnsmasq:master.1175
pull_image 10.84.18.17:5000 contrail-statusmonitor:latest