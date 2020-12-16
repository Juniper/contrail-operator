package contrailtest

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

var config = &v1alpha1.Config{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "config1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: v1alpha1.ConfigSpec{
		ServiceConfiguration: v1alpha1.ConfigConfiguration{
			KeystoneSecretName: "keystone-adminpass-secret",
			AuthMode:           v1alpha1.AuthenticationModeKeystone,
			CassandraInstance:  "cassandra1",
			ZookeeperInstance:  "zookeeper1",
			KeystoneInstance:   "keystone",
		},
	},
}

var configService = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "config1-api",
		Namespace: "default",
		Labels:    map[string]string{"service": "config1"},
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{Port: int32(v1alpha1.ConfigApiPort), Protocol: "TCP", Name: "api"},
			{Port: int32(v1alpha1.AnalyticsApiPort), Protocol: "TCP", Name: "analytics"},
		},
		ClusterIP: "10.10.10.10",
	},
}

var control = &v1alpha1.Control{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "control1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
			"control_role":     "master",
		},
	},
	Spec: v1alpha1.ControlSpec{
		ServiceConfiguration: v1alpha1.ControlConfiguration{
			CassandraInstance: "cassandra1",
		},
	},
}

var kubemanager = &v1alpha1.Kubemanager{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "kubemanager1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: v1alpha1.KubemanagerSpec{
		ServiceConfiguration: v1alpha1.KubemanagerServiceConfiguration{
			KubemanagerConfiguration: v1alpha1.KubemanagerConfiguration{
				AuthMode: "keystone",
			},
			KubemanagerNodesConfiguration: v1alpha1.KubemanagerNodesConfiguration{
				CassandraNodesConfiguration: &v1alpha1.CassandraClusterConfiguration{
					Port:         9160,
					ServerIPList: []string{"1.1.2.1", "1.1.2.2", "1.1.2.3"},
				},
				ZookeeperNodesConfiguration: &v1alpha1.ZookeeperClusterConfiguration{
					ClientPort:   2181,
					ServerIPList: []string{"1.1.3.1", "1.1.3.2", "1.1.3.3"},
				},
				RabbbitmqNodesConfiguration: &v1alpha1.RabbitmqClusterConfiguration{
					SSLPort:      15673,
					ServerIPList: []string{"1.1.4.1", "1.1.4.2", "1.1.4.3"},
					Secret:       "rabbitmq-secret",
				},
				ConfigNodesConfiguration: &v1alpha1.ConfigClusterConfiguration{
					APIServerPort:         8082,
					CollectorPort:         8086,
					APIServerIPList:       []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"},
					CollectorServerIPList: []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"},
					AuthMode:              v1alpha1.AuthenticationModeKeystone,
				},
				KeystoneNodesConfiguration: &v1alpha1.KeystoneClusterConfiguration{
					Port:     5555,
					Endpoint: "10.11.12.13",
				},
			},
		},
	},
}

var webui = &v1alpha1.Webui{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "webui1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: v1alpha1.WebuiSpec{
		ServiceConfiguration: v1alpha1.WebuiConfiguration{
			KeystoneSecretName: "keystone-adminpass-secret",
			CassandraInstance:  "cassandra1",
			KeystoneInstance:   "keystone",
		},
	},
}

var cassandra = &v1alpha1.Cassandra{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cassandra1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Status: v1alpha1.CassandraStatus{ClusterIP: "10.0.0.1"},
}

var zookeeper = &v1alpha1.Zookeeper{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "zookeeper1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
}

var rabbitmq = &v1alpha1.Rabbitmq{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "rabbitmq",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
}

