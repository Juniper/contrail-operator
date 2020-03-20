package webui

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDatawebui_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: webui
spec:
  selector:
    matchLabels:
      app: webui
  serviceName: "webui"
  replicas: 1
  template:
    metadata:
      labels:
        app: webui
        contrail_manager: webui
    spec:
      containers:
      - image: docker.io/michaelhenkel/contrail-controller-webui-web:5.2.0-dev1
        env:
        - name: WEBUI_SSL_KEY_FILE
          value: /etc/contrail/webui_ssl/cs-key.pem
        - name: WEBUI_SSL_CERT_FILE
          value: /etc/contrail/webui_ssl/cs-cert.pem
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: webuiweb
        volumeMounts:
        - mountPath: /var/log/contrail
          name: webui-logs
      - image: docker.io/michaelhenkel/contrail-controller-webui-job:5.2.0-dev1
        env:
        - name: WEBUI_SSL_KEY_FILE
          value: /etc/contrail/webui_ssl/cs-key.pem
        - name: WEBUI_SSL_CERT_FILE
          value: /etc/contrail/webui_ssl/cs-cert.pem
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: webuijob
        volumeMounts:
        - mountPath: /var/log/contrail
          name: webui-logs
      - image: docker.io/michaelhenkel/contrail-external-redis:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: redis
        volumeMounts:
        - mountPath: /var/log/contrail
          name: webui-logs
        - mountPath: /var/lib/redis
          name: webui-data
      dnsPolicy: ClusterFirst
      hostNetwork: true
      initContainers:
      - env:
        - name: CONTRAIL_STATUS_IMAGE
          value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1
        imagePullPolicy: Always
        name: nodeinit
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /host/usr/bin
          name: host-usr-local-bin
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
      nodeSelector:
        node-role.kubernetes.io/master: ""
      restartPolicy: Always
      tolerations:
      - effect: NoSchedule
        operator: Exists
      - effect: NoExecute
        operator: Exists
      volumes:
      - hostPath:
          path: /var/lib/contrail/webui
          type: ""
        name: webui-data
      - hostPath:
          path: /var/log/contrail/webui
          type: ""
        name: webui-logs
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
	err := yaml.Unmarshal([]byte(yamlDatawebui_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDatawebui_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
