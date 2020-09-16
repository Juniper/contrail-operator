#!/bin/bash

while [ ! -f /etc/cni/net.d/10-contrail.conf ]
do
  sleep 2
done
ln -s /etc/cni/net.d/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf
