package ha

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

var versionMap = map[string]string{
	"cassandra":                     "3.11.3",
	"zookeeper":                     "3.5.4-beta",
	"cemVersion":                    "2005.11",
	"python":                        "3.8.2-alpine",
	"redis":                         "4.0.2",
	"busybox":                       "1.31",
	"rabbitmq":                      "3.7",
	"contrail-statusmonitor":        scmBranch + "." + scmRevision,
	"contrail-operator-provisioner": scmBranch + "." + scmRevision,
}

var targetVersionMap = map[string]string{
	"cassandra":                     "3.11.4",
	"zookeeper":                     "3.5.5",
	"cemVersion":                    "2005.42",
	"python":                        "3.8.2-alpine",
	"redis":                         "4.0.2",
	"busybox":                       "1.31",
	"rabbitmq":                      "3.7.16",
	"contrail-statusmonitor":        "R2005.latest",
	"contrail-operator-provisioner": "R2005.latest",
}

func TestHACoreContrailServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true

		cassandras := []*contrail.Cassandra{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-cassandra",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
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

		w := wait.Wait{
			Namespace:     namespace,
			Timeout:       waitTimeout,
			RetryInterval: retryInterval,
			KubeClient:    f.KubeClient,
			Logger:        log,
		}

		t.Run("when manager resource with Config and dependencies are created", func(t *testing.T) {

			var replicas int32 = 1
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})

			require.NoError(t, err)

			t.Run("then a ready Zookeeper StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset", replicas))
			})

			t.Run("then a ready Cassandra StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset", replicas))
			})

			t.Run("then a ready Rabbitmq StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-rabbitmq-rabbitmq-statefulset", replicas))
			})

			t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", replicas))
			})

			t.Run("then a ready WebUI StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", replicas))
			})

			t.Run("then a ready ProvisionManager StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", replicas))
			})
		})

		t.Run("when replicas is set to 3 in manager", func(t *testing.T) {
			t.Skip()
			var replicas int32 = 3
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})

			require.NoError(t, err)

			t.Run("then all services are scaled up from 1 to 3 node", func(t *testing.T) {
				t.Run("then Zookeeper StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset", replicas))
				})

				t.Run("then Cassandra StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset", replicas))
				})

				t.Run("then Rabbitmq StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-rabbitmq-rabbitmq-statefulset", replicas))
				})

				t.Run("then Config StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", replicas))
				})

				t.Run("then WebUI StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", replicas))
				})

				t.Run("then ProvisionManager StatefulSet should be scaled and ready", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", replicas))
				})
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
		t.Run("when manager resource is upgraded", func(t *testing.T) {
			var replicas int32 = 1
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)

			zkContainer := utils.GetContainerFromList("zookeeper", instance.Spec.Services.Zookeepers[0].Spec.ServiceConfiguration.Containers)
			zkContainer.Image = "registry:5000/common-docker-third-party/contrail/zookeeper:" + targetVersionMap["zookeeper"]

			rmqContainer := utils.GetContainerFromList("rabbitmq", instance.Spec.Services.Rabbitmq.Spec.ServiceConfiguration.Containers)
			rmqContainer.Image = "registry:5000/common-docker-third-party/contrail/rabbitmq:" + targetVersionMap["rabbitmq"]

			csContainer := utils.GetContainerFromList("cassandra", instance.Spec.Services.Cassandras[0].Spec.ServiceConfiguration.Containers)
			csContainer.Image = "registry:5000/common-docker-third-party/contrail/cassandra:" + targetVersionMap["cassandra"]

			apiContainer := utils.GetContainerFromList("api", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			devicemanagerContainer := utils.GetContainerFromList("devicemanager", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			dnsmasqContainer := utils.GetContainerFromList("dnsmasq", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			schematransformerContainer := utils.GetContainerFromList("schematransformer", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			servicemonitorContainer := utils.GetContainerFromList("servicemonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			analyticsapiContainer := utils.GetContainerFromList("analyticsapi", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			collectorContainer := utils.GetContainerFromList("collector", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			queryengineContainer := utils.GetContainerFromList("queryengine", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			//Status monitor upgrade is skipped since this fails.
			//TODO: Uncomment this after fixing Statusmonitor issues
			//statusmonitorContainer := utils.GetContainerFromList("statusmonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			apiContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-api:" + targetVersionMap["cemVersion"]
			devicemanagerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + targetVersionMap["cemVersion"]
			dnsmasqContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + targetVersionMap["cemVersion"]
			schematransformerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-schema:" + targetVersionMap["cemVersion"]
			servicemonitorContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + targetVersionMap["cemVersion"]
			analyticsapiContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-api:" + targetVersionMap["cemVersion"]
			collectorContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-collector:" + targetVersionMap["cemVersion"]
			queryengineContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + targetVersionMap["cemVersion"]
			//statusmonitorContainer.Image = "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + intendedVersionMap["contrail-statusmonitor"]

			webuijobContainer := utils.GetContainerFromList("webuijob", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuiwebContainer := utils.GetContainerFromList("webuiweb", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuijobContainer.Image = "registry:5000/contrail-nightly/contrail-controller-webui-job:" + targetVersionMap["cemVersion"]
			webuiwebContainer.Image = "registry:5000/contrail-nightly/contrail-controller-webui-web:" + targetVersionMap["cemVersion"]

			pmContainer := utils.GetContainerFromList("provisioner", instance.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.Containers)
			pmContainer.Image = "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + targetVersionMap["contrail-operator-provisioner"]

			err = f.Client.Update(context.TODO(), instance)
			require.NoError(t, err)

			t.Run("then Zookeeper has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "zookeeper", zkContainer.Image, "zookeeper")
				require.NoError(t, err)
			})

			t.Run("then Zookeeper StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset", replicas))
			})

			t.Run("then Rabbitmq has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "rabbitmq", rmqContainer.Image, "rabbitmq")
				require.NoError(t, err)
			})

			t.Run("then Rabbitmq StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-rabbitmq-rabbitmq-statefulset", replicas))
			})

			t.Run("then Cassandra has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "cassandra", csContainer.Image, "cassandra")
				require.NoError(t, err)
			})

			t.Run("then Cassandra StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset", replicas))
			})

			t.Run("then Config has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "config", apiContainer.Image, "api")
				require.NoError(t, err)
			})

			t.Run("then Config StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", replicas))
			})

			t.Run("then Config pod can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "config=hatest-config",
				})
				assert.NoError(t, err)
				require.NotEmpty(t, configPods.Items)

				configProxy := proxy.NewSecureClient("contrail", configPods.Items[0].Name, 8082)
				req, err := configProxy.NewRequest(http.MethodGet, "/projects", nil)
				assert.NoError(t, err)
				res, err := configProxy.Do(req)
				assert.NoError(t, err)
				assert.Equal(t, 200, res.StatusCode)

			})

			t.Run("then Webui has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "webui", webuijobContainer.Image, "webuijob")
				require.NoError(t, err)
			})

			t.Run("then WebUI StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", replicas))
			})

			t.Run("then ProvisionManager has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "provisionmanager", pmContainer.Image, "provisioner")
				require.NoError(t, err)
			})

			t.Run("then ProvisionManager StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", replicas))
			})

		})
		t.Run("when reference cluster is deleted", func(t *testing.T) {
			pp := meta.DeletePropagationForeground
			err = f.Client.Delete(context.TODO(), cluster, &client.DeleteOptions{
				PropagationPolicy: &pp,
			})
			assert.NoError(t, err)

			t.Run("then manager is cleared in less then 5 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})
		})
	})
}
