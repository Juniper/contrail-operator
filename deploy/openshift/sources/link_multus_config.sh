#!/bin/bash

while [ ! -f /tmp/list.txt ]
do
  sleep 2
done
ln /etc/cni/net.d/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf
