#!/bin/bash

while true;
do
  FQHOSTNAME=$(hostname -f)
  if [[ $FQHOSTNAME != 'localhost' ]] && [[ $FQHOSTNAME != 'localhost.localdomain' ]];
  then
    break
  fi
  sleep 2
done

while true; do
  FQHOSTNAME=$(hostname -f)
  echo "Setting static hostname to $FQHOSTNAME"
  hostnamectl set-hostname $FQHOSTNAME
  if [[ $? -eq 0 ]]; then
    break;
  fi
done

while true;
do
  if [[ -L "/sys/class/net/vhost0" && $(ip address show vhost0 | grep inet[^6]) ]];
  then
          echo "[INFO] Detected vhost0 interface. Stopping NetworkManager..."
          systemctl stop NetworkManager
          echo "[INFO] Networkmanager stopped."
  fi
  sleep 10
done

