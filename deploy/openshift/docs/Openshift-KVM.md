## Contrail with OpenShift 4.x installation on VMs running on KVM

The following procedure works also if bare metal servers are used. If there are existing DNS, DHCP, HTTP, PXE servers, update services following examples [here](bare-metal-prerequisites.md) and jump to [Create Ignition Configs](Openshift-KVM.md#create-ignition-configs).

The procedure follows [helper node installation guide line](https://github.com/RedHatOfficial/ocp4-helpernode/blob/master/docs/quickstart.md). Some modifications occurs when applying Contrail manifests

On the hypervisor host create a working directory

```
# mkdir ~/ocp4-workingdir
# cd ~/ocp4-workingdir
```

### Create a virtual network

Download the virtual network configuration file, virt-net.xml
```
# wget https://raw.githubusercontent.com/RedHatOfficial/ocp4-helpernode/master/docs/examples/virt-net.xml
```
Create a virtual network using this file provided in this repo (modify if you need to).

```
# virsh net-define --file virt-net.xml
```
Set it to autostart on boot

```
# virsh net-autostart openshift4
# virsh net-start openshift4
```

### Create a CentOS 7/8 VM

Download the Kickstart file for either EL 7 or EL 8 for the helper node.

**EL 7**
```
# wget https://raw.githubusercontent.com/RedHatOfficial/ocp4-helpernode/master/docs/examples/helper-ks.cfg -O helper-ks.cfg
```

**EL 8**
```
# wget https://raw.githubusercontent.com/RedHatOfficial/ocp4-helpernode/master/docs/examples/helper-ks8.cfg -O helper-ks.cfg
```
Edit `helper-ks.cfg` for your environment and use it to install the helper. The following command installs it "unattended".

**EL 7**
```
# virt-install --name="ocp4-aHelper" --vcpus=2 --ram=4096 \
--disk path=/var/lib/libvirt/images/ocp4-aHelper.qcow2,bus=virtio,size=30 \
--os-variant centos7.0 --network network=openshift4,model=virtio \
--boot hd,menu=on --location /var/lib/libvirt/iso/CentOS-7-x86_64-Minimal-2003.iso \
--initrd-inject helper-ks.cfg --extra-args "inst.ks=file:/helper-ks.cfg" --noautoconsole
```

**EL 8**
```
# virt-install --name="ocp4-aHelper" --vcpus=2 --ram=4096 \
--disk path=/var/lib/libvirt/images/ocp4-aHelper.qcow2,bus=virtio,size=50 \
--os-variant centos8 --network network=openshift4,model=virtio \
--boot hd,menu=on --location /var/lib/libvirt/iso/CentOS-8.2.2004-x86_64-dvd1.iso \
--initrd-inject helper-ks.cfg --extra-args "inst.ks=file:/helper-ks.cfg" --noautoconsole
```

The provided Kickstart file installs the helper with the following settings (which is based on the virt-net.xml file that was used before).

- HELPER_IP - 192.168.7.77
- NetMask - 255.255.255.0
- Default Gateway - 192.168.7.1
- DNS Server - 8.8.8.8

You can watch the progress by lauching the viewer
```
# virt-viewer --domain-name ocp4-aHelper
```

Once it's done, it'll shut off...turn it on with the following command
```
# virsh start ocp4-aHelper
```

### Prepare the Helper Node

After the helper node is installed; login to it
```
# ssh -l root <HELPER_IP>
```

Install EPEL and update. If kernel is updated, reboot.
```
# yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %rhel).noarch.rpm
# yum -y update
# reboot
```

Install `ansible` and `git` and clone this helpernode repo
```
# yum -y install ansible git
# git clone https://github.com/RedHatOfficial/ocp4-helpernode
# cd ocp4-helpernode
```

Copy vars.yaml file. Edit and change it if is necessary (domain name, mac addresses, ipam in case you used a different subnet...).
```
# cp docs/examples/vars.yaml .
```

To modify the OpenShift version modify `vars/main.yml` file

Run the playbook to setup your helper node
```
# ansible-playbook -e @vars.yaml tasks/main.yml
```

After it is done run the following command to get info about your environment and some install help
```
# /usr/local/bin/helpernodecheck services
```

### Create Ignition Configs

**Note**: make sure NTP server is configured on Hypervisor and Helper node, otherwise installation will fail with error `X509: certificate has expired or is not yet valid`_

Create a place to store your pull-secret
```
# mkdir -p ~/.openshift
```

Visit [try.openshift.com](https://cloud.redhat.com/openshift/install) and select "Bare Metal". Download your pull secret and save it under ~/.openshift/pull-secret

```
# ls -1 ~/.openshift/pull-secret
/root/.openshift/pull-secret
```

This playbook creates an sshkey for you; it's under `~/.ssh/helper_rsa`. You can use this key or create/user another one if you wish.

```
# ls -1 ~/.ssh/helper_rsa
/root/.ssh/helper_rsa
```

Create an install directory

```
# mkdir ~/ocp4
# cd ~/ocp4
```

Next, create an `install-config.yaml` file.

```
# cat <<EOF > install-config.yaml
apiVersion: v1
baseDomain: example.com
compute:
- hyperthreading: Disabled
  name: worker
  replicas: 0
controlPlane:
  hyperthreading: Disabled
  name: master
  replicas: 3
metadata:
  name: ocp4
networking:
  clusterNetworks:
  - cidr: 10.254.0.0/16
    hostPrefix: 24
  networkType: Contrail
  serviceNetwork:
  - 172.30.0.0/16
platform:
  none: {}
pullSecret: '$(< ~/.openshift/pull-secret)'
sshKey: '$(< ~/.ssh/helper_rsa.pub)'
EOF
```

Create the installation manifests
```
# openshift-install create manifests
```

Edit the `manifests/cluster-scheduler-02-config.yml` Kubernetes manifest file to prevent Pods from being scheduled on the control plane machines by setting `mastersSchedulable` to `false`.
```
# sed -i 's/mastersSchedulable: true/mastersSchedulable: false/g' manifests/cluster-scheduler-02-config.yml
```
It should look something like this after you edit it

```
# cat manifests/cluster-scheduler-02-config.yml
apiVersion: config.openshift.io/v1
kind: Scheduler
metadata:
  creationTimestamp: null
  name: cluster
spec:
  mastersSchedulable: false
  policy:
    name: ""
status: {}
```

Install Contrail manifests and configs

Download additional Contrail manifests and configs and add them to the generated manifests directory by executing these commands:
```
bash <<EOF
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-01-namespace.yaml -o manifests/00-contrail-01-namespace.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-02-admin-password.yaml -o manifests/00-contrail-02-admin-password.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-02-rbac-auth.yaml -o manifests/00-contrail-02-rbac-auth.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-02-registry-secret.yaml -o manifests/00-contrail-02-registry-secret.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-03-cluster-role.yaml -o manifests/00-contrail-03-cluster-role.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-04-serviceaccount.yaml -o manifests/00-contrail-04-serviceaccount.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-05-rolebinding.yaml -o manifests/00-contrail-05-rolebinding.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/00-contrail-06-clusterrolebinding.yaml -o manifests/00-contrail-06-clusterrolebinding.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_cassandras_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_cassandras_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_commands_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_commands_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_configs_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_configs_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_contrailcnis_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_contrailcnis_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_fernetkeymanagers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_fernetkeymanagers_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_contrailmonitors_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_contrailmonitors_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_contrailstatusmonitors_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_contrailstatusmonitors_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_controls_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_controls_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_keystones_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_keystones_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_kubemanagers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_kubemanagers_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_managers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_managers_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_memcacheds_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_memcacheds_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_postgres_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_postgres_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_provisionmanagers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_provisionmanagers_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_rabbitmqs_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_rabbitmqs_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_swiftproxies_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swiftproxies_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_swifts_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swifts_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_swiftstorages_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swiftstorages_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_vrouters_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_vrouters_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_webuis_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_webuis_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/crds/contrail.juniper.net_zookeepers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_zookeepers_crd.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/releases/R2011/manifests/00-contrail-08-operator.yaml -o manifests/00-contrail-08-operator.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/releases/R2011/manifests/00-contrail-09-manager.yaml -o manifests/00-contrail-09-manager.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/manifests/cluster-network-02-config.yml -o manifests/cluster-network-02-config.yml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_master-iptables-machine-config.yaml -o openshift/99_master-iptables-machine-config.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_master-kernel-modules-overlay.yaml -o openshift/99_master-kernel-modules-overlay.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_master_network_functions.yaml -o openshift/99_master_network_functions.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_master_network_manager_stop_service.yaml -o openshift/99_master_network_manager_stop_service.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_master-pv-mounts.yaml -o openshift/99_master-pv-mounts.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_worker-iptables-machine-config.yaml -o openshift/99_worker-iptables-machine-config.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_worker-kernel-modules-overlay.yaml -o openshift/99_worker-kernel-modules-overlay.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_worker_network_functions.yaml -o openshift/99_worker_network_functions.yaml;\
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2011/deploy/openshift/openshift/99_worker_network_manager_stop_service.yaml -o openshift/99_worker_network_manager_stop_service.yaml;
EOF
```

Modify `manifests/00-contrail-02-registry-secret.yaml` file providing proper configuration with credentials to *hub.juniper.net* registry.

**NOTE:** You may create base64 encoded value for config with script provided in [here](https://github.com/Juniper/contrail-operator/tree/master/deploy/openshift/tools/docker-config-generate) directory.
Copy output of the script and paste into contrail registry secret manifest

**NOTE**: If your environment has to use a specific NTP server, follow [these](./chrony-ntp-configuration.md) instructions before executing next steps.

Generate the ignition configs

```
# openshift-install create ignition-configs
```

Copy the ignition files in the `ignition` directory for the websever

```
# cp ~/ocp4/*.ign /var/www/html/ignition/
# restorecon -vR /var/www/html/
# restorecon -vR /var/lib/tftpboot/
# chmod o+r /var/www/html/ignition/*.ign
```

### Install VMs

From the hypervisor launch VMs using PXE booting. For BMS, boot the servers using PXE booting.

Launch Bootstrap VM
```
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:60:72:67 --name ocp4-bootstrap --ram=8192 --vcpus=4 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-bootstrap.qcow2,size=120 --vnc
```
This command will create a bootstrap node VM, will connect to PXE server (our Helper), assign the IP address from DHCP and download the RHCOS image from the HTTP server. At the end of the installation it will embed the ignition file. After the node is installed and rebooted we can connect to it from the Helper.

```
# ssh -i ~/.ssh/helper_rsa core@192.168.7.20
```

Use `journalctl -f` to see the logs

On the Bootstrap node a temporary etcd and bootkube is created. When these services are running

```
[core@bootstrap ~]$ sudo crictl ps
CONTAINER           IMAGE                                                                                                                    CREATED              STATE               NAME                             ATTEMPT             POD ID
33762f4a23d7d       976cc3323bd3394e613ff3d9ff02cd2ab55456063e08d6e275e81f71349d6399                                                         54 seconds ago       Running             manager                          3                   29aed2b586f33
ad6f2453d7a16       86694d2cdf8823ae48f13242bbd7a381eaab0218831ed53e9806b5e19608b1ed                                                         About a minute ago   Running             kube-apiserver-insecure-readyz   0                   4cd5138fb8fa7
3bbdf4176882f       quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:8db6b7ce80d002fb2687c6408d0efaae7cd908bb83b7b13ea512ad880747f02c   About a minute ago   Running             kube-scheduler                   0                   b3e7e6831100c
57ad52023300e       quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:8db6b7ce80d002fb2687c6408d0efaae7cd908bb83b7b13ea512ad880747f02c   About a minute ago   Running             kube-controller-manager          0                   596e248e26449
a1dbe7b8950da       quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:8db6b7ce80d002fb2687c6408d0efaae7cd908bb83b7b13ea512ad880747f02c   About a minute ago   Running             kube-apiserver                   0                   4cd5138fb8fa7
5aa7a59a06feb       quay.io/openshift-release-dev/ocp-release@sha256:2e4bbcf4dff0857bec6328c77d3a0480c1ae6778d48c7fba197f54a3e1912c72        About a minute ago   Running             cluster-version-operator         0                   3ab41a6177a8d
ca45790f4a5f6       099c2a95af4ff574a824a5476a960e86deec7e31882294116d195eb186752d36                                                         About a minute ago   Running             etcd-metrics                     0                   081b292dfe92b
e72fb8aaa1606       quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:3af018d2799385f4e516cec73e24915351acf0012a8b775cf852eb05ad34797d   About a minute ago   Running             etcd-member                      0                   081b292dfe92b
ca56bbf2708f7       1ac19399249cf839f48e246869b6932ce1273afb6a11a25e0eccb01092ea3cbf                                                         About a minute ago   Running             machine-config-server            0                   c1127810cd0ed
```

it is time to launch the Masters VMs

```
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:e7:9d:67 --name ocp4-master0 --ram=40960 --vcpus=8 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-master0.qcow2,size=250 --vnc
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:80:16:23 --name ocp4-master1 --ram=40960 --vcpus=8 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-master1.qcow2,size=250 --vnc
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:d5:1c:39 --name ocp4-master2 --ram=40960 --vcpus=8 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-master2.qcow2,size=250 --vnc
```

You can login to the Master from Helper Node

```
# ssh -i ~/.ssh/helper_rsa core@192.168.7.21
# ssh -i ~/.ssh/helper_rsa core@192.168.7.22
# ssh -i ~/.ssh/helper_rsa core@192.168.7.23
```

You can monitor pods creation with `sudo crictl ps` command

### Wait for install

The boostrap VM actually does the install for you; you can track it with the following command from the Helper. Make sure you are in `~/ocp4` directory

```
# openshift-install wait-for bootstrap-complete --log-level debug
```

Once you see this message below...
```
INFO Waiting up to 30m0s for the Kubernetes API at https://api.ocp4.example.com:6443...
INFO API v1.13.4+838b4fa up
INFO Waiting up to 30m0s for bootstrapping to complete...
DEBUG Bootstrap status: complete
INFO It is now safe to remove the bootstrap resources
```

You can delete the bootstrap VM and luanch the Worker nodes from the Hypervisor

```
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:f4:26:a1 --name ocp4-worker0 --ram=16384 --vcpus=6 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-worker0.qcow2,size=120 --vnc
# virt-install --pxe --network bridge=openshift4 --mac=52:54:00:82:90:00 --name ocp4-worker1 --ram=16384 --vcpus=6 --os-variant rhel8.0 --disk path=/var/lib/libvirt/images/ocp4-worker1.qcow2,size=120 --vnc
```

### Finish Install

Login to your cluster

```
# export KUBECONFIG=/root/ocp4/auth/kubeconfig
```

Your install may be waiting for worker nodes to get approved. Normally the machineconfig node approval operator takes care of this for you. However, sometimes this needs to be done manually. Check pending CSRs with the following command

```
# oc get csr
```

You can approve all pending CSRs in "one shot" with the following

```
# oc get csr -o go-template='{{range .items}}{{if not .status}}{{.metadata.name}}{{"\n"}}{{end}}{{end}}' | xargs oc adm certificate approve
```

You may have to run this multiple times depending on how many workers you have and in what order they come in. Keep a watch on these CSRs

```
# watch -n5 oc get csr
```

In order to setup your registry, you first have to set the `managementState` to `Managed` for your cluster
```
# oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"managementState":"Managed"}}'
```

For PoCs, using emptyDir is okay (to use PVs follow [this](https://docs.openshift.com/container-platform/latest/installing/installing_bare_metal/installing-bare-metal.html#registry-configuring-storage-baremetal_installing-bare-metal) doc)

```
# oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"emptyDir":{}}}}'
```

If you need to expose the registry, run this command
```
# oc patch configs.imageregistry.operator.openshift.io/cluster --type merge -p '{"spec":{"defaultRoute":true}}'
```

To finish the install process, run the following (make sure you are in `~/ocp4` directory)

```
openshift-install wait-for install-complete
```

When the follwoing message is displayed the installation has finished
```
# openshift-install wait-for install-complete
INFO Waiting up to 30m0s for the cluster at https://api.ocp4.example.com:6443 to initialize...
INFO Waiting up to 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/root/ocp4/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.ocp4.example.com
INFO Login to the console with user: kubeadmin, password: XXX-XXXX-XXXX-XXXX
```

### Adding a user

By default, OpenShift4 ships with a single kubeadmin user, that could be used during initial cluster configuration. You will create a Custom Resource (CR) to define a HTTPasswd identity provider.

To use the HTPasswd identity provider, you must generate a flat file that contains the user names and passwords for your cluster by using [htpasswd](https://httpd.apache.org/docs/2.4/programs/htpasswd.html).
```
$ htpasswd -c -B -b users.htpasswd testuser MyPassword
```

You should get a file called users.httpasswd

Next we need to define a secret that contains the HTPasswd user file
```
$ oc create secret generic htpass-secret --from-file=htpasswd=/root/ocp4/users.htpasswd -n openshift-config
```

This Custom Resource shows the parameters and acceptable values for an HTPasswd identity provider.

```
$ cat htpasswdCR.yaml
apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
spec:
  identityProviders:
  - name: testuser
    mappingMethod: claim
    type: HTPasswd
    htpasswd:
      fileData:
        name: htpass-secret
```

Apply the defined CR
```
$ oc create -f htpasswdCR.yaml
```

Add the user to `cluster-amdin` role
```
$ oc adm policy add-cluster-role-to-user cluster-admin testuser
```

Login using the user
```
oc login -u testuser
Authentication required for https://api.ocp4.example.com:6443 (openshift)
Username: testuser
Password:
Login successful.
```

Now it is safe to remove kubeadmin user. Details [here](https://docs.openshift.com/container-platform/4.5/authentication/remove-kubeadmin.html).

### Credits: [Ovidiu Valeanu](https://github.com/ovaleanujnpr)
