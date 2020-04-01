package provisionmanager

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlProvisionManager_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: provisionmanager
spec:
  selector:
    matchLabels:
      app: provisionmanager
  serviceName: "provisionmanager"
  replicas: 1
  template:
    metadata:
      labels:
        app: provisionmanager
        contrail_manager: provisionmanager
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
        env:
        - name: CONTRAIL_STATUS_IMAGE
          value: hub.juniper.net/contrail-nightly/contrail-status:5.2.0-0.740
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
      - name: provisioner
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        image: docker.io/kaweue/contrail-provisioner:master.1115
        imagePullPolicy: Always
        volumeMounts:
        - mountPath: /var/lib/provisionmanager
          name: provisionmanager-data
        - mountPath: /var/log/provisionmanager
          name: provisionmanager-logs
      volumes:
      - name: provisionmanager-data
        hostPath:
          path: /var/lib/contrail/provisionManager
      - name: provisionmanager-logs
        hostPath:
          path: /var/log/contrail/provisionmanager
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
	err := yaml.Unmarshal([]byte(yamlProvisionManager_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlProvisionManager_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
