#!/usr/bin/env sh

kind delete cluster
kind create cluster
IP=$(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" kind-control-plane)
echo $IP
kubectl config set-cluster kind-kind --server=https://$IP:6443 --insecure-skip-tls-verify
kubectl cluster-info
kind delete cluster;