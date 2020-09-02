#!/bin/bash

while true; do
  FQHOSTNAME=$(hostname -f)
  if [[ $FQHOSTNAME != 'localhost' ]] && [[ $FQHOSTNAME != 'localhost.localdomain' ]]; then
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

