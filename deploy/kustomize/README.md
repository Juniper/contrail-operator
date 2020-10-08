# Testing environment

Tested on a system setup using this [vagrantfile](misc/vagrant/Vagrantfile).    
Kubernetes installed using kubespray with [inventory](misc/kubespray/inventory.yaml) and    
this [deploy command](misc/kubespray/runKubespray.sh)

# Instructions

## Set Contrail release version
```
RELEASE=R2008
```
## Create Operator, roles, bindings and crds
```
kubectl apply -k \
  github.com/Juniper/contrail-operator/deploy/kustomize/operator/${RELEASE}
```
## Wait for the CRDs to be created
```
kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net
```
## Create Contrail
### Set scale to 1 or 3 node(s)
```
REPLICA=1
```
or    
```
REPLICA=3
```
### Deploy Contrail
```
kubectl apply -k \
  github.com/Juniper/contrail-operator/deploy/kustomize/contrail/${REPLICA}node/${RELEASE}
```
## Wait for Contrail to come up
```
kubectl -n contrail get pods
NAME                                          READY   STATUS      RESTARTS   AGE
cassandra1-cassandra-statefulset-0            1/1     Running     0          89m
cassandra1-cassandra-statefulset-1            1/1     Running     0          87m
cassandra1-cassandra-statefulset-2            1/1     Running     0          86m
cnimasternodes-contrailcni-job-87qrj          0/1     Completed   0          46m
cnimasternodes-contrailcni-job-8tf62          0/1     Completed   0          46m
cnimasternodes-contrailcni-job-pfx65          0/1     Completed   0          46m
cniworkernodes-contrailcni-job-7g8m6          0/1     Completed   0          19m
config1-config-statefulset-0                  10/10   Running     1          86m
config1-config-statefulset-1                  10/10   Running     2          86m
config1-config-statefulset-2                  10/10   Running     2          86m
contrail-operator-7578bd7c6f-rrvpf            1/1     Running     0          20m
control1-control-statefulset-0                4/4     Running     0          81m
control1-control-statefulset-1                4/4     Running     0          81m
control1-control-statefulset-2                4/4     Running     0          81m
kubemanager1-kubemanager-statefulset-0        2/2     Running     1          81m
kubemanager1-kubemanager-statefulset-1        2/2     Running     1          81m
kubemanager1-kubemanager-statefulset-2        2/2     Running     0          81m
provmanager1-provisionmanager-statefulset-0   1/1     Running     0          81m
rabbitmq1-rabbitmq-statefulset-0              1/1     Running     0          89m
rabbitmq1-rabbitmq-statefulset-1              1/1     Running     0          89m
rabbitmq1-rabbitmq-statefulset-2              1/1     Running     0          89m
vroutermaster-vrouter-daemonset-6klvd         1/1     Running     0          79m
vroutermaster-vrouter-daemonset-kz7nq         1/1     Running     0          79m
vroutermaster-vrouter-daemonset-wnhrn         1/1     Running     0          79m
vrouternodes-vrouter-daemonset-4qhcx          1/1     Running     0          79m
webui1-webui-statefulset-0                    3/3     Running     1          81m
webui1-webui-statefulset-1                    3/3     Running     1          81m
webui1-webui-statefulset-2                    3/3     Running     0          81m
zookeeper1-zookeeper-statefulset-0            1/1     Running     0          89m
zookeeper1-zookeeper-statefulset-1            1/1     Running     0          87m
zookeeper1-zookeeper-statefulset-2            1/1     Running     0          87m
```
## Get password for UI
User name is 'admin', password can be retrieved by
```
kubectl -n contrail get secret cluster1-admin-password -ojson |jq .data.password | tr -d '"' |base64 --decode
```
## Cleanup
```
kubectl delete -k github.com/Juniper/contrail-operator/deploy/kustomize/contrail/${REPLICA}node/${RELEASE}
kubectl delete -k github.com/Juniper/contrail-operator/deploy/kustomize/operator/R2008
kubectl delete crds --all
kubectl delete pv --all
```
also, remove /mnt/zookeeper and /mnt/cassandra from the master nodes    
