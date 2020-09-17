#!/bin/bash

# In case of CoreOS reboot /var/run directory is wiped out.
# Since ContrailCNI Job will not be run again then it's necessary to
# copy multus config with oneshot service triggering this script
while [ ! -f /etc/cni/net.d/10-contrail.conf ]
do
  sleep 2
done
cp -f /etc/cni/net.d/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf
