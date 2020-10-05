# Prerequisities
Deployment depends strongly on Openshift installation which is described in this [documentation](https://docs.openshift.com/container-platform/4.1/installing/installing_aws/installing-aws-customizations.html)

Prerequisities that have to be fulfilled in order to dpeloy Contrail with operator on Openshift:
* openshift-install binary (>=4.4.8) ([download](https://cloud.redhat.com/openshift/install))
* Openshift pull secrets ([download](https://cloud.redhat.com/openshift/install/pull-secret))
* Configured AWS account with proper permissions and resolvable base domain configured in Route53 ([documentation](https://docs.openshift.com/container-platform/4.3/installing/installing_aws/installing-aws-account.html#installing-aws-account))
* Any SSH key generated on local machine to provide during installation
* (Optional) `oc` command line tool downloaded ([download](https://cloud.redhat.com/openshift/install))

# Deployment

1. Create install config with:
```
./openshift-install create install-config --dir <name of desired directory>
```
In created YAML file under specified directory setup all settings of cluster
under *networking* section change *networkType* field to *Contrail* (instead of *OpenshiftSDN*)

**NOTE**: Master nodes need larger instances.
For example, If you run cluster on AWS, use e.g. *m5.2xlarge*.

**NOTE**: See *install-config.example* for an example cluster configuration

2. Create manifests with
```
./openshift-install create manifests --dir <name of desired directory>
```

3. Install Contrail manifests and configs

Change directory to openshift install directory and download additional Contrail manifests and configs and add them to the generated manifests directory:
```
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-01-namespace.yaml -o manifests/00-contrail-01-namespace.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-02-admin-password.yaml -o manifests/00-contrail-02-admin-password.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-02-rbac-auth.yaml -o manifests/00-contrail-02-rbac-auth.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-02-registry-secret.yaml -o manifests/00-contrail-02-registry-secret.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-03-cluster-role.yaml -o manifests/00-contrail-03-cluster-role.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-04-serviceaccount.yaml -o manifests/00-contrail-04-serviceaccount.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-05-rolebinding.yaml -o manifests/00-contrail-05-rolebinding.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/0000000-contrail-06-clusterrolebinding.yaml -o manifests/00-contrail-06-clusterrolebinding.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_cassandras_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_cassandras_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_commands_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_commands_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_configs_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_configs_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_contrailmonitors_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_contrailmonitors_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_contrailstatusmonitors_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_contrailstatusmonitors_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_controls_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_controls_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_keystones_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_keystones_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_kubemanagers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_kubemanagers_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_managers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_managers_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_memcacheds_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_memcacheds_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_postgres_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_postgres_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_provisionmanagers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_provisionmanagers_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_rabbitmqs_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_rabbitmqs_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_swiftproxies_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swiftproxies_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_swifts_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swifts_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_swiftstorages_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_swiftstorages_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_vrouters_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_vrouters_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_webuis_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_webuis_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/crds/contrail.juniper.net_zookeepers_crd.yaml -o manifests/00-contrail-07-contrail.juniper.net_zookeepers_crd.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/master/deploy/openshift/releases/R2008/manifests/00-contrail-08-operator.yaml -o manifests/00-contrail-08-operator.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/master/deploy/openshift/releases/R2008/manifests/00-contrail-09-manager.yaml -o manifests/00-contrail-09-manager.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/manifests/cluster-network-02-config.yml -o manifests/cluster-network-02-config.yml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_master-iptables-machine-config.yaml -o openshift/99_master-iptables-machine-config.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_master-kernel-modules-overlay.yaml -o openshift/99_master-kernel-modules-overlay.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_master_network_functions.yaml -o openshift/99_master_network_functions.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_master_network_manager_stop_service.yaml -o openshift/99_master_network_manager_stop_service.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_master-pv-mounts.yaml -o openshift/99_master-pv-mounts.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_worker-iptables-machine-config.yaml -o openshift/99_worker-iptables-machine-config.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_worker-kernel-modules-overlay.yaml -o openshift/99_worker-kernel-modules-overlay.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_worker_network_functions.yaml -o openshift/99_worker_network_functions.yaml
curl https://raw.githubusercontent.com/Juniper/contrail-operator/R2008/deploy/openshift/openshift/99_worker_network_manager_stop_service.yaml -o openshift/99_worker_network_manager_stop_service.yaml
```

Modify `manifests/00-contrail-02-registry-secret.yaml` file providing proper configuration with credentials to *hub.juniper.net* registry.

**NOTE:** You may create base64 encoded value for config with script provided in [here](https://github.com/Juniper/contrail-operator/tree/master/deploy/openshift/tools/docker-config-generate) directory.
Copy output of the script and paste into contrail registry secret manifest.

4. Modify manifests if neccessary:

If you use pod/service network CIDRs other then the default values open the  **deploy/openshift/manifests/cluster-network-02-config.yml** in text editor and update CIDR values.

If you deploy more/less master nodes than the default 3, modify the **deploy/openshift/manifests/00-contrail-09-manager.yaml** in text editor and set the
spec.commonConfiguration.replicas field to the number of master nodes (modify only the top-level replicas field in the file).

5. Install manifest

Install manifests with (you can execute the scripts from anywhere, the example assumes that you are in the contrail-operator repository root directory):
```
./deploy/openshift/install-manifests.sh --dir <name of openshift install directory> --config <path to config file>
```
**NOTE** If **--config** is not provided by default script will try to read config from script directory's file **config**

6. Install Openshift
Run this command to start Openshift install:
```
./openshift-install create cluster --dir <name of openshift install directory>
```

7. Open security groups:

Login to AWS Console and find *master* instance created by the *openshift-installer*. Select Security Group attached to it and edit it's inbound rules to accept all traffic. **Do the same for the security group attached to worker nodes, after they are created.** To automatically open port required by Contrail, you can use the [contrail-sc-open](https://github.com/Juniper/contrail-operator/tree/master/deploy/openshift/tools/contrail-sc-open) tool.


1. Patch the externalTrafficPolicy

Verify that the **router-default** service has been created, by running:
```
kubectl -n openshift-ingress describe service router-default
```
If it is not present yet, wait until it is created. Then patch the externalTrafficPolicy by running this command:

```
kubectl -n openshift-ingress patch service router-default --patch '{"spec": {"externalTrafficPolicy": "Cluster"}}'
```

# Access cluster
In order to access export **KUBECONFIG** environment variable.
**KUBECONFIG** file may be found under **<Openshift install directory>/auth/kubeconfig**
E.x.
```
export KUBECONFIG=<Openshift install directory>/auth/kubeconfig
```
Afterwards cluster may be accessed with `kubectl` command line tool.

It's also possible to access cluster with dedicated Openshift command line tool: `oc`.
However, `oc` requires to login before.
After successful deployment **openshift-install** binary prints out username (**kubeadmin**) and password to cluster.
Password may be also found also under **<Openshift install directory>/auth/** directory.

Login into `oc` may be performed with this command:
```
oc login -u kubeadmin -p <cluster password>
```

Last method to access Openshift cluster is web console.
URL to web console will be displayed by **openshift-install** binary at the end of deployment.
Login into console with the same credentials as for `oc`.

## Notes

* Contrail Operator creates Persistent Volumes that are used by some of the deployed pods. After deletion of Contrail resources (e.g. after deleting the Manager Custom Resource), those Persistent Volumes will not be deleted. Administrator has to delete them manually and make sure that directories created by these volumes on cluster nodes are in the expected state. Example Persistent Volumes deletion command:
```
kubectl delete pv $(kubectl get pv -o=jsonpath='{.items[?(@.spec.storageClassName=="local-storage")].metadata.name}')
