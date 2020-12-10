# E2e test executor

This directory contains ansible scripts that can be used to configure node for e2e tests execution. It installs docker and sets up proxy registry. After configuration node is ready to be added to jenkins.

The working directory is /var/lib/jenkins

## Dependencies

1. Ansible
1. Role geerlingguy.docker from ansible-galaxy
1. CentOS 7 with at least 20 vCPU, 48 GB RAM and 96GB of storage

## Usage

Adjust _inventory_ file and then execute ansible:

    ansible-playbook -i inventory configure.yaml
