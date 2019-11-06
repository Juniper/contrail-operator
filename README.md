# Contrail Operator

## TOC

[TOC]

## Prerequisites

- An installed kubernetes cluster (>=1.15.0)    
- An account on Atom gitlab (https://ssd-git.juniper.net/atom/atom)

### Preparation

#### Get an Access Token

- go to https://ssd-git.juniper.net/profile/personal_access_tokens    
- enter a name    
- enter an expiration date    
- check all scopes     
- click on "Create personal access token"    
- store the generated token on a save place     

#### Install the Harbour container registry certificates on ALL kubernetes nodes (including masters)

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \ 
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/cacerts%2FAddTrustExternalCARoot%2Ecrt/raw?ref=master \
  -o /etc/pki/ca-trust/source/anchors/AddTrustExternalCARoot.crt
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/cacerts%2FSectigoRSADomainValidationSecureServerCA%2Ecrt/raw?ref=master \
  -o /etc/pki/ca-trust/source/anchors/SectigoRSADomainValidationSecureServerCA.crt
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/cacerts%2FUSERTrustRSAAddTrustCA%2Ecrt/raw?ref=master \
  -o /etc/pki/ca-trust/source/anchors/USERTrustRSAAddTrustCA.crt
update-ca-trust
/bin/systemctl restart docker.service
```

Again, this has to happen on ALL k8s nodes!    

#### Create an image pull secret

```
kubectl create secret docker-registry atom-registry \
  --docker-server=atom-registry-prod.juniper.net/atom \
  --docker-username $YOUR_LDAP_USER_WITHOUT_DOMAIN \
  --docker-password '$YOUR_LDAP_PASSWORD'
```

### Quick Start Installation

#### Set branch

```
BRANCH=master
```

#### Create a Image Pull Secret for contrail/contrail-nightly registry

```
kubectl create secret docker-registry contrail-nightly \
  --docker-server=hub.juniper.net/contrail-nightly \
  --docker-username=$YOUR_REGISTRY_USER \
  --docker-password=$YOUR_REGISTRY_PASSWORD
```

More configuration not needed for quick start

#### Create CRDs, Service Account, Role, Bindings and Operator

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/contrail%2Foperator%2Fdeploy%2F1-create-operator%2Eyaml/raw?ref=${BRANCH} | kubectl apply -f -
```

Wait for Contrail Operator deployment to run:    

```
[root@kvm1 ~]# kubectl get pods
NAME                                 READY   STATUS    RESTARTS   AGE
contrail-operator-7bbb99845c-qktvf   1/1     Running   0          16m
```

#### Apply a 1 or a 3 node manifest

Note: THIS WILL INSTALL CONTRAIL USING DEFAULTS!

##### 1 Node

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/contrail%2Foperator%2Fdeploy%2F2-start-operator-1node%2Eyaml/raw?ref=${BRANCH}| kubectl apply -f -
```

#### 3 Node

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/contrail%2Foperator%2Fdeploy%2F2-start-operator-3node-minimal%2Eyaml/raw?ref=${BRANCH}| kubectl apply -f -
```

### Deploy a custom configuration

#### Get a manifest

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/contrail%2Foperator%2Fdeploy%2F2-start-operator-3node%2Eyaml/raw?ref=${BRANCH} \
  -o 2-start-operator-3node-custom.yaml
```

### Deploy a custom configuration

#### Get a manifest

```
curl --header "PRIVATE-TOKEN: $YOUR_ACCESS_TOKEN" \
  https://svl-ssd-git.juniper.net/api/v4/projects/3584/repository/files/contrail%2Foperator%2Fdeploy%2F2-start-operator-3node%2Eyaml/raw?ref=${BRANCH} \
  -o 2-start-operator-3node-custom.yaml
```

### Edit manifest

```
vi 2-start-operator-3node-custom.yaml

---
apiVersion: contrail.juniper.net/v1alpha1
kind: Manager
metadata:
  ## (Mandatory) defines the name of the cluster
  name: cluster1
spec:
  commonConfiguration:
    ## (Optional - defaults to 1). Defines the number of instances globally.
    ## Can be overwritten in each service
    replicas: 3
    ## (Optional - defaults to true). DO NOT CHANGE FOR NOW!
    hostNetwork: true
    ## (Optional). Needed if images are pulled from a password potected registry.
    imagePullSecrets:
    - contrail-nightly
  services:
    ## Cassandra Services
    cassandras:
    - metadata:
        name: cassandra1
        labels:
          ## (Manadatory). Has to match managers metadata.name
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          ## (Optional). Selects the node to run the service on.
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          images:
            cassandra: cassandra:3.11.4
            init: python:alpine
          ## (Optional). Cassandra service configuration
          listenAddress: auto
          startRpc: true
          port: 9160
          cqlPort: 9042
          sslStoragePort: 7001
          storagePort: 7000
          jmxLocalPort: 7199
          clusterName: ContrailConfigDB
          maxHeapSize: 1024M
          minHeapSize: 100M
    zookeepers:
    - metadata:
        name: zookeeper1
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          tolerations:
          - effect: NoSchedule
            operator: Exists
          - effect: NoExecute
            operator: Exists
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          images:
            zookeeper: docker.io/zookeeper:3.5.5
            init: python:alpine
          clientPort: 2181
          electionPort: 3888
          serverPort: 2888
    rabbitmq:
      metadata:
        name: rabbitmq1
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          images:
            rabbitmq: rabbitmq:3.7
            init: python:alpine
          port: 5673
          erlangCookie: 47EFF3BB-4786-46E0-A5BB-58455B3C2CB4
    config:
      metadata:
        name: config1
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          ## (Manadatory). Defines which cassandra/zookeeper instance to use.
          ## Has to match cassandras and zookeepers metadata.name
          cassandraInstance: cassandra1
          zookeeperInstance: zookeeper1
          images:
            api: hub.juniper.net/contrail-nightly/contrail-controller-config-api:1908.47
            devicemanager: hub.juniper.net/contrail-nightly/contrail-controller-config-devicemgr:1908.47
            schematransformer: hub.juniper.net/contrail-nightly/contrail-controller-config-schema:1908.47
            servicemonitor: hub.juniper.net/contrail-nightly/contrail-controller-config-svcmonitor:1908.47
            analyticsapi: hub.juniper.net/contrail-nightly/contrail-analytics-api:1908.47
            collector: hub.juniper.net/contrail-nightly/contrail-analytics-collector:1908.47
            redis: redis:4.0.2
            nodemanagerconfig: hub.juniper.net/contrail-nightly/contrail-nodemgr:1908.47
            nodemanageranalytics: hub.juniper.net/contrail-nightly/contrail-nodemgr:1908.47
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
            init: python:alpine
    kubemanagers:
    - metadata:
        name: kubemanager1
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          cassandraInstance: cassandra1
          zookeeperInstance: zookeeper1
          images:
            kubemanager: hub.juniper.net/contrail-nightly/contrail-kubernetes-kube-manager:1908.47
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
            init: python:alpine
          ## (Optional). If kubeadm has the required information useKubeadmConfig can be set to true.
          ## In that case kubernetesAPIServer, kubernetesAPIPort, podSubnet, serviceSubnet, kubernetesClusterName
          ## don't need to be provided. Can be checked with kubeadm.
          useKubeadmConfig: false
          serviceAccount: contrail-service-account
          clusterRole: contrail-cluster-role
          clusterRoleBinding: contrail-cluster-role-binding
          cloudOrchestrator: kubernetes
          kubernetesAPIServer: "10.96.0.1"
          kubernetesAPIPort: 443
          podSubnet: 10.32.0.0/12
          serviceSubnet: 10.96.0.0/12
          kubernetesClusterName: kubernetes
          ipFabricForwarding: true
          ipFabricSnat: true
          kubernetesTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    controls:
    - metadata:
        name: control1
        labels:
          contrail_cluster: cluster1
          control_role: master
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          cassandraInstance: cassandra1
          zookeeperInstance: zookeeper1
          images:
            control: hub.juniper.net/contrail-nightly/contrail-controller-control-control:1908.47
            dns: hub.juniper.net/contrail-nightly/contrail-controller-control-dns:1908.47
            named: hub.juniper.net/contrail-nightly/contrail-controller-control-named:1908.47
            nodemanager: hub.juniper.net/contrail-nightly/contrail-nodemgr:1908.47
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
            init: python:alpine
    webui:
      metadata:
        name: webui1
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""   
        serviceConfiguration:
          cassandraInstance: cassandra1
          images:
            webuiweb: hub.juniper.net/contrail-nightly/contrail-controller-webui-web:1908.47
            webuijob: hub.juniper.net/contrail-nightly/contrail-controller-webui-job:1908.47
            redis: redis:4.0.2
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
    ## There can be multiple vrouter daemonsets with different configurations.
    ## This example has two DS', one for the master and one for the nodes.
    ## In this example, nodes are be identified by the nodeselector and the
    ## label node-role.opencontrail.org: "vrouter". Hence, nodes must be labeled first
    vrouters:
    - metadata:
        name: vroutermaster
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: ""    
        serviceConfiguration:
          cassandraInstance: cassandra1
          controlInstance: control1
          images:
            vrouteragent: hub.juniper.net/contrail-nightly/contrail-vrouter-agent:1908.47
            vrouterkernelinit: hub.juniper.net/contrail-nightly/contrail-vrouter-kernel-init:1908.47
            vroutercni: hub.juniper.net/contrail-nightly/contrail-kubernetes-cni-init:1908.47
            nodemanager: hub.juniper.net/contrail-nightly/contrail-nodemgr:1908.47
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
            init: python:alpine
    - metadata:
        name: vrouternodes
        labels:
          contrail_cluster: cluster1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.opencontrail.org: "vrouter"  
        serviceConfiguration:
          cassandraInstance: cassandra1
          controlInstance: control1
          images:
            vrouteragent: hub.juniper.net/contrail-nightly/contrail-vrouter-agent:1908.47
            vrouterkernelinit: hub.juniper.net/contrail-nightly/contrail-vrouter-kernel-init:1908.47
            vroutercni: hub.juniper.net/contrail-nightly/contrail-kubernetes-cni-init:1908.47
            nodemanager: hub.juniper.net/contrail-nightly/contrail-nodemgr:1908.47
            nodeinit: hub.juniper.net/contrail-nightly/contrail-node-init:1908.47
            init: python:alpine 
```

#### Apply custom manifest

```
kubectl apply -f 2-start-operator-3node-custom.yaml
```