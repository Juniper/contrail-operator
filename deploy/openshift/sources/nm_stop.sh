#!/bin/bash

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

