#!/usr/bin/env sh

kind delete cluster
kind create cluster --config "/workspace/test/gcb/cluster.yaml"
IP=$(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" kind-external-load-balancer)
echo $IP
kubectl config set-cluster kind-kind --server=https://$IP:6443 --insecure-skip-tls-verify
sleep 20
kubectl cluster-info
kubectl get nodes -o wide
kind get nodes
kind delete cluster