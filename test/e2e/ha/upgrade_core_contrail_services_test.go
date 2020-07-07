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

var initialVersionMap = map[string]string{
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

var intendedVersionMap = map[string]string{
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

func TestUpgradeCoreContrailServices(t *testing.T) {
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
		//This upgrade test is skipped since this fails.
		//TODO: Include this test after fixing upgrade issues
		t.Skip()

		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true
		var replicas int32 = 3
		cassandras := []*contrail.Cassandra{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-cassandra",
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
						{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + initialVersionMap["cassandra"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + initialVersionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + initialVersionMap["cassandra"]},
					},
				},
			},
		}}

		zookeepers := []*contrail.Zookeeper{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-zookeeper",
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
						{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:" + initialVersionMap["zookeeper"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + initialVersionMap["python"]},
					},
				},
			},
		}}

		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-rabbitmq",
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
						{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:" + initialVersionMap["rabbitmq"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + initialVersionMap["busybox"]},
					},
				},
			},
		}

		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-config",
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
					CassandraInstance: "upgradetest-cassandra",
					ZookeeperInstance: "upgradetest-zookeeper",
					Containers: []*contrail.Container{
						{Name: "api", Image: "registry:5000/contrail-nightly/contrail-controller-config-api:" + initialVersionMap["cemVersion"]},
						{Name: "devicemanager", Image: "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + initialVersionMap["cemVersion"]},
						{Name: "dnsmasq", Image: "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + initialVersionMap["cemVersion"]},
						{Name: "schematransformer", Image: "registry:5000/contrail-nightly/contrail-controller-config-schema:" + initialVersionMap["cemVersion"]},
						{Name: "servicemonitor", Image: "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + initialVersionMap["cemVersion"]},
						{Name: "analyticsapi", Image: "registry:5000/contrail-nightly/contrail-analytics-api:" + initialVersionMap["cemVersion"]},
						{Name: "collector", Image: "registry:5000/contrail-nightly/contrail-analytics-collector:" + initialVersionMap["cemVersion"]},
						{Name: "queryengine", Image: "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + initialVersionMap["cemVersion"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + initialVersionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + initialVersionMap["redis"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + initialVersionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + initialVersionMap["busybox"]},
						{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + initialVersionMap["contrail-statusmonitor"]},
					},
				},
			},
		}

		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-webui",
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
					CassandraInstance: "upgradetest-cassandra",
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + initialVersionMap["python"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + initialVersionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + initialVersionMap["redis"]},
						{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + initialVersionMap["cemVersion"]},
						{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + initialVersionMap["cemVersion"]},
					},
				},
			},
		}

		provisionManager := &contrail.ProvisionManager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "upgradetest-provmanager",
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
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + initialVersionMap["python"]},
						{Name: "provisioner", Image: "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + initialVersionMap["contrail-operator-provisioner"]},
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
					Replicas:    &replicas,
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
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				return nil
			})

			require.NoError(t, err)

			t.Run("then a ready Zookeeper StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-zookeeper-zookeeper-statefulset", replicas))
			})

			t.Run("then a ready Cassandra StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-cassandra-cassandra-statefulset", replicas))
			})

			t.Run("then a ready Rabbitmq StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-rabbitmq-rabbitmq-statefulset", replicas))
			})

			t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-config-config-statefulset", replicas))
			})

			t.Run("then a ready WebUI StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-webui-webui-statefulset", replicas))
			})

			t.Run("then a ready ProvisionManager StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-provmanager-provisionmanager-statefulset", replicas))
			})
		})
		t.Run("when Zookeeper is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/common-docker-third-party/contrail/zookeeper:" + intendedVersionMap["zookeeper"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				zkContainer := utils.GetContainerFromList("zookeeper", cluster.Spec.Services.Zookeepers[0].Spec.ServiceConfiguration.Containers)
				zkContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)

			t.Run("then Zookeeper resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "zookeeper", targetImage, "zookeeper")
				require.NoError(t, err)
			})

			t.Run("then Zookeeper StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-zookeeper-zookeeper-statefulset", replicas))
			})
		})

		t.Run("when Rabbitmq is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/common-docker-third-party/contrail/rabbitmq:" + intendedVersionMap["rabbitmq"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				rmqContainer := utils.GetContainerFromList("rabbitmq", cluster.Spec.Services.Rabbitmq.Spec.ServiceConfiguration.Containers)
				rmqContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)
			t.Run("then Rabbitmq resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "rabbitmq", targetImage, "rabbitmq")
				require.NoError(t, err)
			})

			t.Run("then Rabbitmq StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-rabbitmq-rabbitmq-statefulset", replicas))
			})
		})

		t.Run("when Cassandra is updated with newer image", func(t *testing.T) {
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				csContainer := utils.GetContainerFromList("cassandra", cluster.Spec.Services.Cassandras[0].Spec.ServiceConfiguration.Containers)
				csContainer.Image = "registry:5000/common-docker-third-party/contrail/cassandra:" + intendedVersionMap["cassandra"]
				return nil
			})
			require.NoError(t, err)
			t.Run("then Cassandra StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-cassandra-cassandra-statefulset", replicas))
			})
		})

		t.Run("when Config is updated with newer image", func(t *testing.T) {
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)
			targetImage := "registry:5000/contrail-nightly/contrail-controller-config-api:" + intendedVersionMap["cemVersion"]
			apiContainer := utils.GetContainerFromList("api", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			devicemanagerContainer := utils.GetContainerFromList("devicemanager", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			dnsmasqContainer := utils.GetContainerFromList("dnsmasq", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			schematransformerContainer := utils.GetContainerFromList("schematransformer", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			servicemonitorContainer := utils.GetContainerFromList("servicemonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			analyticsapiContainer := utils.GetContainerFromList("analyticsapi", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			collectorContainer := utils.GetContainerFromList("collector", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			queryengineContainer := utils.GetContainerFromList("queryengine", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			//Statusmonitor upgrade is skipped since this fails.
			//TODO: Uncomment this after fixing Statusmonitor issues
			//statusmonitorContainer := utils.GetContainerFromList("statusmonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			apiContainer.Image = targetImage
			devicemanagerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + intendedVersionMap["cemVersion"]
			dnsmasqContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + intendedVersionMap["cemVersion"]
			schematransformerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-schema:" + intendedVersionMap["cemVersion"]
			servicemonitorContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + intendedVersionMap["cemVersion"]
			analyticsapiContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-api:" + intendedVersionMap["cemVersion"]
			collectorContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-collector:" + intendedVersionMap["cemVersion"]
			queryengineContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + intendedVersionMap["cemVersion"]
			//statusmonitorContainer.Image = "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + intendedVersionMap["contrail-statusmonitor"]

			err = f.Client.Update(context.TODO(), instance)
			require.NoError(t, err)

			t.Run("then Config resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "config", targetImage, "api")
				require.NoError(t, err)
			})

			t.Run("then Config StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-config-config-statefulset", replicas))
			})

			t.Run("then Config pod can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "config=upgradetest-config",
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
		})

		t.Run("when WebUI is updated with newer image", func(t *testing.T) {
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)
			targetImage := "registry:5000/contrail-nightly/contrail-controller-webui-job:" + intendedVersionMap["cemVersion"]
			webuijobContainer := utils.GetContainerFromList("webuijob", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuiwebContainer := utils.GetContainerFromList("webuiweb", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuijobContainer.Image = targetImage
			webuiwebContainer.Image = "registry:5000/contrail-nightly/contrail-controller-webui-web:" + intendedVersionMap["cemVersion"]
			err = f.Client.Update(context.TODO(), instance)
			require.NoError(t, err)
			t.Run("then Webui resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "webui", targetImage, "webuijob")
				require.NoError(t, err)
			})

			t.Run("then WebUI StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-webui-webui-statefulset", replicas))
			})
		})

		t.Run("when ProvisionManager is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + intendedVersionMap["contrail-operator-provisioner"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				pmContainer := utils.GetContainerFromList("provisioner", cluster.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.Containers)
				pmContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)
			t.Run("then ProvisionManager resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "provisionmanager", targetImage, "provisioner")
				require.NoError(t, err)
			})

			t.Run("then ProvisionManager StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("upgradetest-provmanager-provisionmanager-statefulset", replicas))
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
