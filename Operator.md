# Contrail Operator 101

## What is an operator?

Operator is an automated software extension which allow to easily manage applications and its' components.
It's used to package, deploy and manage applications based on Kubernetes custom resources.

Operator acts as a controller which extend Kubernetes API to manage lifecycle of dependent application resources. Controller implements logic which periodically compare desired state of a cluster to it's actual state and apply corrections to meet declared state.

Compared to other methods of managing application deployment on Kubernetes clusters, operator allows to cover completly lifetime of application.
Below graph compares operator to Helm and Ansible.
Helm provides only installation process with some mechanisms of upgrade.
Ansible additionaly allow to cover some aspects of managing application resources lifecycle.
However, only operator allows to have full insights into created cluster resources and perform custom operations on them during lifecycle.

<< HELM_ANSIBLE_OPERATOR >>

Custom operator may be written in multiple languages as in fact it's just a piece of code that periodically act with requests to Kubernetes API.
Currently, the most popular language to write operators is the language that was used to write Kubernetes project itself - Go.
Helpful tool for creating operators is [Operator Framework](https://github.com/operator-framework) which distributes [Operator-SDK](https://github.com/operator-framework/operator-sdk) commonly used to create operators - for example this operator.

Every custom resource built with operator contains 2 elements.
First element is API which defines how resource is defined and what's it's structure.
This definition is used afterwards by Operator Framework to generate CRD manifests applied on a cluster and defines how user should write manifests to deploy custom resource successfully.
Second element is controller which runs on operator pod and handles logic of custom resource.
Every controller has Reconcile methos which is run periodically and contain code which defines what to do every loop.
Commonly, it creates, deletes or updates standard Kubernetes resources like pods, sets, secrets etc.

## How does it work here?

This operator implements custom resources for Contrail deployment.
Before contrail-operator there was only project contrail-ansible-deployer which contained
set of ansible playbooks that configured instances based on file with 

```
$ kubectl get po -n contrail
NAME                                          READY   STATUS             RESTARTS   AGE
cassandra1-cassandra-statefulset-0            0/1     Running            0          39m
cassandra1-cassandra-statefulset-1            0/1     Running            0          39m
cassandra1-cassandra-statefulset-2            0/1     Running            0          39m
config1-config-statefulset-0                  10/10   Running            0          38m
config1-config-statefulset-1                  10/10   Running            0          38m
config1-config-statefulset-2                  0/10    Running            0          39m
contrail-operator-dd5bb5c-klqwb               1/1     Running            0          42m
control1-control-statefulset-0                4/4     Running            0          30m
control1-control-statefulset-1                4/4     Running            0          30m
control1-control-statefulset-2                4/4     Running            0          30m
kubemanager1-kubemanager-statefulset-0        2/2     Running            0          30m
kubemanager1-kubemanager-statefulset-1        2/2     Running            0          30m
kubemanager1-kubemanager-statefulset-2        2/2     Running            0          30m
provmanager1-provisionmanager-statefulset-0   1/1     Running            0          30m
rabbitmq1-rabbitmq-statefulset-0              1/1     Running            0          39m
rabbitmq1-rabbitmq-statefulset-1              1/1     Running            0          39m
rabbitmq1-rabbitmq-statefulset-2              1/1     Running            0          39m
vroutermasternodes-vrouter-daemonset-rgl4t    1/1     Running            0          28m
vroutermasternodes-vrouter-daemonset-ttc7c    1/1     Running            0          28m
vroutermasternodes-vrouter-daemonset-wn6qg    1/1     Running            0          28m
vrouterworkernodes-vrouter-daemonset-gs4bw    1/1     Running            0          5m
vrouterworkernodes-vrouter-daemonset-p7zkw    1/1     Running            0          5m
vrouterworkernodes-vrouter-daemonset-pqfw9    1/1     Running            0          5m
webui1-webui-statefulset-0                    3/3     Running            0          30m
webui1-webui-statefulset-1                    2/3     Running            0          30m
webui1-webui-statefulset-2                    2/3     Running            0          30m
zookeeper1-zookeeper-statefulset-0            1/1     Running            0          8m
zookeeper1-zookeeper-statefulset-1            1/1     Running            0          8m
zookeeper1-zookeeper-statefulset-2            1/1     Running            0          8m

```