var vrouter = &v1alpha1.Vrouter{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "vrouter",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: v1alpha1.VrouterSpec{
		ServiceConfiguration: v1alpha1.VrouterServiceConfiguration{
			VrouterConfiguration: v1alpha1.VrouterConfiguration{
				Gateway: "1.1.8.254",
			},
			VrouterNodesConfiguration: v1alpha1.VrouterNodesConfiguration{
				ConfigNodesConfiguration: &v1alpha1.ConfigClusterConfiguration{
					APIServerIPList:       []string{"1.1.5.1", "1.1.5.2", "1.1.5.3"},
					AnalyticsServerIPList: []string{"1.1.5.1", "1.1.5.2", "1.1.5.3"},
					CollectorServerIPList: []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"},
				},
				ControlNodesConfiguration: &v1alpha1.ControlClusterConfiguration{
					ControlServerIPList: []string{"1.1.5.1", "1.1.5.2", "1.1.5.3"},
				},
			},
		},
	},
}

var keystone = &v1alpha1.Keystone{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "keystone",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: v1alpha1.KeystoneSpec{
		ServiceConfiguration: v1alpha1.KeystoneConfiguration{
			ListenPort:        5555,
			AuthProtocol:      "https",
			UserDomainName:    "Default",
			ProjectDomainName: "Default",
			Region:            "RegionOne",
		},
	},
	Status: v1alpha1.KeystoneStatus{
		Endpoint: "10.11.12.14",
	},
}

var rabbitmqList = &v1alpha1.RabbitmqList{}
var zookeeperList = &v1alpha1.ZookeeperList{}
var cassandraList = &v1alpha1.CassandraList{}
var configList = &v1alpha1.ConfigList{}
var controlList = &v1alpha1.ControlList{}
var kubemanagerList = &v1alpha1.KubemanagerList{}
var webuiList = &v1alpha1.WebuiList{}
var vrouterList = &v1alpha1.VrouterList{}

var configMap = &corev1.ConfigMap{}
var secret = &corev1.Secret{}

type Environment struct {
	client                   *client.Client
	configPodList            corev1.PodList
	rabbitmqPodList          corev1.PodList
	zookeeperPodList         corev1.PodList
	cassandraPodList         corev1.PodList
	controlPodList           corev1.PodList
	kubemanbagerPodList      corev1.PodList
	webuiPodList             corev1.PodList
	vrouterPodList           corev1.PodList
	configResource           v1alpha1.Config
	controlResource          v1alpha1.Control
	cassandraResource        v1alpha1.Cassandra
	zookeeperResource        v1alpha1.Zookeeper
	rabbitmqResource         v1alpha1.Rabbitmq
	kubemanagerResource      v1alpha1.Kubemanager
	webuiResource            v1alpha1.Webui
	vrouterResource          v1alpha1.Vrouter
	configConfigMap          corev1.ConfigMap
	controlConfigMap         corev1.ConfigMap
	cassandraConfigMap       corev1.ConfigMap
	zookeeperConfigMap       corev1.ConfigMap
	zookeeperConfigMap2      corev1.ConfigMap
	rabbitmqConfigMap        corev1.ConfigMap
	rabbitmqConfigMap2       corev1.ConfigMap
	kubemanagerConfigMap     corev1.ConfigMap
	kubemanagerConfigMapEnvs corev1.ConfigMap
	kubemanagerSecret        corev1.Secret
	webuiConfigMap           corev1.ConfigMap
	vrouterConfigMap         corev1.ConfigMap
	vrouterConfigMap2        corev1.ConfigMap
}

