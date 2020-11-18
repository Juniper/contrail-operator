# Openshift 4.x chrony configuration

By default, Openshift 4 will configure its nodes to use public NTP servers. To use a specific NTP server, for example from an internal network, user has to modify the chrony configuration at /etc/chrony.conf.


## Steps

Execute these steps after you've run the `openshift-install create manifests` command and the `install-manifests.sh` script:

1. Encode the chrony configuration you want to apply with base64 and store it in an environment variable (instead of `<NTP_SERVER>` specify the NTP server you want to use):
```
export CHRONY_CONF_BASE64=$(cat << EOF | base64 -w 0
server <NTP_SERVER> iburst
driftfile /var/lib/chrony/drift
makestep 1.0 3
rtcsync
logdir /var/log/chrony
EOF          
)
```
2. Create MachineConfig for master nodes inside the `openshift` folder with other MachineConfigs:
```
cat << EOF > ./openshift/99_masters-chrony-configuration.yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: master
  name: masters-chrony-configuration
spec:
  config:
    ignition:
      config: {}
      security:
        tls: {}
      timeouts: {}
      version: 3.1.0
    networkd: {}
    passwd: {}
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,$CHRONY_CONF_BASE64
          verification: {}
        filesystem: root
        mode: 420
        path: /etc/chrony.conf
EOF
```
3. Create MachineConfig for worker nodes:
```
cat << EOF > ./openshift/99_workers-chrony-configuration.yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: workers-chrony-configuration
spec:
  config:
    ignition:
      config: {}
      security:
        tls: {}
      timeouts: {}
      version: 3.1.0
    networkd: {}
    passwd: {}
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,$CHRONY_CONF_BASE64
          verification: {}
        filesystem: root
        mode: 420
        path: /etc/chrony.conf
EOF
```
4. Verify the files in the `openshift` directory:
```
cat ./openshift/99_masters-chrony-configuration.yaml
cat ./openshift/99_workers-chrony-configuration.yaml
```

5. Continue Openshift deployment as specified in the [guide](./Openshift-KVM.md)

6. As soon as the bootstrap node comes up, ssh to the node and run these commands as the root user (instead of `<NTP_SERVER>` specify the NTP server you want to use):
```
cat << EOF > /etc/chrony.conf
server <NTP_SERVER> iburst
driftfile /var/lib/chrony/drift
makestep 1.0 3
rtcsync
logdir /var/log/chrony
EOF          

systemctl restart chronyd
```
