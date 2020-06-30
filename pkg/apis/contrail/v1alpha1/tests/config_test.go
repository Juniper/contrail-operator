package contrailtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
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
		},
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
}

var kubemanager = &v1alpha1.Kubemanager{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "kubemanager1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
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
		ServiceConfiguration: v1alpha1.VrouterConfiguration{
			ControlInstance: "control1",
			Gateway:         "1.1.8.254",
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
			ListenPort: 5555,
		},
	},
	Status: v1alpha1.KeystoneStatus{
		IPs: []string{"10.11.12.13"},
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
	client               *client.Client
	configPodList        corev1.PodList
	rabbitmqPodList      corev1.PodList
	zookeeperPodList     corev1.PodList
	cassandraPodList     corev1.PodList
	controlPodList       corev1.PodList
	kubemanbagerPodList  corev1.PodList
	webuiPodList         corev1.PodList
	vrouterPodList       corev1.PodList
	configResource       v1alpha1.Config
	controlResource      v1alpha1.Control
	cassandraResource    v1alpha1.Cassandra
	zookeeperResource    v1alpha1.Zookeeper
	rabbitmqResource     v1alpha1.Rabbitmq
	kubemanagerResource  v1alpha1.Kubemanager
	webuiResource        v1alpha1.Webui
	vrouterResource      v1alpha1.Vrouter
	configConfigMap      corev1.ConfigMap
	controlConfigMap     corev1.ConfigMap
	cassandraConfigMap   corev1.ConfigMap
	zookeeperConfigMap   corev1.ConfigMap
	zookeeperConfigMap2  corev1.ConfigMap
	rabbitmqConfigMap    corev1.ConfigMap
	rabbitmqConfigMap2   corev1.ConfigMap
	kubemanagerConfigMap corev1.ConfigMap
	kubemanagerSecret    corev1.Secret
	webuiConfigMap       corev1.ConfigMap
	vrouterConfigMap     corev1.ConfigMap
	vrouterConfigMap2    corev1.ConfigMap
}

func SetupEnv() Environment {
	logf.SetLogger(logf.ZapLogger(true))
	configConfigMap := *configMap
	rabbitmqConfigMap := *configMap
	rabbitmqConfigMap2 := *configMap
	zookeeperConfigMap := *configMap
	zookeeperConfigMap2 := *configMap
	cassandraConfigMap := *configMap
	controlConfigMap := *configMap
	kubemanagerConfigMap := *configMap
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

	zookeeperConfigMap.Name = "zookeeper1-zookeeper-configmap-1"
	zookeeperConfigMap.Namespace = "default"

	zookeeperConfigMap2.Name = "zookeeper1-zookeeper-configmap"
	zookeeperConfigMap2.Namespace = "default"

	cassandraConfigMap.Name = "cassandra1-cassandra-configmap"
	cassandraConfigMap.Namespace = "default"

	controlConfigMap.Name = "control1-control-configmap"
	controlConfigMap.Namespace = "default"

	kubemanagerConfigMap.Name = "kubemanager1-kubemanager-configmap"
	kubemanagerConfigMap.Namespace = "default"

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
		&zookeeperConfigMap2,
		&rabbitmqConfigMap,
		&rabbitmqConfigMap2,
		&kubemanagerConfigMap,
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
		configPods:      map[string]string{"pod1": "1.1.1.1", "pod2": "1.1.1.2", "pod3": "1.1.1.3"},
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
	rabbitmqResource.ManageNodeStatus(podMap.rabbitmqPods, cl)

	cassandraResource.ManageNodeStatus(podMap.cassandraPods, cl)
	zookeeperResource.ManageNodeStatus(podMap.zookeeperPods, cl)
	controlResource.ManageNodeStatus(podMap.controlPods, cl)
	kubemanagerResource.ManageNodeStatus(podMap.kubemanagerPods, cl)
	webuiResource.ManageNodeStatus(podMap.webuiPods, cl)
	vrouterResource.ManageNodeStatus(podMap.vrouterPods, cl)

	environment := Environment{
		client:               &cl,
		configPodList:        configPodList,
		cassandraPodList:     cassandraPodList,
		zookeeperPodList:     zookeeperPodList,
		rabbitmqPodList:      rabbitmqPodList,
		controlPodList:       controlPodList,
		kubemanbagerPodList:  kubemanagerPodList,
		webuiPodList:         webuiPodList,
		vrouterPodList:       vrouterPodList,
		configResource:       *configResource,
		controlResource:      *controlResource,
		cassandraResource:    *cassandraResource,
		zookeeperResource:    *zookeeperResource,
		rabbitmqResource:     *rabbitmqResource,
		kubemanagerResource:  *kubemanagerResource,
		webuiResource:        *webuiResource,
		vrouterResource:      *vrouterResource,
		configConfigMap:      configConfigMap,
		controlConfigMap:     controlConfigMap,
		cassandraConfigMap:   cassandraConfigMap,
		zookeeperConfigMap:   zookeeperConfigMap,
		zookeeperConfigMap2:  zookeeperConfigMap2,
		rabbitmqConfigMap:    rabbitmqConfigMap,
		rabbitmqConfigMap2:   rabbitmqConfigMap2,
		kubemanagerConfigMap: kubemanagerConfigMap,
		webuiConfigMap:       webuiConfigMap,
		vrouterConfigMap:     vrouterConfigMap,
		vrouterConfigMap2:    vrouterConfigMap2,
	}
	return environment
}
func TestConfigConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	request := reconcile.Request{types.NamespacedName{Name: "config1", Namespace: "default"}}
	configMapNamespacedName := types.NamespacedName{Name: "config1-config-configmap", Namespace: "default"}
	t.Run("default setup", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client
		err := environment.configResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "config1", Namespace: "default"}}, &environment.configPodList, cl)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		err = cl.Get(context.TODO(),
			types.NamespacedName{Name: "config1-config-configmap", Namespace: "default"},
			&environment.configConfigMap)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		if environment.configConfigMap.Data["api.1.1.1.1"] != configConfigHa {
			diff := diff.Diff(environment.configConfigMap.Data["api.1.1.1.1"], configConfigHa)
			t.Fatalf("get api config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["vnc.1.1.1.1"] != vncApiConfig {
			diff := diff.Diff(environment.configConfigMap.Data["vnc.1.1.1.1"], vncApiConfig)
			t.Fatalf("get vncapi config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["devicemanager.1.1.1.1"] != devicemanagerConfig {
			diff := diff.Diff(environment.configConfigMap.Data["devicemanager.1.1.1.1"], devicemanagerConfig)
			t.Fatalf("get devicemanager config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["dnsmasq.1.1.1.1"] != dnsmasqConfig {
			diff := diff.Diff(environment.configConfigMap.Data["dnsmasq.1.1.1.1"], dnsmasqConfig)
			t.Fatalf("get dnsmasq config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["schematransformer.1.1.1.1"] != schematransformerConfig {
			diff := diff.Diff(environment.configConfigMap.Data["schematransformer.1.1.1.1"], schematransformerConfig)
			t.Fatalf("get schematransformer config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["servicemonitor.1.1.1.1"] != servicemonitorConfig {
			diff := diff.Diff(environment.configConfigMap.Data["servicemonitor.1.1.1.1"], servicemonitorConfig)
			t.Fatalf("get servicemonitor config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["analyticsapi.1.1.1.1"] != analyticsapiConfig {
			diff := diff.Diff(environment.configConfigMap.Data["analyticsapi.1.1.1.1"], analyticsapiConfig)
			t.Fatalf("get analyticsapi config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["queryengine.1.1.1.1"] != queryengineConfig {
			diff := diff.Diff(environment.configConfigMap.Data["queryengine.1.1.1.1"], queryengineConfig)
			t.Fatalf("get queryengine config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["collector.1.1.1.1"] != collectorConfig {
			diff := diff.Diff(environment.configConfigMap.Data["collector.1.1.1.1"], collectorConfig)
			t.Fatalf("get collector config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["nodemanagerconfig.1.1.1.1"] != confignodemanagerConfig {
			diff := diff.Diff(environment.configConfigMap.Data["nodemanagerconfig.1.1.1.1"], confignodemanagerConfig)
			t.Fatalf("get nodemanagerconfig config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["nodemanageranalytics.1.1.1.1"] != confignodemanagerAnalytics {
			diff := diff.Diff(environment.configConfigMap.Data["nodemanageranalytics.1.1.1.1"], confignodemanagerAnalytics)
			t.Fatalf("get nodemanageranalytics config: \n%v\n", diff)
		}
	})

	t.Run("device manager host ip is the same as fabric IP stored in config spec", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client
		environment.configResource.Spec.ServiceConfiguration.FabricMgmtIP = "2.2.2.2"

		err := environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl)
		assert.NoError(t, err, "cannot configure instance")

		err = cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap)
		assert.NoError(t, err, "cannot get configmap:")

		actual := environment.configConfigMap.Data["devicemanager.1.1.1.1"]
		assert.Equal(t, actual, devicemanagerWithFabricConfig)
	})
}

func TestControlConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	err := environment.controlResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "control1", Namespace: "default"}}, &environment.controlPodList, cl)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "control1-control-configmap", Namespace: "default"},
		&environment.controlConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.controlConfigMap.Data["control.1.1.5.1"] != controlConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["control.1.1.5.1"], controlConfig)
		t.Fatalf("get control config: \n%v\n", diff)
	}

	if environment.controlConfigMap.Data["named.1.1.5.1"] != namedConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["named.1.1.5.1"], namedConfig)
		t.Fatalf("get named config: \n%v\n", diff)
	}

	if environment.controlConfigMap.Data["dns.1.1.5.1"] != dnsConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["dns.1.1.5.1"], dnsConfig)
		t.Fatalf("get dns config: \n%v\n", diff)
	}

	if environment.controlConfigMap.Data["nodemanager.1.1.5.1"] != controlNodemanagerConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["nodemanager.1.1.5.1"], controlNodemanagerConfig)
		t.Fatalf("get nodemanager config: \n%v\n", diff)
	}

	if environment.controlConfigMap.Data["provision.sh.1.1.5.1"] != controlProvisioningConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["provision.sh.1.1.5.1"], controlProvisioningConfig)
		t.Fatalf("get provision config: \n%v\n", diff)
	}
	if environment.controlConfigMap.Data["deprovision.py.1.1.5.1"] != controlDeProvisioningConfig {
		diff := diff.Diff(environment.controlConfigMap.Data["deprovision.py.1.1.5.1"], controlDeProvisioningConfig)
		t.Fatalf("get deprovision config: \n%v\n", diff)
	}
}

func TestKubemanagerConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	clientset := kubernetes.Clientset{}
	err := environment.kubemanagerResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "kubemanager1", Namespace: "default"}},
		&environment.kubemanbagerPodList, cl, k8s.ClusterConfig{Client: clientset.CoreV1()})
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "kubemanagersecret", Namespace: "default"},
		&environment.kubemanagerSecret)
	if err != nil {
		t.Fatalf("get secret: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "kubemanager1-kubemanager-configmap", Namespace: "default"},
		&environment.kubemanagerConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.kubemanagerConfigMap.Data["kubemanager.1.1.6.1"] != kubemanagerConfig {
		diff := diff.Diff(environment.kubemanagerConfigMap.Data["kubemanager.1.1.6.1"], kubemanagerConfig)
		t.Fatalf("get kubemanager config: \n%v\n", diff)
	}
}

func TestZookeeperConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	err := environment.zookeeperResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "zookeeper1", Namespace: "default"}}, &environment.zookeeperPodList, cl)
	if err != nil {
		t.Fatalf("create config for zookeeper failed: (%v)", err)
	}

	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "zookeeper1-zookeeper-configmap-1", Namespace: "default"},
		&environment.zookeeperConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.zookeeperConfigMap.Data["zoo.cfg"] != zookeeperConfig {
		configDiff := diff.Diff(environment.zookeeperConfigMap.Data["zoo.cfg"], zookeeperConfig)
		t.Fatalf("get zoo.cfg config: \n%v\n", configDiff)
	}

	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "zookeeper1-zookeeper-configmap", Namespace: "default"},
		&environment.zookeeperConfigMap2)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.zookeeperConfigMap2.Data["zoo.cfg.dynamic.100000000"] != zookeeperDynamicConfig {
		configDiff := diff.Diff(environment.zookeeperConfigMap2.Data["zoo.cfg.dynamic.100000000"], zookeeperDynamicConfig)
		t.Fatalf("get zoo.cfg.dynamic.100000000 config: \n%v\n", configDiff)
	}
}

func TestWebuiConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	err := environment.webuiResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "webui1", Namespace: "default"}}, &environment.webuiPodList, cl)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "webui1-webui-configmap", Namespace: "default"},
		&environment.webuiConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.webuiConfigMap.Data["config.global.js.1.1.7.1"] != webuiConfigHa {
		configDiff := diff.Diff(environment.webuiConfigMap.Data["config.global.js.1.1.7.1"], webuiConfigHa)
		t.Fatalf("get webui config: \n%v\n", configDiff)
	}

	if environment.webuiConfigMap.Data["contrail-webui-userauth.js"] != webuiAuthConfig {
		configDiff := diff.Diff(environment.webuiConfigMap.Data["contrail-webui-userauth.js"], webuiAuthConfig)
		t.Fatalf("get webui auth config: \n%v\n", configDiff)
	}
}



func TestVrouterConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	clientset := kubernetes.Clientset{}

	err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, k8s.ClusterConfig{Client: clientset.CoreV1()})
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"},
		&environment.vrouterConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"},
		&environment.vrouterConfigMap2)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.vrouterConfigMap.Data["vrouter.1.1.8.1"] != vrouterConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["vrouter.1.1.8.1"], vrouterConfig)
		t.Fatalf("get vrouter config: \n%v\n", configDiff)
	}
}

func TestCassandraConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	err := environment.cassandraResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "cassandra1", Namespace: "default"}}, &environment.cassandraPodList, cl)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "cassandra1-cassandra-configmap", Namespace: "default"},
		&environment.cassandraConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.cassandraConfigMap.Data["1.1.2.1.yaml"] != cassandraConfig {
		configDiff := diff.Diff(environment.cassandraConfigMap.Data["1.1.2.1.yaml"], cassandraConfig)
		t.Fatalf("get cassandra config: \n%v\n", configDiff)
	}
}

func TestRabbitmqConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	err := environment.rabbitmqResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "rabbitmq1", Namespace: "default"}}, &environment.rabbitmqPodList, cl)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "rabbitmq1-rabbitmq-configmap", Namespace: "default"},
		&environment.rabbitmqConfigMap)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	for _, k := range []string{
		"rabbitmq-1.1.4.1.conf", "rabbitmq-1.1.4.2.conf", "rabbitmq-1.1.4.3.conf", "rabbitmq.nodes",
		"RABBITMQ_ERLANG_COOKIE", "RABBITMQ_USE_LONGNAME", "RABBITMQ_CONFIG_FILE", "RABBITMQ_PID_FILE",
		"RABBITMQ_PID_FILE", "RABBITMQ_CONF_ENV_FILE", "plugins.conf",
		// "definitions.json" TODO: Handle random password_hash in test
	} {
		if configDiff := diff.Diff(environment.rabbitmqConfigMap.Data[k], rabbitmqConfig[k]); configDiff != "" {
			t.Fatalf("get rabbitmq config key = %v: \n%v\n", k, configDiff)
		}
	}

	err = cl.Get(context.TODO(),
		types.NamespacedName{Name: "rabbitmq1-rabbitmq-configmap-runner", Namespace: "default"},
		&environment.rabbitmqConfigMap2)
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.rabbitmqConfigMap2.Data["run.sh"] != rabbitmqConfigRunner {
		configDiff := diff.Diff(environment.rabbitmqConfigMap2.Data["run.sh"], rabbitmqConfigRunner)
		t.Fatalf("get rabbitmq config: \n%v\n", configDiff)
	}
}