func SetupEnv() Environment {
	logf.SetLogger(logf.ZapLogger(true))
	configConfigMap := *configMap
	rabbitmqConfigMap := *configMap
	rabbitmqConfigMap2 := *configMap
	cassandraConfigMap := *configMap
	zookeeperConfigMap := *configMap
	controlConfigMap := *configMap
	kubemanagerConfigMap := *configMap
	kubemanagerConfigMapEnvs := *configMap
	webuiConfigMap := *configMap
	vrouterConfigMap := *configMap
	vrouterConfigMap2 := *configMap
	kubemanagerSecret := *secret
	keystoneAdminSecret := *secret

	kubemanagerSecret.Name = "kubemanagersecret"
	kubemanagerSecret.Namespace = "default"
	kubemanagerSecret.Annotations = map[string]string{"kubernetes.io/service-account.name": "contrail-service-account"}
	kubemanagerSecret.Type = corev1.SecretType("kubernetes.io/service-account-token")
	var data = make(map[string][]byte)
	data["token"] = []byte("THISISATOKEN")
	kubemanagerSecret.Data = data

	rabbitmqSecret := *secret

	rabbitmqSecret.Name = "rabbitmq-secret"
	rabbitmqSecret.Namespace = "default"
	rabbitmqSecret.Annotations = map[string]string{"kubernetes.io/service-account.name": "contrail-service-account"}

	rabbitmqSecret.Data = map[string][]byte{
		"user":     []byte("user"),
		"password": []byte("password"),
		"vhost":    []byte("vhost"),
	}

	cassandraSecret := *secret

	cassandraSecret.Name = "cassandra1-secret"
	cassandraSecret.Namespace = "default"
	cassandraSecret.Annotations = map[string]string{"kubernetes.io/service-account.name": "contrail-service-account"}

	cassandraSecret.Data = map[string][]byte{
		"keystorePassword":   []byte("keystorePassword"),
		"truststorePassword": []byte("truststorePassword"),
	}

	configConfigMap.Name = "config1-config-configmap"
	configConfigMap.Namespace = "default"

	rabbitmqConfigMap.Name = "rabbitmq1-rabbitmq-configmap"
	rabbitmqConfigMap.Namespace = "default"

	rabbitmqConfigMap2.Name = "rabbitmq1-rabbitmq-configmap-runner"
	rabbitmqConfigMap2.Namespace = "default"

	cassandraConfigMap.Name = "cassandra1-cassandra-configmap"
	cassandraConfigMap.Namespace = "default"

	zookeeperConfigMap.Name = "zookeeper1-zookeeper-configmap"
	zookeeperConfigMap.Namespace = "default"

	controlConfigMap.Name = "control1-control-configmap"
	controlConfigMap.Namespace = "default"

	kubemanagerConfigMap.Name = "kubemanager1-kubemanager-configmap"
	kubemanagerConfigMap.Namespace = "default"

	kubemanagerConfigMapEnvs.Name = "kubemanager1-kubemanager-configmap-envs"
	kubemanagerConfigMapEnvs.Namespace = "default"

	webuiConfigMap.Name = "webui1-webui-configmap"
	webuiConfigMap.Namespace = "default"

	vrouterConfigMap.Name = "vrouter1-vrouter-configmap"
	vrouterConfigMap.Namespace = "default"

	vrouterConfigMap2.Name = "vrouter1-vrouter-configmap-1"
	vrouterConfigMap2.Namespace = "default"

	keystoneAdminSecret.Name = "keystone-adminpass-secret"
	keystoneAdminSecret.Namespace = "default"
	keystoneAdminSecret.Annotations = map[string]string{"kubernetes.io/service-account.name": "contrail-service-account"}
	keystoneAdminSecret.Data = map[string][]byte{
		"password": []byte("test123"),
	}

	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion,
		config,
		cassandra,
		zookeeper,
		rabbitmq,
		control,
		kubemanager,
		webui,
		vrouter,
		keystone,
		rabbitmqList,
		zookeeperList,
		cassandraList,
		configList,
		controlList,
		kubemanagerList,
		webuiList,
		vrouterList)

	objs := []runtime.Object{config,
		cassandra,
		zookeeper,
		rabbitmq,
		control,
		kubemanager,
		webui,
		vrouter,
		keystone,
		&configConfigMap,
		&controlConfigMap,
		&cassandraConfigMap,
		&zookeeperConfigMap,
		&rabbitmqConfigMap,
		&rabbitmqConfigMap2,
		&kubemanagerConfigMap,
		&kubemanagerConfigMapEnvs,
		&webuiConfigMap,
		&vrouterConfigMap,
		&vrouterConfigMap2,
		&rabbitmqSecret,
		&cassandraSecret,
		&kubemanagerSecret,
		&keystoneAdminSecret}

	cl := fake.NewFakeClient(objs...)

	configResource := config

	controlResource := control

	rabbitmqResource := rabbitmq

	zookeeperResource := zookeeper

	cassandraResource := cassandra

	kubemanagerResource := kubemanager

	webuiResource := webui

	vrouterResource := vrouter

	var podServiceMap = make(map[string]map[string]string)
	podServiceMap["configPods"] = map[string]string{"pod1": "1.1.1.1", "pod2": "1.1.1.2", "pod3": "1.1.1.3"}
	podServiceMap["rabbitmqPods"] = map[string]string{"pod1": "1.1.4.1", "pod2": "1.1.4.2", "pod3": "1.1.4.3"}
	podServiceMap["cassandraPods"] = map[string]string{"pod1": "1.1.2.1", "pod2": "1.1.2.2", "pod3": "1.1.2.3"}
	podServiceMap["zookeeperPods"] = map[string]string{"pod1": "1.1.3.1", "pod2": "1.1.3.2", "pod3": "1.1.3.3"}
	podServiceMap["controlPods"] = map[string]string{"pod1": "1.1.5.1", "pod2": "1.1.5.2", "pod3": "1.1.5.3"}
	podServiceMap["kubemanagerPods"] = map[string]string{"pod1": "1.1.6.1", "pod2": "1.1.6.2", "pod3": "1.1.6.3"}
	podServiceMap["webuiPods"] = map[string]string{"pod1": "1.1.7.1", "pod2": "1.1.7.2", "pod3": "1.1.7.3"}
	podServiceMap["vrouterPods"] = map[string]string{"pod1": "1.1.8.1", "pod2": "1.1.8.2", "pod3": "1.1.8.3"}

	type PodMap struct {
		configPods      map[string]string
		rabbitmqPods    map[string]string
		cassandraPods   map[string]string
		zookeeperPods   map[string]string
		controlPods     map[string]string
		kubemanagerPods map[string]string
		webuiPods       map[string]string
		vrouterPods     map[string]string
	}
	podMap := PodMap{
		configPods:      map[string]string{"pod-0": "1.1.1.1", "pod-1": "1.1.1.2", "pod-2": "1.1.1.3"},
		rabbitmqPods:    map[string]string{"pod-0": "1.1.4.1", "pod-1": "1.1.4.2", "pod-2": "1.1.4.3"},
		cassandraPods:   map[string]string{"pod-0": "1.1.2.1", "pod-1": "1.1.2.2", "pod-2": "1.1.2.3"},
		zookeeperPods:   map[string]string{"pod-0": "1.1.3.1", "pod-1": "1.1.3.2", "pod-2": "1.1.3.3"},
		controlPods:     map[string]string{"pod-0": "1.1.5.1", "pod-1": "1.1.5.2", "pod-2": "1.1.5.3"},
		kubemanagerPods: map[string]string{"pod-0": "1.1.6.1", "pod-1": "1.1.6.2", "pod-2": "1.1.6.3"},
		webuiPods:       map[string]string{"pod-0": "1.1.7.1", "pod-1": "1.1.7.2", "pod-2": "1.1.7.3"},
		vrouterPods:     map[string]string{"pod-0": "1.1.8.1", "pod-1": "1.1.8.2", "pod-2": "1.1.8.3"},
	}

	podTemplate := corev1.Pod{}

	configPodItems := []corev1.Pod{}
	for pod, ip := range podMap.configPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		podTemplate.SetAnnotations(map[string]string{"hostname": "host1"})
		configPodItems = append(configPodItems, podTemplate)
	}
	configPodList := corev1.PodList{
		Items: configPodItems,
	}
	rabbitmqPodItems := []corev1.Pod{}
	for pod, ip := range podMap.rabbitmqPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		rabbitmqPodItems = append(rabbitmqPodItems, podTemplate)
	}
	rabbitmqPodList := corev1.PodList{
		Items: rabbitmqPodItems,
	}
	cassandraPodItems := []corev1.Pod{}
	for pod, ip := range podMap.cassandraPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		podTemplate.Status.Conditions = []corev1.PodCondition{{
			Type:   corev1.PodReady,
			Status: corev1.ConditionTrue,
		}}
		cassandraPodItems = append(cassandraPodItems, podTemplate)
	}
	cassandraPodList := corev1.PodList{
		Items: cassandraPodItems,
	}
	zookeeperPodItems := []corev1.Pod{}
	for pod, ip := range podMap.zookeeperPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		zookeeperPodItems = append(zookeeperPodItems, podTemplate)
	}
	zookeeperPodList := corev1.PodList{
		Items: zookeeperPodItems,
	}
	controlPodItems := []corev1.Pod{}
	for pod, ip := range podMap.controlPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		podTemplate.SetAnnotations(map[string]string{"hostname": "host1"})
		controlPodItems = append(controlPodItems, podTemplate)
	}
	controlPodList := corev1.PodList{
		Items: controlPodItems,
	}
	kubemanagerPodItems := []corev1.Pod{}
	for pod, ip := range podMap.kubemanagerPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		kubemanagerPodItems = append(kubemanagerPodItems, podTemplate)
	}
	kubemanagerPodList := corev1.PodList{
		Items: kubemanagerPodItems,
	}
	webuiPodItems := []corev1.Pod{}
	for pod, ip := range podMap.webuiPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		webuiPodItems = append(webuiPodItems, podTemplate)
	}
	webuiPodList := corev1.PodList{
		Items: webuiPodItems,
	}

	vrouterPodItems := []corev1.Pod{}
	for pod, ip := range podMap.vrouterPods {
		podTemplate.Name = pod
		podTemplate.Namespace = "default"
		podTemplate.Status.PodIP = ip
		podTemplate.SetAnnotations(map[string]string{"hostname": "host1", "physicalInterface": "eth0", "physicalInterfaceMac": "de:ad:be:ef:ba:be", "prefixLength": "24"})
		vrouterPodItems = append(vrouterPodItems, podTemplate)
	}
	vrouterPodList := corev1.PodList{
		Items: vrouterPodItems,
	}

	configResource.ManageNodeStatus(podMap.configPods, cl)
	configResource.SetEndpointInStatus(cl, configService.Spec.ClusterIP)
	rabbitmqResource.ManageNodeStatus(podMap.rabbitmqPods, cl)

	cassandraResource.ManageNodeStatus(podMap.cassandraPods, cl)
	zookeeperResource.ManageNodeStatus(podMap.zookeeperPods, cl)
	controlResource.ManageNodeStatus(podMap.controlPods, cl)
	kubemanagerResource.ManageNodeStatus(podMap.kubemanagerPods, cl)
	webuiResource.ManageNodeStatus(podMap.webuiPods, cl)
	vrouterResource.ManageNodeStatus(podMap.vrouterPods, cl)

	environment := Environment{
		client:                   &cl,
		configPodList:            configPodList,
		cassandraPodList:         cassandraPodList,
		zookeeperPodList:         zookeeperPodList,
		rabbitmqPodList:          rabbitmqPodList,
		controlPodList:           controlPodList,
		kubemanbagerPodList:      kubemanagerPodList,
		webuiPodList:             webuiPodList,
		vrouterPodList:           vrouterPodList,
		configResource:           *configResource,
		controlResource:          *controlResource,
		cassandraResource:        *cassandraResource,
		zookeeperResource:        *zookeeperResource,
		rabbitmqResource:         *rabbitmqResource,
		kubemanagerResource:      *kubemanagerResource,
		webuiResource:            *webuiResource,
		vrouterResource:          *vrouterResource,
		configConfigMap:          configConfigMap,
		controlConfigMap:         controlConfigMap,
		cassandraConfigMap:       cassandraConfigMap,
		zookeeperConfigMap:       zookeeperConfigMap,
		rabbitmqConfigMap:        rabbitmqConfigMap,
		rabbitmqConfigMap2:       rabbitmqConfigMap2,
		kubemanagerConfigMap:     kubemanagerConfigMap,
		kubemanagerConfigMapEnvs: kubemanagerConfigMapEnvs,
		webuiConfigMap:           webuiConfigMap,
		vrouterConfigMap:         vrouterConfigMap,
		vrouterConfigMap2:        vrouterConfigMap2,
	}
	return environment
}
