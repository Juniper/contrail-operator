package ha

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

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

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail operator is running", func(t *testing.T) {
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

		cluster := getHACluster(namespace)

		t.Run("when manager resource with Config and dependencies are created", func(t *testing.T) {
			//TODO fix 1->3 node scale scenario
			t.Skip()
			var replicas int32 = 1
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})

			require.NoError(t, err)
			assertReplicasReady(t, w, 1)
		})

		t.Run("when replicas is set to 3 in manager", func(t *testing.T) {
			var replicas int32 = 3
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})

			require.NoError(t, err)

			t.Run("then all services are scaled up from 1 to 3 node", func(t *testing.T) {
				assertReplicasReady(t, w, 3)
			})

			t.Run("then all Config pods can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "config=hatest-config",
				})
				assert.NoError(t, err)
				require.NotEmpty(t, configPods.Items)

				for _, pod := range configPods.Items {
					assertConfigIsHealthy(t, proxy, &pod)
				}
			})
		})

		t.Run("when one of the nodes fails", func(t *testing.T) {
			//TODO fix failover scenario
			t.Skip()
			nodes, err := f.KubeClient.CoreV1().Nodes().List(meta.ListOptions{
				LabelSelector: "node-role.kubernetes.io/master=",
			})
			assert.NoError(t, err)
			require.NotEmpty(t, nodes.Items)
			node := nodes.Items[0]
			node.Spec.Taints = append(node.Spec.Taints, core.Taint{
				Key:    "e2e.test/failure",
				Effect: core.TaintEffectNoExecute,
			})

			_, err = f.KubeClient.CoreV1().Nodes().Update(&node)
			assert.NoError(t, err)
			t.Run("then all services should have 2 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       time.Minute * 5,
					RetryInterval: time.Second * 15,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertReplicasReady(t, w, 2)
			})

			t.Run("then ready Config pods can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "config=hatest-config",
				})
				assert.NoError(t, err)
				require.NotEmpty(t, configPods.Items)

				healthyConfigs := 0
				for _, pod := range configPods.Items {
					for _, c := range pod.Status.Conditions {
						if c.Type == core.PodReady && c.Status == core.ConditionTrue {
							assertConfigIsHealthy(t, proxy, &pod)
							healthyConfigs++
						}
					}
				}
				assert.Equalf(t, 2, healthyConfigs, "expected 2 healthy configs, got %d", healthyConfigs)
			})
		})

		t.Run("when all nodes are back operational", func(t *testing.T) {
			err := untaintNodes(f.KubeClient, "e2e.test/failure")
			assert.NoError(t, err)
			t.Run("then all services should have 3 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       time.Minute * 5,
					RetryInterval: retryInterval,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertReplicasReady(t, w, 3)
			})

			t.Run("then all Config pods can process requests", func(t *testing.T) {
				configPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "config=hatest-config",
				})
				assert.NoError(t, err)
				require.NotEmpty(t, configPods.Items)

				for _, pod := range configPods.Items {
					assertConfigIsHealthy(t, proxy, &pod)
				}
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

func untaintNodes(k kubernetes.Interface, taintKey string) error {
	nodes, err := k.CoreV1().Nodes().List(meta.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master=",
	})

	if err != nil {
		return err
	}

	for _, n := range nodes.Items {
		for i, tn := range n.Spec.Taints {
			if tn.Key != "e2e.test/failure" {
				continue
			}
			s := n.Spec.Taints
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			n.Spec.Taints = s[:len(s)-1]
			_, err = k.CoreV1().Nodes().Update(&n)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func assertReplicasReady(t *testing.T, w wait.Wait, r int32) {
	t.Run(fmt.Sprintf("then a Zookeeper StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a Cassandra StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a Rabbit StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-rabbitmq-rabbitmq-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a Config StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a WebUI StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a ProvisionManager StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", r))
	})
}

func assertConfigIsHealthy(t *testing.T, proxy *kubeproxy.HTTPProxy, p *core.Pod) {
	configProxy := proxy.NewSecureClient("contrail", p.Name, 8082)
	req, err := configProxy.NewRequest(http.MethodGet, "/projects", nil)
	assert.NoError(t, err)
	res, err := configProxy.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func getHACluster(namespace string) *contrail.Manager {
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

	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cluster1",
			Namespace: namespace,
		},
		Spec: contrail.ManagerSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				HostNetwork: &trueVal,
				Tolerations: []core.Toleration{
					{
						Effect:   core.TaintEffectNoSchedule,
						Operator: core.TolerationOpExists,
					},
				},
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
}
