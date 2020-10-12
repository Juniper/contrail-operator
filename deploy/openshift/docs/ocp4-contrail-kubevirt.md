## OpenShift Virtualisation (KubeVirt) on OpenShift 4.5 with Contrail cluster

This tutorial walks you through steps to install OpenShift Virtualisation 2.4 on OpenShift 4.5 with Contrail.

### Prerequisites
A Red Hat OpenShift 4.5 with Contrail 2008 or later, running on BMS or VMs supporting nested virtualisation. For installation procedure please check [here](https://github.com/ovaleanujnpr/openshift4.x/blob/master/docs/ocp4-contrail-vm-bms.md).

```
$ oc get pods -n contrail
NAME                                          READY   STATUS      RESTARTS   AGE
cassandra1-cassandra-statefulset-0            1/1     Running     0          27h
cassandra1-cassandra-statefulset-1            1/1     Running     0          27h
cassandra1-cassandra-statefulset-2            1/1     Running     0          27h
cnimasternodes-contrailcni-job-gjxjv          0/1     Completed   0          27h
cnimasternodes-contrailcni-job-r9z4v          0/1     Completed   0          27h
cnimasternodes-contrailcni-job-zgvz7          0/1     Completed   0          27h
cniworkernodes-contrailcni-job-mq9rb          0/1     Completed   0          27h
cniworkernodes-contrailcni-job-sr7v8          0/1     Completed   0          27h
config1-config-statefulset-0                  10/10   Running     1          27h
config1-config-statefulset-1                  10/10   Running     0          27h
config1-config-statefulset-2                  10/10   Running     2          27h
contrail-operator-76488c8cb9-7k5lc            1/1     Running     0          27h
contrail-operator-76488c8cb9-fgllk            1/1     Running     1          27h
contrail-operator-76488c8cb9-fmp4c            1/1     Running     0          27h
control1-control-statefulset-0                4/4     Running     0          27h
control1-control-statefulset-1                4/4     Running     0          27h
control1-control-statefulset-2                4/4     Running     0          27h
kubemanager1-kubemanager-statefulset-0        2/2     Running     1          27h
kubemanager1-kubemanager-statefulset-1        2/2     Running     0          27h
kubemanager1-kubemanager-statefulset-2        2/2     Running     1          27h
provmanager1-provisionmanager-statefulset-0   1/1     Running     1          27h
rabbitmq1-rabbitmq-statefulset-0              1/1     Running     0          27h
rabbitmq1-rabbitmq-statefulset-1              1/1     Running     0          27h
rabbitmq1-rabbitmq-statefulset-2              1/1     Running     0          27h
vroutermasternodes-vrouter-daemonset-7kxmx    1/1     Running     0          27h
vroutermasternodes-vrouter-daemonset-9l4wt    1/1     Running     0          27h
vroutermasternodes-vrouter-daemonset-xlzr7    1/1     Running     0          27h
vrouterworkernodes-vrouter-daemonset-qjgrs    1/1     Running     0          27h
vrouterworkernodes-vrouter-daemonset-zm6ml    1/1     Running     0          27h
webui1-webui-statefulset-0                    3/3     Running     0          27h
webui1-webui-statefulset-1                    3/3     Running     0          27h
webui1-webui-statefulset-2                    3/3     Running     0          27h
zookeeper1-zookeeper-statefulset-0            1/1     Running     0          27h
zookeeper1-zookeeper-statefulset-1            1/1     Running     0          27h
zookeeper1-zookeeper-statefulset-2            1/1     Running     0          27h
```


### Installing OpenShift Virtualisation operator

Login as a user with `cluster-admin` privileges

Create a yaml file containing the following manifest

```
$ cat <<EOF > cnv.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: openshift-cnv
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: kubevirt-hyperconverged-group
  namespace: openshift-cnv
spec:
  targetNamespaces:
    - openshift-cnv
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: hco-operatorhub
  namespace: openshift-cnv
spec:
  source: redhat-operators
  sourceNamespace: openshift-marketplace
  name: kubevirt-hyperconverged
  startingCSV: kubevirt-hyperconverged-operator.v2.4.1
  channel: "2.4"
EOF
```
This yaml file will create the required `Namespace`, `OperatorGroup` and `Subscription` for the OpenShift Virtualisation.
```
$ oc apply -f cmv.yaml
```

Now we can deploy the OpenShift Virtualisation operator. Create a yaml file with the following content
```
$ cat <<EOF > kubevirt-hyperconverged.yaml
apiVersion: hco.kubevirt.io/v1alpha1
kind: HyperConverged
metadata:
  name: kubevirt-hyperconverged
  namespace: openshift-cnv
spec:
  BareMetalPlatform: true
EOF
```
```
$ oc apply -f kubevirt-hyperconverged.yaml
```

Check if all the pods are Running in `openshift-cnv` namespace
```
$ oc get pods -n openshift-cnv
NAME                                                  READY   STATUS    RESTARTS   AGE
bridge-marker-5tndk                                   1/1     Running   0          22h
bridge-marker-d2gff                                   1/1     Running   0          22h
bridge-marker-d8cgd                                   1/1     Running   0          22h
bridge-marker-r6glh                                   1/1     Running   0          22h
bridge-marker-rt5lb                                   1/1     Running   0          22h
cdi-apiserver-7c4566c98c-z89qz                        1/1     Running   0          22h
cdi-deployment-79fdcfdccb-xmphs                       1/1     Running   0          22h
cdi-operator-7785b655bb-7q5k6                         1/1     Running   0          22h
cdi-uploadproxy-5d4cc54b4c-g2ztz                      1/1     Running   0          22h
cluster-network-addons-operator-67d7f76cbd-8kl6l      1/1     Running   0          22h
hco-operator-854f5988c8-v2qbm                         1/1     Running   0          22h
hostpath-provisioner-operator-595b955c9d-zxngg        1/1     Running   0          22h
kube-cni-linux-bridge-plugin-5w67f                    1/1     Running   0          22h
kube-cni-linux-bridge-plugin-kjm8b                    1/1     Running   0          22h
kube-cni-linux-bridge-plugin-rgrn8                    1/1     Running   0          22h
kube-cni-linux-bridge-plugin-s6xkz                    1/1     Running   0          22h
kube-cni-linux-bridge-plugin-ssw29                    1/1     Running   0          22h
kubemacpool-mac-controller-manager-6f9c447bbd-phd5n   1/1     Running   0          22h
kubevirt-node-labeller-297nh                          1/1     Running   0          22h
kubevirt-node-labeller-cbjnl                          1/1     Running   0          22h
kubevirt-ssp-operator-75d54556b9-zq2kb                1/1     Running   0          22h
nmstate-handler-9prp8                                 1/1     Running   1          22h
nmstate-handler-dk4ht                                 1/1     Running   0          22h
nmstate-handler-fzjmk                                 1/1     Running   0          22h
nmstate-handler-rqwmq                                 1/1     Running   1          22h
nmstate-handler-spx7w                                 1/1     Running   0          22h
node-maintenance-operator-6486bcbfcd-rhn4l            1/1     Running   0          22h
ovs-cni-amd64-4t9ld                                   1/1     Running   0          22h
ovs-cni-amd64-5mdmq                                   1/1     Running   0          22h
ovs-cni-amd64-bz5d9                                   1/1     Running   0          22h
ovs-cni-amd64-h9j6j                                   1/1     Running   0          22h
ovs-cni-amd64-k8hwf                                   1/1     Running   0          22h
virt-api-7686f978db-ngwn2                             1/1     Running   0          22h
virt-api-7686f978db-nkl4d                             1/1     Running   0          22h
virt-controller-7d567db8c6-bbdjk                      1/1     Running   0          22h
virt-controller-7d567db8c6-n2vgk                      1/1     Running   0          22h
virt-handler-lkpsq                                    1/1     Running   0          5h30m
virt-handler-vfcbd                                    1/1     Running   0          5h30m
virt-operator-7995d994c4-9bxw9                        1/1     Running   0          22h
virt-operator-7995d994c4-q8wnv                        1/1     Running   0          22h
virt-template-validator-5d9bbfbcc7-g2zph              1/1     Running   0          22h
virt-template-validator-5d9bbfbcc7-lhhrw              1/1     Running   0          22h
vm-import-controller-58469cdfcf-kwkgb                 1/1     Running   0          22h
vm-import-operator-9495bd74c-dkw2h                    1/1     Running   0          22h
```
and `PHASE` of the ClusterServiceVersion (CSV) is Succeded.

```
$ oc get csv -n openshift-cnv
NAME                                      DISPLAY                    VERSION   REPLACES   PHASE
kubevirt-hyperconverged-operator.v2.4.1   OpenShift Virtualization   2.4.1                Succeeded
```

If you are running OpenShift Virtualisation in a nested enviroment, kubevirt-config ConfigMap must be updated to support [software emulation](https://github.com/kubevirt/kubevirt/blob/master/docs/software-emulation.md#software-emulation).

Add to kubevirt-config ConfigMap
```
data:
  debug.useEmulation: "true"
```
```
$ oc edit cm kubevirt-config -n openshift-cnv

apiVersion: v1
kind: ConfigMap
metadata:
  name: kubevirt-config
  namespace: openshift-cnv
data:
  debug.useEmulation: "true"
```
Then you need to restart `virt-handler` pods

### Creating Virtual Machines on OpenShift Virtualisation

Create namespace for the demo. I will call it `cnv-demo`.

```
$ oc new-project cnv-demo
```

Using Virtual Machine Instance (VMI) custom resources you can create VMs fully integrated in OpenShift.

Create a Virtual Machine with Centos 7 using the following manifest:

```
cat <<EOF > kubevirt-centos.yaml
apiVersion: kubevirt.io/v1alpha3
kind: VirtualMachineInstance
metadata:
  labels:
    special: vmi-centos7
  name: vmi-centos7
  namespace: cnv-demo
spec:
  domain:
    devices:
      disks:
      - disk:
          bus: virtio
        name: containerdisk
      - disk:
          bus: virtio
        name: cloudinitdisk
      interfaces:
      - name: default
        bridge: {}
    resources:
      requests:
        memory: 1024M
  networks:
  - name: default
    pod: {}
  volumes:
  - containerDisk:
      image: ovaleanu/centos:latest
    name: containerdisk
  - cloudInitNoCloud:
      userData: |-
        #cloud-config
        password: centos
        ssh_pwauth: True
        chpasswd: { expire: False }
    name: cloudinitdisk
EOF

$ oc apply -f kubevirt-centos.yaml
virtualmachineinstance.kubevirt.io/vmi-centos7 created
```

Create also a Virtual Machine with Fedora using the following manifest:

```
cat <<EOF > kubevirt-fedora.yaml
apiVersion: kubevirt.io/v1alpha3
kind: VirtualMachineInstance
metadata:
  labels:
    special: vmi-fedora
  name: vmi-fedora
spec:
  domain:
    devices:
      disks:
      - disk:
          bus: virtio
        name: containerdisk
      - disk:
          bus: virtio
        name: cloudinitdisk
      interfaces:
      - name: default
        bridge: {}
    resources:
      requests:
        memory: 1024M
  networks:
  - name: default
    pod: {}
  volumes:
  - containerDisk:
      image: kubevirt/fedora-cloud-registry-disk-demo
    name: containerdisk
  - cloudInitNoCloud:
      userData: |-
        #cloud-config
        password: fedora
        ssh_pwauth: True
        chpasswd: { expire: False }
    name: cloudinitdisk
EOF

$ oc apply -f kubevirt-fedora.yaml
virtualmachineinstance.kubevirt.io/vmi-fedora created
```

Check if the pods and VirtualMachineInstance were created
```
$ oc get pods -n cnv-demo
NAME                              READY   STATUS    RESTARTS   AGE
virt-launcher-vmi-centos7-q2jr6   2/2     Running   0          47s
virt-launcher-vmi-fedora-fml6q    2/2     Running   0          39s

$ oc get vmi
NAME          AGE    PHASE     IP              NODENAME
vmi-centos7   2m2s   Running   10.254.255.85   worker1.ocp4.example.com
vmi-fedora    114s   Running   10.254.255.67   worker1.ocp4.example.com
```

Create a service for each VM to connect with ssh through NodePort using node ip

```
cat <<EOF > kubevirt-centos-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: vmi-centos-ssh-svc
  namespace: cnv-demo
spec:
  ports:
  - name: centos-ssh-svc
    nodePort: 30000
    port: 27017
    protocol: TCP
    targetPort: 22
  selector:
    special: vmi-centos7
  type: NodePort
EOF

$ oc apply -f kubevirt-centos-svc.yaml
service/vmi-centos-ssh-svc created
```

```
cat <<EOF > kubevirt-fedora-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: vmi-fedora-ssh-svc
  namespace: cnv-demo
spec:
  ports:
  - name: fedora-ssh-svc
    nodePort: 31000
    port: 25025
    protocol: TCP
    targetPort: 22
  selector:
    special: vmi-fedora
  type: NodePort
EOF

$ oc apply -f kubevirt-fedora-svc.yaml
service/vmi-fedora-ssh-svc created
```
```
$ oc get svc -n cnv-demo
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
vmi-centos-ssh-svc   NodePort   172.30.18.44    <none>        27017:30000/TCP   31s
vmi-fedora-ssh-svc   NodePort   172.30.55.247   <none>        25025:31000/TCP   14s
```

### Test Virtual Machines connectivity

Connect to VMs with ssh via service NodePort using worker node IP address

```
$ ssh centos@192.168.7.12 -p 30000
The authenticity of host '[192.168.7.12]:30000 ([192.168.7.12]:30000)' can't be established.
ECDSA key fingerprint is SHA256:kk+9dbMqzpXDoPucnxiYozBgDt75IBSNS8Y4hUcEEmI.
ECDSA key fingerprint is MD5:86:b6:e9:3b:f0:55:ee:e7:fd:56:96:c3:4a:c6:fd:e0.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[192.168.7.12]:30000' (ECDSA) to the list of known hosts.
centos@192.168.7.12's password:

[centos@vmi-centos7 ~]$ uname -sr
Linux 3.10.0-957.12.2.el7.x86_64
```
Confirm the VM has access to outside world
```
[centos@vmi-centos7 ~]$ ping www.google.com
PING www.google.com (142.250.73.196) 56(84) bytes of data.
64 bytes from iad23s87-in-f4.1e100.net (142.250.73.196): icmp_seq=1 ttl=108 time=13.1 ms
64 bytes from iad23s87-in-f4.1e100.net (142.250.73.196): icmp_seq=2 ttl=108 time=11.9 ms
^C
--- www.google.com ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1003ms
rtt min/avg/max/mdev = 11.990/12.547/13.104/0.557 ms
```

Confirm the Centos VM can ping the other Fedora VM
```
[centos@vmi-centos7 ~]$ ping 10.254.255.67
PING 10.254.255.67 (10.254.255.67) 56(84) bytes of data.
64 bytes from 10.254.255.67: icmp_seq=1 ttl=63 time=8.33 ms
64 bytes from 10.254.255.67: icmp_seq=2 ttl=63 time=3.19 ms
^C
--- 10.254.255.67 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1002ms
rtt min/avg/max/mdev = 3.190/5.760/8.331/2.571 ms
```

Repeat the same steps for Fedora VM

```
$ ssh fedora@192.168.7.12 -p 31000
The authenticity of host '[192.168.7.12]:31000 ([192.168.7.12]:31000)' can't be established.
ECDSA key fingerprint is SHA256:JlhysyH0XiHXszLLqu8GmuSHB4msOYWPAJjZhv5j3FM.
ECDSA key fingerprint is MD5:62:ca:0b:b9:21:c9:2b:73:db:b6:09:e2:b0:b4:81:60.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[192.168.7.12]:31000' (ECDSA) to the list of known hosts.
fedora@192.168.7.12's password:

[fedora@vmi-fedora ~]$ uname -sr
Linux 4.13.9-300.fc27.x86_64

[fedora@vmi-fedora ~]$ ping www.google.com
PING www.google.com (142.250.73.196) 56(84) bytes of data.
64 bytes from iad23s87-in-f4.1e100.net (142.250.73.196): icmp_seq=1 ttl=108 time=14.3 ms
64 bytes from iad23s87-in-f4.1e100.net (142.250.73.196): icmp_seq=2 ttl=108 time=12.3 ms
^C
--- www.google.com ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1004ms
rtt min/avg/max/mdev = 12.326/13.360/14.394/1.034 ms
```

### Contrail Security policy

Create a network security policy to isolate the VM in its namespace. You allow only ssh for ingress and egress only to podNetwork.

```
cat <<EOF > kubevirt-centos-netpol.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
 name: netpol
 namespace: cnv-demo
spec:
 podSelector:
   matchLabels:
    special: vmi-centos7
 policyTypes:
 - Ingress
 - Egress
 ingress:
 - from:
   ports:
   - port: 22
 egress:
 - to:
   - ipBlock:
       cidr: 10.254.255.0/16
EOF

$ oc apply -f kubevirt-centos-netpol.yaml
networkpolicy.networking.k8s.io/netpol
```

Connect to VM again and ping Fedora VM ip. Pinging www.google.com will not work.
```
$ ssh centos@192.168.7.12 -p 30000
centos@192.168.7.12's password:
Last login: Wed Sep 30 11:09:49 2020 from 192.168.7.12
[centos@vmi-centos7 ~]$ ping 10.254.255.67
PING 10.254.255.67 (10.254.255.67) 56(84) bytes of data.
64 bytes from 10.254.255.67: icmp_seq=1 ttl=63 time=8.52 ms
64 bytes from 10.254.255.67: icmp_seq=2 ttl=63 time=4.12 ms
^C
--- 10.254.255.67 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1007ms
rtt min/avg/max/mdev = 4.120/6.320/8.521/2.201 ms
[centos@vmi-centos7 ~]$ ping www.google.com
^C
```

### Credits: [Ovidiu Valeanu](https://github.com/ovaleanujnpr)
