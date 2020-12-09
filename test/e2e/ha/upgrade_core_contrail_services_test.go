package ha

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
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

func TestUpgradeCoreContrailServices(t *testing.T) {
	if testing.Short() {
		t.Skip("it is a long test")
	}
	ctx := test.NewContext(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetOperatorNamespace()
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

		w := wait.Wait{
			Namespace:     namespace,
			Timeout:       waitTimeout,
			RetryInterval: retryInterval,
			KubeClient:    f.KubeClient,
			Logger:        log,
		}

		nodeLabelKey := "test-ha-upgrade"
		storagePath := "/mnt/storage/" + uuid.New().String()
		cluster := getHACluster(namespace, nodeLabelKey, storagePath)
		var replicas int32 = 3

		t.Run("when manager resource with Config and dependencies are created", func(t *testing.T) {
			err := labelAllNodes(f.KubeClient, nodeLabelKey)
			require.NoError(t, err)
			_, err = controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				return nil
			})

			require.NoError(t, err)
			assertReplicasReady(t, w, replicas)
		})

		t.Run("when Zookeeper is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/common-docker-third-party/contrail/zookeeper:" + targetVersionMap["zookeeper"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				zkContainer := utils.GetContainerFromList("zookeeper", cluster.Spec.Services.Zookeepers[0].Spec.ServiceConfiguration.Containers)
				zkContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)

			t.Run("then Zookeeper resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "zookeeper", targetImage, "zookeeper")
				require.NoError(t, err)
			})

			t.Run("then Zookeeper StatefulSet should be ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset", replicas))
			})
		})

		t.Run("when Rabbitmq is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/common-docker-third-party/contrail/rabbitmq:" + targetVersionMap["rabbitmq"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				rmqContainer := utils.GetContainerFromList("rabbitmq", cluster.Spec.Services.Rabbitmq.Spec.ServiceConfiguration.Containers)
				rmqContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)
			t.Run("then Rabbitmq resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "rabbitmq", targetImage, "rabbitmq")
				require.NoError(t, err)
			})

			t.Run("then Rabbitmq StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-rabbitmq-rabbitmq-statefulset", replicas))
			})
		})

		t.Run("when Cassandra is updated with newer image", func(t *testing.T) {
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				csContainer := utils.GetContainerFromList("cassandra", cluster.Spec.Services.Cassandras[0].Spec.ServiceConfiguration.Containers)
				csContainer.Image = "registry:5000/common-docker-third-party/contrail/cassandra:" + targetVersionMap["cassandra"]
				return nil
			})
			require.NoError(t, err)
			t.Run("then Cassandra StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset", replicas))
			})
		})

		t.Run("when Config is updated with newer image", func(t *testing.T) {
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)
			targetImage := "registry:5000/contrail-nightly/contrail-controller-config-api:" + targetVersionMap["cemVersion"]
			apiContainer := utils.GetContainerFromList("api", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			devicemanagerContainer := utils.GetContainerFromList("devicemanager", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			dnsmasqContainer := utils.GetContainerFromList("dnsmasq", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			schematransformerContainer := utils.GetContainerFromList("schematransformer", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			servicemonitorContainer := utils.GetContainerFromList("servicemonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			analyticsapiContainer := utils.GetContainerFromList("analyticsapi", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			collectorContainer := utils.GetContainerFromList("collector", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			queryengineContainer := utils.GetContainerFromList("queryengine", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			statusmonitorContainer := utils.GetContainerFromList("statusmonitor", instance.Spec.Services.Config.Spec.ServiceConfiguration.Containers)
			apiContainer.Image = targetImage
			devicemanagerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + targetVersionMap["cemVersion"]
			dnsmasqContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + targetVersionMap["cemVersion"]
			schematransformerContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-schema:" + targetVersionMap["cemVersion"]
			servicemonitorContainer.Image = "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + targetVersionMap["cemVersion"]
			analyticsapiContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-api:" + targetVersionMap["cemVersion"]
			collectorContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-collector:" + targetVersionMap["cemVersion"]
			queryengineContainer.Image = "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + targetVersionMap["cemVersion"]
			statusmonitorContainer.Image = "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + targetVersionMap["contrail-statusmonitor"]

			err = f.Client.Update(context.TODO(), instance)
			require.NoError(t, err)

			t.Run("then Config resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "config", targetImage, "api")
				require.NoError(t, err)
			})

			t.Run("then Config StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", replicas))
			})

			t.Run("then Config pod can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(context.Background(), meta.ListOptions{
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
		})

		t.Run("when WebUI is updated with newer image", func(t *testing.T) {
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)
			targetImage := "registry:5000/contrail-nightly/contrail-controller-webui-job:" + targetVersionMap["cemVersion"]
			webuijobContainer := utils.GetContainerFromList("webuijob", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuiwebContainer := utils.GetContainerFromList("webuiweb", instance.Spec.Services.Webui.Spec.ServiceConfiguration.Containers)
			webuijobContainer.Image = targetImage
			webuiwebContainer.Image = "registry:5000/contrail-nightly/contrail-controller-webui-web:" + targetVersionMap["cemVersion"]
			err = f.Client.Update(context.TODO(), instance)
			require.NoError(t, err)
			t.Run("then Webui resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "webui", targetImage, "webuijob")
				require.NoError(t, err)
			})

			t.Run("then WebUI StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", replicas))
			})
		})

		t.Run("when ProvisionManager is updated with newer image", func(t *testing.T) {
			targetImage := "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + targetVersionMap["contrail-operator-provisioner"]
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				pmContainer := utils.GetContainerFromList("provisioner", cluster.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.Containers)
				pmContainer.Image = targetImage
				return nil
			})
			require.NoError(t, err)
			t.Run("then ProvisionManager resource has updated image", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPodImageChange(f.KubeClient, "provisionmanager", targetImage, "provisioner")
				require.NoError(t, err)
			})

			t.Run("then ProvisionManager StatefulSet should be updated and ready", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", 1))
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
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})

			t.Run("then persistent volumes are removed", func(t *testing.T) {
				err := deleteAllPVs(f.KubeClient, "local-storage")
				require.NoError(t, err)
			})

			t.Run("then test label is removed from nodes", func(t *testing.T) {
				err := removeLabel(f.KubeClient, nodeLabelKey)
				require.NoError(t, err)
			})
		})
	})
}
