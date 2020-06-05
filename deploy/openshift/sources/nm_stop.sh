#!/bin/bash

while true;
do
  if [[ -L "/sys/class/net/vhost0" ]];
  then
          echo "[INFO] Detected vhost0 interface. Stopping NetworkManager..."
          systemctl stop NetworkManager
          echo "[INFO] Networkmanager stopped."
  fi
  sleep 10
done