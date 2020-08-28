#!/bin/bash

while true;
do
  if [[ -L "/sys/class/net/vhost0" && $(ip address show vhost0 | grep inet[^6]) ]];
  then
          echo "[INFO] Detected vhost0 interface. Stopping NetworkManager..."
          # On RHCOS NetworkManager manages the hostname (using dhcp), so before stopping it
          # we have to set a static hostname to the current fqdn hostname.
          hostnamectl set-hostname $(hostname -A)
          systemctl stop NetworkManager
          echo "[INFO] Networkmanager stopped."
  fi
  sleep 10
done

