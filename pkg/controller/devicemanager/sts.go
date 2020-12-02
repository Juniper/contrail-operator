package devicemanager

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDataconfig_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: devicemanager
spec:
  selector:
    matchLabels:
      app: config
  serviceName: "devicemanager"
  replicas: 1
  template:
    metadata:
      labels:
        app: devicemanager
        contrail_manager: devicemanager
    spec:
      initContainers:
        - name: init
          image: busybox
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
        - name: init2
          image: busybox
          command:
            - sh
            - -c
            - until grep true /tmp/podinfo/peers_ready > /dev/null 2>&1; do sleep 1; done
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
      containers:
        - name: devicemanager
          image: docker.io/michaelhenkel/contrail-controller-config-devicemgr:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail/device-manager
              name: config-device-manager-logs
        - name: dnsmasq
          image: docker.io/michaelhenkel/contrail-external-dnsmasq:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: CONTROLLER_NODES
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail/device-manager
              name: config-device-manager-logs
        - name: statusmonitor
          image: docker.io/kaweue/contrail-statusmonitor:debug
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail/device-manager
              name: config-device-manager-logs
      dnsPolicy: ClusterFirst
      hostNetwork: true
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
      volumes:
        - emptyDir: {}
          name: tftp
        - emptyDir: {}
          name: dnsmasq
        - hostPath:
            path: /var/log/contrail/config-device-manager
            type: ""
          name: config-device-manager-logs
        - hostPath:
            path: /var/contrail/crashes
            type: ""
          name: crashes
        - hostPath:
            path: /var/lib/contrail/config
            type: ""
          name: config-data
        - hostPath:
            path: /var/run
            type: ""
          name: docker-unix-socket
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
              path: peers_ready
            - fieldRef:
                apiVersion: v1
                fieldPath: metadata.labels
              path: pod_labelsx
          name: status`

func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlDataconfig_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataconfig_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