var webuiConfigHa = `/*
* Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
*/
var config = {};
config.orchestration = {};
config.orchestration.Manager = "openstack";
config.orchestrationModuleEndPointFromConfig = false;
config.contrailEndPointFromConfig = true;
config.regionsFromConfig = false;
config.endpoints = {};
config.endpoints.apiServiceType = "ApiServer";
config.endpoints.opServiceType = "OpServer";
config.regions = {};
config.regions.RegionOne = "https://10.11.12.13:5555/v3";
config.serviceEndPointTakePublicURL = true;
config.networkManager = {};
config.networkManager.ip = "127.0.0.1";
config.networkManager.port = "9696";
config.networkManager.authProtocol = "http";
config.networkManager.apiVersion = [];
config.networkManager.strictSSL = false;
config.networkManager.ca = "";
config.imageManager = {};
config.imageManager.ip = "127.0.0.1";
config.imageManager.port = "9292";
config.imageManager.authProtocol = "http";
config.imageManager.apiVersion = ['v1', 'v2'];
config.imageManager.strictSSL = false;
config.imageManager.ca = "";
config.computeManager = {};
config.computeManager.ip = "127.0.0.1";
config.computeManager.port = "8774";
config.computeManager.authProtocol = "http";
config.computeManager.apiVersion = ['v1.1', 'v2'];
config.computeManager.strictSSL = false;
config.computeManager.ca = "";
config.identityManager = {};
config.identityManager.ip = "10.11.12.13";
config.identityManager.port = "5555";
config.identityManager.authProtocol = "https";
config.identityManager.apiVersion = ['v3'];
config.identityManager.strictSSL = false;
config.identityManager.defaultDomain = "Default";
config.identityManager.ca = "/etc/ssl/certs/kubernetes/ca-bundle.crt";
config.storageManager = {};
config.storageManager.ip = "127.0.0.1";
config.storageManager.port = "8776";
config.storageManager.authProtocol = "http";
config.storageManager.apiVersion = ['v1'];
config.storageManager.strictSSL = false;
config.storageManager.ca = "";
config.cnfg = {};
config.cnfg.server_ip = ['1.1.1.1','1.1.1.2','1.1.1.3'];
config.cnfg.server_port = "8082";
config.cnfg.authProtocol = "https";
config.cnfg.strictSSL = false;
config.cnfg.ca = "/etc/ssl/certs/kubernetes/ca-bundle.crt";
config.cnfg.statusURL = '/global-system-configs';
config.analytics = {};
config.analytics.server_ip = ['1.1.1.1','1.1.1.2','1.1.1.3'];
config.analytics.server_port = "8081";
config.analytics.authProtocol = "https";
config.analytics.strictSSL = false;
config.analytics.ca = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.analytics.statusURL = '/analytics/uves/bgp-peers';
config.dns = {};
config.dns.server_ip = ['1.1.5.1','1.1.5.2','1.1.5.3'];
config.dns.server_port = '8092';
config.dns.statusURL = '/Snh_PageReq?x=AllEntries%20VdnsServersReq';
config.vcenter = {};
config.vcenter.server_ip = "127.0.0.1";         //vCenter IP
config.vcenter.server_port = "443";                                //Port
config.vcenter.authProtocol = "https";   //http or https
config.vcenter.datacenter = "vcenter";      //datacenter name
config.vcenter.dvsswitch = "vswitch";         //dvsswitch name
config.vcenter.strictSSL = false;                                  //Validate the certificate or ignore
config.vcenter.ca = '';                                            //specify the certificate key file
config.vcenter.wsdl = "/usr/src/contrail/contrail-web-core/webroot/js/vim.wsdl";
config.introspect = {};
config.introspect.ssl = {};
config.introspect.ssl.enabled = true;
config.introspect.ssl.key = '/etc/certificates/server-key-1.1.7.1.pem';
config.introspect.ssl.cert = '/etc/certificates/server-1.1.7.1.crt';
config.introspect.ssl.ca = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.introspect.ssl.strictSSL = false;
config.jobServer = {};
config.jobServer.server_ip = '127.0.0.1';
config.jobServer.server_port = '3000';
config.files = {};
config.files.download_path = '/tmp';
config.cassandra = {};
config.cassandra.server_ips = ['1.1.2.1','1.1.2.2','1.1.2.3'];
config.cassandra.server_port = '9042';
config.cassandra.enable_edit = false;
config.cassandra.use_ssl = true;
config.cassandra.ca_certs = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.kue = {};
config.kue.ui_port = '3002'
config.webui_addresses = {};
config.insecure_access = false;
config.http_port = '8180';
config.https_port = '8143';
config.require_auth = false;
config.node_worker_count = 1;
config.maxActiveJobs = 10;
config.redisDBIndex = 3;
config.CONTRAIL_SERVICE_RETRY_TIME = 300000; //5 minutes
config.redis_server_port = '6380';
config.redis_server_ip = '127.0.0.1';
config.redis_dump_file = '/var/lib/redis/dump-webui.rdb';
config.redis_password = '';
config.logo_file = '/opt/contrail/images/logo.png';
config.favicon_file = '/opt/contrail/images/favicon.ico';
config.featurePkg = {};
config.featurePkg.webController = {};
config.featurePkg.webController.path = '/usr/src/contrail/contrail-web-controller';
config.featurePkg.webController.enable = true;
config.qe = {};
config.qe.enable_stat_queries = false;
config.logs = {};
config.logs.level = 'debug';
config.getDomainProjectsFromApiServer = false;
config.network = {};
config.network.L2_enable = false;
config.getDomainsFromApiServer = false;
config.jsonSchemaPath = "/usr/src/contrail/contrail-web-core/src/serverroot/configJsonSchemas";
config.server_options = {};
config.server_options.key_file = '/etc/certificates/server-key-1.1.7.1.pem';
config.server_options.cert_file = '/etc/certificates/server-1.1.7.1.crt';
config.server_options.ciphers = 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES256-SHA';
module.exports = config;
config.staticAuth = [];
config.staticAuth[0] = {};
config.staticAuth[0].username = 'admin';
config.staticAuth[0].password = 'test123';
config.staticAuth[0].roles = ['cloudAdmin'];
`

