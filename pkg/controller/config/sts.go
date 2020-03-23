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
      initContainers:
        - name: init
          image: busybox
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
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
        - name: nodeinit
          image: docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1
          env:
            - name: CONTRAIL_STATUS_IMAGE
              value: docker.io/michaelhenkel/contrail-status:5.2.0-dev1
          imagePullPolicy: Always
          securityContext:
            privileged: true
          volumeMounts:
            - mountPath: /host/usr/bin
              name: host-usr-local-bin
      containers:
        - name: api
          image: docker.io/michaelhenkel/contrail-controller-config-api:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          readinessProbe:
            httpGet:
              scheme: HTTPS
              path: /
              port: 8082
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: devicemanager
          image: docker.io/michaelhenkel/contrail-controller-config-devicemgr:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
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
            - mountPath: /var/log/contrail
              name: config-logs
        - name: schematransformer
          image: docker.io/michaelhenkel/contrail-controller-config-schema:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: servicemonitor
          image: docker.io/michaelhenkel/contrail-controller-config-svcmonitor:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: analyticsapi
          image: docker.io/michaelhenkel/contrail-analytics-api:5.2.0-dev1
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
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: queryengine
          image: docker.io/michaelhenkel/contrail-analytics-query-engine:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: collector
          image: docker.io/michaelhenkel/contrail-analytics-collector:5.2.0-dev1
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
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
                name: config-logs
              - mountPath: /var/lib/redis
                name: config-data
        - name: nodemanagerconfig
          image: docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1
          env:
            - name: DOCKER_HOST
              value: unix://mnt/docker.sock
            - name: NODE_TYPE
              value: config
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /mnt
              name: docker-unix-socket
            - mountPath: /var/crashes
              name: crashes
        - name: nodemanageranalytics
          image: docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1
          env:
            - name: DOCKER_HOST
              value: unix://mnt/docker.sock
            - name: NODE_TYPE
              value: analytics
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /mnt
              name: docker-unix-socket
            - mountPath: /var/crashes
              name: crashes
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
        - persistentVolumeClaim: {}
          name: tftp
        - persistentVolumeClaim: {}
          name: dnsmasq
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
