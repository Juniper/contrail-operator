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

This operator implements custom resources for 