var webuiAuthConfig = `/*
* Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
*/
var auth = {};
auth.admin_user = 'admin';
auth.admin_password = 'test123';
auth.admin_token = '';
auth.admin_tenant_name = 'admin';
auth.project_domain_name = 'Default';
auth.user_domain_name = 'Default';
module.exports = auth;
`

var configConfigHa = `[DEFAULTS]
listen_ip_addr=0.0.0.0
listen_port=8082
http_server_port=8084
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-api.log
log_level=SYS_NOTICE
log_local=1
list_optimization_enabled=True
auth=keystone
aaa_mode=rbac
cloud_admin_role=admin
global_read_only_role=
config_api_ssl_enable=True
config_api_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
config_api_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
config_api_ssl_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var cassandraConfig = `cluster_name: ContrailConfigDB
num_tokens: 32
hinted_handoff_enabled: true
max_hint_window_in_ms: 10800000 # 3 hours
hinted_handoff_throttle_in_kb: 1024
max_hints_delivery_threads: 2
hints_directory: /var/lib/cassandra/hints
hints_flush_period_in_ms: 10000
max_hints_file_size_in_mb: 128
batchlog_replay_throttle_in_kb: 1024
authenticator: AllowAllAuthenticator
authorizer: AllowAllAuthorizer
role_manager: CassandraRoleManager
roles_validity_in_ms: 2000
permissions_validity_in_ms: 2000
credentials_validity_in_ms: 2000
partitioner: org.apache.cassandra.dht.Murmur3Partitioner
data_file_directories:
- /var/lib/cassandra/data
commitlog_directory: /var/lib/cassandra/commitlog
disk_failure_policy: stop
commit_failure_policy: stop
key_cache_size_in_mb:
key_cache_save_period: 14400
row_cache_size_in_mb: 0
row_cache_save_period: 0
counter_cache_size_in_mb:
counter_cache_save_period: 7200
saved_caches_directory: /var/lib/cassandra/saved_caches
commitlog_sync: periodic
commitlog_sync_period_in_ms: 10000
commitlog_segment_size_in_mb: 32
seed_provider:
- class_name: org.apache.cassandra.locator.SimpleSeedProvider
  parameters:
  - seeds: 1.1.2.1,1.1.2.2
concurrent_reads: 32
concurrent_writes: 32
concurrent_counter_writes: 32
concurrent_materialized_view_writes: 32
disk_optimization_strategy: ssd
memtable_allocation_type: heap_buffers
index_summary_capacity_in_mb:
index_summary_resize_interval_in_minutes: 60
trickle_fsync: false
trickle_fsync_interval_in_kb: 10240
storage_port: 7000
ssl_storage_port: 7001
listen_address: 1.1.2.1
broadcast_address: 1.1.2.1
start_native_transport: true
native_transport_port: 9042
start_rpc: true
rpc_address: 1.1.2.1
rpc_port: 9160
broadcast_rpc_address: 1.1.2.1
rpc_keepalive: true
rpc_server_type: sync
thrift_framed_transport_size_in_mb: 15
incremental_backups: false
snapshot_before_compaction: false
auto_snapshot: true
tombstone_warn_threshold: 1000
tombstone_failure_threshold: 100000
column_index_size_in_kb: 64
batch_size_warn_threshold_in_kb: 5
batch_size_fail_threshold_in_kb: 50
compaction_throughput_mb_per_sec: 16
compaction_large_partition_warning_threshold_mb: 100
sstable_preemptive_open_interval_in_mb: 50
read_request_timeout_in_ms: 5000
range_request_timeout_in_ms: 10000
write_request_timeout_in_ms: 2000
counter_write_request_timeout_in_ms: 5000
cas_contention_timeout_in_ms: 1000
truncate_request_timeout_in_ms: 60000
request_timeout_in_ms: 10000
cross_node_timeout: false
endpoint_snitch: SimpleSnitch
dynamic_snitch_update_interval_in_ms: 100
dynamic_snitch_reset_interval_in_ms: 600000
dynamic_snitch_badness_threshold: 0.1
request_scheduler: org.apache.cassandra.scheduler.NoScheduler
server_encryption_options:
  internode_encryption: all
  keystore: /etc/keystore/server-keystore.jks
  keystore_password: keystorePassword
  truststore: /etc/keystore/server-truststore.jks
  truststore_password: truststorePassword
  require_client_auth: true
  store_type: JKS
client_encryption_options:
  enabled: true
  optional: false
  keystore: /etc/keystore/server-keystore.jks
  keystore_password: keystorePassword
internode_compression: all
inter_dc_tcp_nodelay: false
tracetype_query_ttl: 86400
tracetype_repair_ttl: 604800
gc_warn_threshold_in_ms: 1000
enable_user_defined_functions: false
enable_scripted_user_defined_functions: false
windows_timer_interval: 1
transparent_data_encryption_options:
  enabled: false
  chunk_length_kb: 64
  cipher: AES/CBC/PKCS5Padding
  key_alias: testing:1
  key_provider:
  - class_name: org.apache.cassandra.security.JKSKeyProvider
    parameters:
    - keystore: conf/.keystore
      keystore_password: cassandra
      store_type: JCEKS
      key_password: cassandra
