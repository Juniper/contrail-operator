#!/bin/bash

mkdir -p /var/run/multus/cni/net.d
touch /etc/cni/net.d/10-contrail.conf
ln -s /etc/cni/net.d/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf
