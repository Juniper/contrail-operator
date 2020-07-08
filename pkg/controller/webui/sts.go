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
      initContainers:
        - name: nodeinit
          image: docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: docker.io/opencontrailnightly/contrail-status:latest
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          securityContext:
            privileged: true
          volumeMounts:
            - mountPath: /host/usr/bin
              name: host-usr-local-bin
        - name: init
          image: busybox
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: hub.juniper.net/contrail-nightly/contrail-status:5.2.0-0.740
          imagePullPolicy: Always
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
        - name: webuiweb
          image: docker.io/michaelhenkel/contrail-controller-webui-web:5.2.0-dev1
          env:
            - name: WEBUI_SSL_KEY_FILE
              value: /etc/contrail/webui_ssl/cs-key.pem
            - name: WEBUI_SSL_CERT_FILE
              value: /etc/contrail/webui_ssl/cs-cert.pem
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: ANALYTICSDB_ENABLE
              value: "true"
            - name: ANALYTICS_SNMP_ENABLE
              value: "true"
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: webui-logs
        - name: webuijob
          image: docker.io/michaelhenkel/contrail-controller-webui-job:5.2.0-dev1
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
          volumeMounts:
            - mountPath: /var/log/contrail
              name: webui-logs
        - name: redis
          image: docker.io/michaelhenkel/contrail-external-redis:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: webui-logs
            - mountPath: /var/lib/redis
              name: webui-data
      dnsPolicy: ClusterFirst
      hostNetwork: true
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

// GetSTS returns StatesfulSet object created from YAML yamlDatawebui_sts
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
