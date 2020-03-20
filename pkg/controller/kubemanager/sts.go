package kubemanager

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDatakubemanager_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kubemanager
spec:
  selector:
    matchLabels:
      app: kubemanager
  serviceName: "kubemanager"
  replicas: 1
  template:
    metadata:
    labels:
      app: kubemanager
      contrail_manager: kubemanager
    spec:
      nodeSelector:
        node-role.kubernetes.io/master: ""
      serviceAccount: contrail-service-account-kubemanager
      serviceAccountName: contrail-service-account-kubemanager
      dnsPolicy: ClusterFirst
      hostNetwork: true
      initContainers:
        - name: init
          image: busybox
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
          imagePullPolicy: Always
        - name: nodeinit
          image: docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
          volumeMounts:
            - mountPath: /host/usr/bin
              name: host-usr-local-bin
          securityContext:
            privileged: true
          imagePullPolicy: Always
      containers:
        - name: kubemanager
          image: docker.io/michaelhenkel/contrail-kubernetes-kube-manager:5.2.0-dev1
          volumeMounts:
            - mountPath: /var/log/contrail
              name: kubemanager-logs
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
      volumes:
        - hostPath:
            path: /var/log/contrail/kubemanager
            type: ""
          name: kubemanager-logs
        - hostPath:
            path: /usr/local/bin
            type: ""
          name: host-usr-local-bin
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
	err := yaml.Unmarshal([]byte(yamlDatakubemanager_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDatakubemanager_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
