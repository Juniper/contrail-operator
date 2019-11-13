package config

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDataconfig_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: config
spec:
  selector:
    matchLabels:
      app: config
  serviceName: "config"
  replicas: 1
  template:
    metadata:
      labels:
        app: config
        contrail_manager: config
    spec:
      containers:
      - image: docker.io/michaelhenkel/contrail-controller-config-api:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: api
        readinessProbe:
          httpGet:
            path: /
            port: 8082
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-controller-config-devicemgr:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: devicemanager
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-controller-config-schema:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: schematransformer
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-controller-config-svcmonitor:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: servicemonitor
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-analytics-api:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: ANALYTICSDB_ENABLE
          value: "true"
        - name: ANALYTICS_ALARM_ENABLE
          value: "true"
        imagePullPolicy: Always
        name: analyticsapi
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-analytics-query-engine:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: queryengine
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
      - image: docker.io/michaelhenkel/contrail-analytics-collector:5.2.0-dev1
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        imagePullPolicy: Always
        name: collector
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
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
          name: config-logs
        - mountPath: /var/lib/redis
          name: config-data
      - env:
        - name: DOCKER_HOST
          value: unix://mnt/docker.sock
        - name: NODE_TYPE
          value: config
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1
        imagePullPolicy: Always
        name: nodemanagerconfig
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
        - mountPath: /mnt
          name: docker-unix-socket
        - mountPath: /var/crashes
          name: crashes
      - env:
        - name: DOCKER_HOST
          value: unix://mnt/docker.sock
        - name: NODE_TYPE
          value: analytics
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1
        imagePullPolicy: Always
        name: nodemanageranalytics
        volumeMounts:
        - mountPath: /var/log/contrail
          name: config-logs
        - mountPath: /mnt
          name: docker-unix-socket
        - mountPath: /var/crashes
          name: crashes
      dnsPolicy: ClusterFirst
      hostNetwork: true
      initContainers:
      - command:
        - sh
        - -c
        - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
        env:
        - name: CONTRAIL_STATUS_IMAGE
          value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
        image: busybox
        imagePullPolicy: Always
        name: init
        volumeMounts:
        - mountPath: /tmp/podinfo
          name: status
      - command:
        - sh
        - -c
        - until grep true /tmp/podinfo/peers_ready > /dev/null 2>&1; do sleep 1; done
        image: busybox
        imagePullPolicy: Always
        name: init2
        volumeMounts:
        - mountPath: /tmp/podinfo
          name: status
      - env:
        - name: CONTRAIL_STATUS_IMAGE
          value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
        image: docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1
        imagePullPolicy: Always
        name: nodeinit
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /host/usr/bin
          name: host-usr-bin
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
      - effect: NoSchedule
        operator: Exists
      - effect: NoExecute
        operator: Exists
      volumes:
      - hostPath:
          path: /var/log/contrail/config
          type: ""
        name: config-logs
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
          path: /usr/bin
          type: ""
        name: host-usr-bin
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

var vncApiLib = `[global]
;WEB_SERVER = 127.0.0.1
;WEB_PORT = 9696  ; connection through quantum plugin

WEB_SERVER = localhost
WEB_PORT = 8082 ; connection to api-server directly
BASE_URL = /
;BASE_URL = /tenants/infra ; common-prefix for all URLs

; Authentication settings (optional)
[auth]
;AUTHN_TYPE = keystone
;AUTHN_PROTOCOL = http
;AUTHN_SERVER = 127.0.0.1
;AUTHN_PORT = 35357
;AUTHN_URL = /v2.0/tokens
;AUTHN_TOKEN_URL = http://127.0.0.1:35357/v2.0/tokens
`
