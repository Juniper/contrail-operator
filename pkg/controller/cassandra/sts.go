package cassandra

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDatacassandra_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: cassandra
spec:
  selector:
    matchLabels:
      app: cassandra
  serviceName: "cassandra"
  replicas: 1
  template:
    metadata:
      labels:
        app: cassandra
        cassandra_cr: cassandra
        contrail_manager: cassandra
    spec:
      terminationGracePeriodSeconds: 1800
      containers:
      - image: hub.juniper.net/contrail-nightly/contrail-external-cassandra:5.2.0-0.740
        imagePullPolicy: Always
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        lifecycle:
          preStop:
            exec:
              command: 
              - /bin/sh
              - -c
              - nodetool drain
              #- nodetool decommission
        readinessProbe:
          exec:
            command:
            - /bin/bash
            - -c
            - "seeds=$(for i in $(ls /mydata/*.yaml); do echo $(basename $i .yaml); done) &&  for seed in $(echo $seeds); do if [[ $(nodetool status | grep $seed |awk '{print $1}') != 'UN' ]]; then exit -1; fi; done"
          initialDelaySeconds: 15
          timeoutSeconds: 5
        name: cassandra
        securityContext:
          capabilities:
            add:
              - IPC_LOCK
          privileged: false
          procMount: Default
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        #volumeMounts:
        #- mountPath: /var/log/cassandra
        #  name: cassandra-logs
        #- mountPath: /var/lib/cassandra
        #  name: cassandra-data
      dnsPolicy: ClusterFirst
      hostNetwork: true
      initContainers:
      - command:
        - sh
        - -c
        - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: busybox:1.31
        imagePullPolicy: Always
        name: init
        resources: {}
        securityContext:
          privileged: false
          procMount: Default
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /tmp/podinfo
          name: status
      - command:
        - sh
        - -c
        - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: busybox:1.31
        imagePullPolicy: Always
        name: init2
        resources: {}
        securityContext:
          privileged: false
          procMount: Default
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /tmp/podinfo
          name: status
      nodeSelector:
        node-role.kubernetes.io/master: ""
      restartPolicy: Always
      schedulerName: default-scheduler
      tolerations:
      - effect: NoSchedule
        operator: Exists
      - effect: NoExecute
        operator: Exists
      volumes:
      #- hostPath:
      #    path: /var/log/contrail/cassandra
      #    type: ""
      #  name: cassandra-logs
      #- hostPath:
      #    path: /var/lib/contrail/cassandra
      #    type: ""
      #  name: cassandra-data
      - downwardAPI:
          defaultMode: 420
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels
            path: pod_labels
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels
            path: pod_labelsx
        name: status
  volumeClaimTemplates:
  - metadata:
      name: cassandra-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 5G
  - metadata:
      name: cassandra-logs
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 5G`

func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlDatacassandra_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDatacassandra_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