auto_bootstrap: true
`

var zookeeperConfig = `clientPort=2181
clientPortAddress=
dataDir=/var/lib/zookeeper
tickTime=2000
initLimit=5
syncLimit=2
maxClientCnxns=60
admin.enableServer=true
standaloneEnabled=false
4lw.commands.whitelist=stat,ruok,conf,isro
reconfigEnabled=true
dynamicConfigFile=/mydata/zoo.cfg.dynamic.100000000
`

var zookeeperDynamicConfig = `server.1=1.1.3.1:2888:3888:participant
server.2=1.1.3.2:2888:3888:participant
server.3=1.1.3.3:2888:3888:participant
`

var rabbitmqConfigRunner = `#!/bin/bash
echo $RABBITMQ_ERLANG_COOKIE > /var/lib/rabbitmq/.erlang.cookie
chmod 0600 /var/lib/rabbitmq/.erlang.cookie
export RABBITMQ_NODENAME=rabbit@${POD_IP}
if [[ $(grep $POD_IP /etc/rabbitmq/0) ]] ; then
  rabbitmq-server
else
  rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) ping
  while [[ $? -ne 0 ]]; do
	rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) ping
  done
  rabbitmq-server -detached
  rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) node_health_check
  while [[ $? -ne 0 ]]; do
	rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) node_health_check
  done
  rabbitmqctl stop_app
  sleep 2
  rabbitmqctl join_cluster rabbit@$(cat /etc/rabbitmq/0)
  rabbitmqctl shutdown
  rabbitmq-server
fi
`

var rabbitmqConfig = map[string]string{
	"rabbitmq-1.1.4.1.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.1.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.1.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
`,
	"rabbitmq-1.1.4.2.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.2.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.2.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
`,
	"rabbitmq-1.1.4.3.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.3.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.3.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
`,
	"rabbitmq.nodes":         fmt.Sprintf("1.1.4.1\n1.1.4.2\n1.1.4.3\n"),
	"0":                      "1.1.4.1",
	"1":                      "1.1.4.2",
	"2":                      "1.1.4.3",
	"RABBITMQ_ERLANG_COOKIE": "47EFF3BB-4786-46E0-A5BB-58455B3C2CB4",
	"RABBITMQ_USE_LONGNAME":  "true",
	"RABBITMQ_CONFIG_FILE":   "/etc/rabbitmq/rabbitmq-${POD_IP}.conf",
	"RABBITMQ_PID_FILE":      "/var/run/rabbitmq.pid",
	"RABBITMQ_CONF_ENV_FILE": "/var/lib/rabbitmq/rabbitmq.env",
	"definitions.json":       rabbitmqDefinition,
	"plugins.conf":           "[rabbitmq_management,rabbitmq_management_agent].",
}

var rabbitmqDefinition = `{
  "users": [
    {
      "name": "user",
      "password_hash": "0TeE2b17AedPdBgaZSwAtAfgwXDSqqoDki44i3AcXqxU1DIB",
      "tags": "administrator"
    }
  ],
  "vhosts": [
    {
      "name": "vhost"
    }
  ],
  "permissions": [
    {
      "user": "user",
      "vhost": "vhost",
      "configure": ".*",
      "write": ".*",
      "read": ".*"
    }
  ],
}
`

var devicemanagerConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
http_server_port=8096
api_server_use_ssl=True
analytics_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/contrail-device-manager.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/var/lib/dnsmasq
tftp_dir=/var/lib/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
dnsmasq_reload_by_signal=True
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
dm_run_mode=Full
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var dnsmasqConfig = `
log-facility=/dev/stdout
bogus-priv
log-dhcp
enable-tftp
tftp-root=/etc/tftp
dhcp-leasefile=/var/lib/dnsmasq/dnsmasq.leases
conf-dir=/var/lib/dnsmasq/,*.conf
`
var schematransformerConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
http_server_port=8087
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-schema.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
[SECURITY]
use_certs=True
ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
certfile=/etc/certificates/server-1.1.1.1.crt
keyfile=/etc/certificates/server-key-1.1.1.1.pem`

var servicemonitorConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
http_server_port=8088
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-svc-monitor.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
analytics_api_ssl_enable = True
analytics_api_insecure_enable = False
analytics_api_ssl_certfile = /etc/certificates/server-1.1.1.1.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-1.1.1.1.pem
analytics_api_ssl_ca_cert = /etc/ssl/certs/kubernetes/ca-bundle.crt
[SECURITY]
use_certs=True
keyfile=/etc/certificates/server-key-1.1.1.1.pem
certfile=/etc/certificates/server-1.1.1.1.crt
ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[SCHEDULER]
# Analytics server list used to get vrouter status and schedule service instance
analytics_server_list=1.1.1.1:8081 1.1.1.2:8081 1.1.1.3:8081
aaa_mode=rbac
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var queryengineConfig = `[DEFAULT]
analytics_data_ttl=48
hostip=1.1.1.1
hostname=host1
http_server_ip=0.0.0.0
http_server_port=8091
log_file=/var/log/contrail/contrail-query-engine.log
log_level=SYS_DEBUG
log_local=1
max_slice=100
max_tasks=16
start_time=0
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
cassandra_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[REDIS]
server_list=1.1.1.1:6379 1.1.1.2:6379 1.1.1.3:6379
password=
redis_ssl_enable=False
[SANDESH]
introspect_ssl_enable=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var analyticsapiConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_port=8090
http_server_ip=0.0.0.0
rest_api_port=8081
rest_api_ip=0.0.0.0
aaa_mode=rbac
log_file=/var/log/contrail/contrail-analytics-api.log
log_level=SYS_NOTICE
log_local=1
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
#sandesh_send_rate_limit =
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
api_server=1.1.1.1:8082 1.1.1.2:8082 1.1.1.3:8082
api_server_use_ssl=True
zk_list=1.1.3.1:2181 1.1.3.2:2181 1.1.3.3:2181
analytics_api_ssl_enable = True
analytics_api_insecure_enable = True
analytics_api_ssl_certfile = /etc/certificates/server-1.1.1.1.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-1.1.1.1.pem
analytics_api_ssl_ca_cert = /etc/ssl/certs/kubernetes/ca-bundle.crt
[REDIS]
redis_uve_list=1.1.1.1:6379 1.1.1.2:6379 1.1.1.3:6379
redis_password=
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var collectorConfig = `[DEFAULT]
analytics_data_ttl=48
analytics_config_audit_ttl=2160
analytics_statistics_ttl=168
analytics_flow_ttl=2
partitions=30
hostip=1.1.1.1
hostname=host1
http_server_port=8089
http_server_ip=0.0.0.0
syslog_port=514
sflow_port=6343
ipfix_port=4739
# log_category=
log_file=/var/log/contrail/contrail-collector.log
log_files_count=10
log_file_size=1048576
log_level=SYS_NOTICE
log_local=1
# sandesh_send_rate_limit=
cassandra_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
zookeeper_server_list=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[COLLECTOR]
port=8086
server=0.0.0.0
protobuf_port=3333
[STRUCTURED_SYSLOG_COLLECTOR]
# TCP & UDP port to listen on for receiving structured syslog messages
port=3514
# List of external syslog receivers to forward structured syslog messages in ip:port format separated by space
# tcp_forward_destination=10.213.17.53:514
[API_SERVER]
# List of api-servers in ip:port format separated by space
api_server_list=1.1.1.1:8082 1.1.1.2:8082 1.1.1.3:8082
api_server_use_ssl=True
[REDIS]
port=6379
server=127.0.0.1
password=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
config_db_use_ssl=true
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
rabbitmq_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var confignodemanagerConfig = `[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=`

var confignodemanagerAnalytics = `[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=`

var controlConfig = `[DEFAULT]
# bgp_config_file=bgp_config.xml
bgp_port=179
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
# gr_helper_bgp_disable=0
# gr_helper_xmpp_disable=0
hostip=1.1.5.1
hostname=host1
http_server_ip=0.0.0.0
http_server_port=8083
log_file=/var/log/contrail/contrail-control.log
log_level=SYS_NOTICE
log_local=1
# log_files_count=10
# log_file_size=10485760 # 10MB
# log_category=
# log_disable=0
xmpp_server_port=5269
xmpp_auth_enable=True
xmpp_server_cert=/etc/certificates/server-1.1.5.1.crt
xmpp_server_key=/etc/certificates/server-key-1.1.5.1.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt

# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
# config_db_username=
# config_db_password=
config_db_use_ssl=True
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-1.1.5.1.pem
rabbitmq_ssl_certfile=/etc/certificates/server-1.1.5.1.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.5.1.pem
sandesh_certfile=/etc/certificates/server-1.1.5.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var dnsConfig = `[DEFAULT]
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
named_config_file = /etc/mycontrail/named.1.1.5.1
named_config_directory = /etc/contrail/dns
named_log_file = /var/log/contrail/contrail-named.log
rndc_config_file = contrail-rndc.conf
named_max_cache_size=32M # max-cache-size (bytes) per view, can be in K or M
named_max_retransmissions=12
named_retransmission_interval=1000 # msec
hostip=1.1.5.1
hostname=host1
http_server_port=8092
http_server_ip=0.0.0.0
dns_server_port=53
log_file=/var/log/contrail/contrail-dns.log
log_level=SYS_NOTICE
log_local=1
# log_files_count=10
# log_file_size=10485760 # 10MB
# log_category=
# log_disable=0
xmpp_dns_auth_enable=True
xmpp_server_cert=/etc/certificates/server-1.1.5.1.crt
xmpp_server_key=/etc/certificates/server-key-1.1.5.1.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
# config_db_username=
# config_db_password=
config_db_use_ssl=True
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-1.1.5.1.pem
rabbitmq_ssl_certfile=/etc/certificates/server-1.1.5.1.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.5.1.pem
sandesh_certfile=/etc/certificates/server-1.1.5.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var controlNodemanagerConfig = `[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-control-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip=1.1.5.1
db_port=9042
db_jmx_port=7200
db_use_ssl=True
[COLLECTOR]
server_list=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.5.1.pem
sandesh_certfile=/etc/certificates/server-1.1.5.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var namedConfig = `options {
    directory "/etc/contrail/dns";
    managed-keys-directory "/etc/contrail/dns";
    empty-zones-enable no;
    pid-file "/etc/contrail/dns/contrail-named.pid";
    session-keyfile "/etc/contrail/dns/session.key";
    listen-on port 53 { any; };
    allow-query { any; };
    allow-recursion { any; };
    allow-query-cache { any; };
    max-cache-size 32M;
};
key "rndc-key" {
    algorithm hmac-md5;
    secret "xvysmOR8lnUQRBcunkC6vg==";
};
controls {
    inet 127.0.0.1 port 8094
    allow { 127.0.0.1; }  keys { "rndc-key"; };
};
logging {
    channel debug_log {
        file "/var/log/contrail/contrail-named.log" versions 3 size 5m;
        severity debug;
        print-time yes;
        print-severity yes;
        print-category yes;
    };
    category default {
        debug_log;
    };
    category queries {
        debug_log;
    };
};`

var controlProvisioningConfig = `#!/bin/bash
sed "s/hostip=.*/hostip=${POD_IP}/g" /etc/mycontrail/nodemanager.${POD_IP} > /etc/contrail/contrail-control-nodemgr.conf
servers=$(echo 1.1.1.1,1.1.1.2,1.1.1.3 | tr ',' ' ')
for server in $servers ; do
  python /opt/contrail/utils/provision_control.py --oper $1 \
  --api_server_use_ssl true \
  --host_ip 1.1.5.1 \
  --router_asn 64512 \
  --bgp_server_port 179 \
  --api_server_ip $server \
  --api_server_port 8082 \
  --host_name host1
  if [[ $? -eq 0 ]]; then
	break
  fi
done
`

var controlDeProvisioningConfig = `#!/usr/bin/python
from vnc_api import vnc_api
import socket
vncServerList = ['1.1.1.1','1.1.1.2','1.1.1.3']
vnc_client = vnc_api.VncApi(
            username = 'admin',
            password = 'contrail123',
            tenant_name = 'admin',
            api_server_host= vncServerList[0],
            api_server_port=8082)
vnc_client.bgp_router_delete(fq_name=['default-domain','default-project','ip-fabric','__default__', 'host1' ])
`

var kubemanagerConfig = `[DEFAULTS]
host_ip=1.1.6.1
orchestrator=kubernetes
token=THISISATOKEN
log_file=/var/log/contrail/contrail-kube-manager.log
log_level=SYS_DEBUG
log_local=1
nested_mode=0
http_server_ip=0.0.0.0
[KUBERNETES]
kubernetes_api_server=10.96.0.1
kubernetes_api_port=8080
kubernetes_api_secure_port=6443
cluster_name=kubernetes
cluster_project={}
cluster_network={}
pod_subnets=10.32.0.0/12
ip_fabric_subnets=10.64.0.0/12
service_subnets=10.96.0.0/12
ip_fabric_forwarding=false
ip_fabric_snat=true
host_network_service=false
[VNC]
public_fip_pool={}
vnc_endpoint_ip=1.1.1.1,1.1.1.2,1.1.1.3
vnc_endpoint_port=8082
rabbit_server=1.1.4.1,1.1.4.2,1.1.4.3
rabbit_port=15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.6.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.6.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=True
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.6.1.pem
sandesh_certfile=/etc/certificates/server-1.1.6.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var vrouterConfig = `[CONTROL-NODE]
servers=1.1.5.1:5269 1.1.5.2:5269 1.1.5.3:5269
[DEFAULT]
http_server_ip=0.0.0.0
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
log_file=/var/log/contrail/contrail-vrouter-agent.log
log_level=SYS_NOTICE
log_local=1
hostname=host1
agent_name=host1
xmpp_dns_auth_enable=True
xmpp_auth_enable=True
xmpp_server_cert=/etc/certificates/server-1.1.8.1.crt
xmpp_server_key=/etc/certificates/server-key-1.1.8.1.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
physical_interface_mac = de:ad:be:ef:ba:be
tsn_servers = []
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.8.1.pem
sandesh_certfile=/etc/certificates/server-1.1.8.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
[NETWORKS]
control_network_ip=1.1.8.1
[DNS]
servers=1.1.5.1:53 1.1.5.2:53 1.1.5.3:53
[METADATA]
metadata_proxy_secret=contrail
[VIRTUAL-HOST-INTERFACE]
name=vhost0
ip=1.1.8.1/24
physical_interface=eth0
compute_node_address=1.1.8.1
gateway=1.1.8.254
[SERVICE-INSTANCE]
netns_command=/usr/bin/opencontrail-vrouter-netns
docker_command=/usr/bin/opencontrail-vrouter-docker
[HYPERVISOR]
type = kvm
[FLOWS]
fabric_snat_hash_table_size = 4096
[SESSION]
slo_destination = collector
sample_destination = collector`

var devicemanagerWithFabricConfig = `[DEFAULTS]
host_ip=2.2.2.2
http_server_ip=0.0.0.0
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
http_server_port=8096
api_server_use_ssl=True
analytics_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/contrail-device-manager.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/var/lib/dnsmasq
tftp_dir=/var/lib/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
dnsmasq_reload_by_signal=True
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
dm_run_mode=Full
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var vncApiConfig = `[global]
WEB_SERVER = 1.1.1.1
WEB_PORT = 8082 ; connection to api-server directly
BASE_URL = /
use_ssl = True
cafile = /etc/ssl/certs/kubernetes/ca-bundle.crt
; Authentication settings (optional)
[auth]
AUTHN_TYPE = keystone
AUTHN_PROTOCOL = https
AUTHN_SERVER = 10.11.12.13
AUTHN_PORT = 5555
AUTHN_URL = /v3/auth/tokens
AUTHN_DOMAIN = Default
cafile = /etc/ssl/certs/kubernetes/ca-bundle.crt
;AUTHN_TOKEN_URL = http://127.0.0.1:35357/v2.0/tokens
`
