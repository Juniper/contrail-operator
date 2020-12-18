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

3. Create configuration file
Create file with configuration parameters that looks similar to this:
```
contrail_version=master.1460-ubi
contrail_registry=hub.juniper.net/contrail-nightly
docker_config=example_json_config
```
or
```
CONTRAIL_VERSION=2011.138
CONTRAIL_REGISTRY=hub.juniper.net/contrail
DOCKER_CONFIG=example_json_config
```

Under *CONTRAIL_VERSION* field enter proper Contrail container build tag, available in the hub.juniper.net/contrail-nightly registry.

Choose source registry for container images with *CONTRAIL_REGISTRY* field.

*DOCKER_CONFIG* is configuration for registry secret to closed container registry (if registry is wide open then no credentials are required)
Set *DOCKER_CONFIG* to registry secret with proper data in base64.

**NOTE:** You may create base64 encoded value for config with script provided in [deploy/openshift/tools/docker-config-generate](tools/docker-config-generate) directory.
Copy output of the script and paste into config used to install-manifests script.

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

Login to AWS Console and find *master* instance created by the *openshift-installer*. Select Security Group attached to it and edit it's inbound rules to accept all traffic. **Do the same for the security group attached to worker nodes, after they are created.** To automatically open port required by Contrail, you can use the [contrail-sc-open](tools/contrail-sc-open/) tool. 


8. Patch the externalTrafficPolicy

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
