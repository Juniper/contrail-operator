package rabbitmq

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlRabbitmq_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: rabbitmq
spec:
  selector:
    matchLabels:
      app: rabbitmq
  serviceName: "rabbitmq"
  replicas: 1
  template:
    metadata:
      labels:
        app: rabbitmq
        contrail_manager: rabbitmq
    spec:
      nodeSelector:
        node-role.kubernetes.io/master: ''
      tolerations:
      - operator: Exists
        effect: NoSchedule
      - operator: Exists
        effect: NoExecute
      hostNetwork: true
      initContainers:
      - command:
        - sh
        - -c
        - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
        image: busybox
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
      containers:
      - name: rabbitmq
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        image: docker.io/michaelhenkel/contrail-external-rabbitmq:5.2.0-dev1
        imagePullPolicy: ""
        volumeMounts:
        - mountPath: /var/lib/rabbitmq
          name: rabbitmq-data
        - mountPath: /var/log/rabbitmq
          name: rabbitmq-logs
        readinessProbe:
          exec:
            command:
            - /bin/bash
            - -c
            - "export RABBITMQ_NODENAME=rabbit@$POD_IP; cluster_status=$(rabbitmqctl cluster_status);nodes=$(echo $cluster_status | sed -e 's/.*disc,\\[\\(.*\\)]}]}, {.*/\\1/' | grep -oP \"(?<=rabbit@).*?(?=')\"); for node in $(cat /etc/rabbitmq/rabbitmq.nodes); do echo ${nodes} |grep ${node}; if [[ $? -ne 0 ]]; then exit -1; fi; done"
          initialDelaySeconds: 15
          timeoutSeconds: 5
      volumes:
      - name: rabbitmq-data
        hostPath:
          path: /var/lib/contrail/rabbitmq
      - name: rabbitmq-logs
        hostPath:
          path: /var/log/contrail/rabbitmq
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
        name: status`

func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlRabbitmq_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlRabbitmq_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
