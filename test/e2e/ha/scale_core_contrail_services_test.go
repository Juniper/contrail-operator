package ha

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/test/e2e"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

var scmRevision = getEnv("BUILD_SCM_REVISION", "latest")
var scmBranch = getEnv("BUILD_SCM_BRANCH", "master")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var versionMap = map[string]string{
	"cassandra":                     "3.11.4",
	"zookeeper":                     "3.5.5",
	"cemVersion":                    "2005.42",
	"python":                        "3.8.2-alpine",
	"redis":                         "4.0.2",
	"busybox":                       "1.31",
	"rabbitmq":                      "3.7",
	"contrail-statusmonitor":        scmBranch + "." + scmRevision,
	"contrail-operator-provisioner": scmBranch + "." + scmRevision,
}

func TestHACoreContrailServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: e2e.CleanupTimeout, RetryInterval: e2e.CleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, e2e.RetryInterval, e2e.WaitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true
		oneVal := int32(1)

		cassandras := []*contrail.Cassandra{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-cassandra",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
					},
				},
			},
		}}

		zookeepers := []*contrail.Zookeeper{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-zookeeper",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:" + versionMap["zookeeper"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
					},
				},
			},
		}}

		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-rabbitmq",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.RabbitmqSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:" + versionMap["rabbitmq"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
					},
				},
			},
		}

		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-config",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ConfigSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.ConfigConfiguration{
					CassandraInstance: "hatest-cassandra",
					ZookeeperInstance: "hatest-zookeeper",
					Containers: []*contrail.Container{
						{Name: "api", Image: "registry:5000/contrail-nightly/contrail-controller-config-api:" + versionMap["cemVersion"]},
						{Name: "devicemanager", Image: "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + versionMap["cemVersion"]},
						{Name: "dnsmasq", Image: "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + versionMap["cemVersion"]},
						{Name: "schematransformer", Image: "registry:5000/contrail-nightly/contrail-controller-config-schema:" + versionMap["cemVersion"]},
						{Name: "servicemonitor", Image: "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + versionMap["cemVersion"]},
						{Name: "analyticsapi", Image: "registry:5000/contrail-nightly/contrail-analytics-api:" + versionMap["cemVersion"]},
						{Name: "collector", Image: "registry:5000/contrail-nightly/contrail-analytics-collector:" + versionMap["cemVersion"]},
						{Name: "queryengine", Image: "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + versionMap["cemVersion"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + versionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + versionMap["redis"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
						{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
					},
				},
			},
		}

		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-webui",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					CassandraInstance: "hatest-cassandra",
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + versionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + versionMap["redis"]},
						{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + versionMap["cemVersion"]},
						{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + versionMap["cemVersion"]},
					},
				},
			},
		}

		provisionManager := &contrail.ProvisionManager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-provmanager",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ProvisionManagerSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					Replicas:     &oneVal,
				},
				ServiceConfiguration: contrail.ProvisionManagerConfiguration{
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "provisioner", Image: "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + versionMap["contrail-operator-provisioner"]},
					},
				},
			},
		}

		cluster := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cluster1",
				Namespace: namespace,
			},
			Spec: contrail.ManagerSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Replicas:    &oneVal,
					HostNetwork: &trueVal,
				},
				Services: contrail.Services{
					Cassandras:       cassandras,
					Zookeepers:       zookeepers,
					Config:           config,
					Webui:            webui,
					Rabbitmq:         rabbitmq,
					ProvisionManager: provisionManager,
				},
			},
		}

		t.Run("when manager resource with Config and dependencies is created", func(t *testing.T) {
			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: e2e.CleanupTimeout, RetryInterval: e2e.CleanupRetryInterval})
			assert.NoError(t, err)

			w := wait.Wait{
				Namespace:     namespace,
				Timeout:       e2e.WaitTimeout,
				RetryInterval: e2e.RetryInterval,
				KubeClient:    f.KubeClient,
				Logger:        log,
			}

			t.Run("then a ready Zookeeper StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset"))
			})

			t.Run("then a ready Cassandra StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset"))
			})

			t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset"))
			})

			t.Run("then a ready webui StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset"))
			})

			t.Run("then a ready provisionmanager StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset"))
			})

			t.Run("and when zookeeper is scaled up from 1 to 3 node", func(t *testing.T) {
				zookeeperInstance := &contrail.Zookeeper{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-zookeeper", Namespace: namespace}, zookeeperInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				zookeeperInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), zookeeperInstance)
				assert.NoError(t, err)

				t.Run("then a ready Zookeeper StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset"))
				})

			})

			t.Run("and when Cassandra is scaled up from 1 to 3 node", func(t *testing.T) {
				cassandraInstance := &contrail.Cassandra{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-cassandra", Namespace: namespace}, cassandraInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				cassandraInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), cassandraInstance)
				assert.NoError(t, err)

				t.Run("then a ready Cassandra StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset"))
				})

			})

			t.Run("and when config is scaled up from 1 to 3 node", func(t *testing.T) {
				configInstance := &contrail.Config{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-config", Namespace: namespace}, configInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				configInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), configInstance)
				assert.NoError(t, err)

				t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset"))
				})

				t.Run("then all Config pods can process requests", func(t *testing.T) {
					configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
						LabelSelector: "config=hatest-config",
					})
					assert.NoError(t, err)
					require.NotEmpty(t, configPods.Items)

					for _, pod := range configPods.Items {
						configProxy := proxy.NewSecureClient("contrail", pod.Name, 8082)
						req, err := configProxy.NewRequest(http.MethodGet, "/projects", nil)
						assert.NoError(t, err)
						res, err := configProxy.Do(req)
						assert.NoError(t, err)
						assert.Equal(t, 200, res.StatusCode)
					}
				})

			})

			t.Run("and when webui is scaled up from 1 to 3 node", func(t *testing.T) {
				webuiInstance := &contrail.Webui{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-webui", Namespace: namespace}, webuiInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				webuiInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), webuiInstance)
				assert.NoError(t, err)

				t.Run("then a ready webui StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset"))
				})

			})

			t.Run("and when provisionmanager is scaled up from 1 to 3 node", func(t *testing.T) {
				provisionManagerInstance := &contrail.ProvisionManager{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-provmanager", Namespace: namespace}, provisionManagerInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				provisionManagerInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), provisionManagerInstance)
				assert.NoError(t, err)

				t.Run("then a ready provisionmanager StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset"))
				})

			})

		})

	})
}
