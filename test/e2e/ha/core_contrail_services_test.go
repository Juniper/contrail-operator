package ha

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8swait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
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
	"cemVersion":                    "2008.10",
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
	"cemVersion":                    cemRelease,
	"python":                        "3.8.2-alpine",
	"redis":                         "4.0.2",
	"busybox":                       "1.31",
	"rabbitmq":                      "3.7.16",
	"contrail-statusmonitor":        "R2008.latest",
	"contrail-operator-provisioner": "R2008.latest",
}

func TestHACoreContrailServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	require.NoError(t, err)

	log := logger.New(t, namespace, test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}

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

		nodeLabelKey := "test-ha"
		cluster := getHACluster(namespace, nodeLabelKey)

		t.Run("when manager resource with Config and dependencies are created", func(t *testing.T) {
			_, err = controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				return nil
			})
			require.NoError(t, err)
		})

		t.Run("when manager resource with Config and dependencies are created", func(t *testing.T) {
			err := labelOneNode(f.KubeClient, nodeLabelKey)
			require.NoError(t, err)
			t.Run("then all services are started with replica 1", func(t *testing.T) {
				assertReplicasReady(t, w, 1)
			})
		})

		t.Run("when nodes are replicated to 3", func(t *testing.T) {
			err := labelAllNodes(f.KubeClient, nodeLabelKey)
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

		t.Run("when manager resource is upgraded", func(t *testing.T) {
			t.Skip()
			var replicas int32 = 3
			instance := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, instance)
			assert.NoError(t, err)

			updateManagerImages(t, f, instance)

			t.Run("then all Pods have updated image", func(t *testing.T) {
				requirePodsHaveUpdatedImages(t, f, namespace, log)
			})

			t.Run("then all services should have 3 ready replicas", func(t *testing.T) {
				assertReplicasReady(t, w, replicas)
			})

			t.Run("then Config pod can process requests", func(t *testing.T) {
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
			nodes, err := f.KubeClient.CoreV1().Nodes().List(meta.ListOptions{
				LabelSelector: labelKeyToSelector(nodeLabelKey),
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
			err := untaintNodes(f.KubeClient, labelKeyToSelector(nodeLabelKey))
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
			t.Run("then test label is removed from nodes", func(t *testing.T) {
				err := removeLabel(f.KubeClient, nodeLabelKey)
				require.NoError(t, err)
			})
		})
	})
}

func labelOneNode(k kubernetes.Interface, labelKey string) error {
	nodes, err := k.CoreV1().Nodes().List(meta.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master=",
	})
	if err != nil {
		return err
	}
	if len(nodes.Items) == 0 {
		return fmt.Errorf("no master nodes found")
	}
	node := nodes.Items[0]
	node.Labels[labelKey] = ""
	_, err = k.CoreV1().Nodes().Update(&node)

	return err
}

func labelAllNodes(k kubernetes.Interface, labelKey string) error {
	nodes, err := k.CoreV1().Nodes().List(meta.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master=",
	})
	if err != nil {
		return err
	}
	for _, n := range nodes.Items {
		n.Labels[labelKey] = ""
		_, err = k.CoreV1().Nodes().Update(&n)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeLabel(k kubernetes.Interface, labelKey string) error {
	nodes, err := k.CoreV1().Nodes().List(meta.ListOptions{
		LabelSelector: labelKeyToSelector(labelKey),
	})

	if err != nil {
		return err
	}

	for _, n := range nodes.Items {
		delete(n.Labels, labelKey)
		_, err = k.CoreV1().Nodes().Update(&n)
		if err != nil {
			return err
		}
	}
	return nil
}

func untaintNodes(k kubernetes.Interface, nodeLabelSelector string) error {
	nodes, err := k.CoreV1().Nodes().List(meta.ListOptions{
		LabelSelector: nodeLabelSelector,
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

	t.Run(fmt.Sprintf("then a Control StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-control-control-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a Config StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset", r))
	})

	t.Run(fmt.Sprintf("then a WebUI StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset", r))
	})

	// ProvisionManager is not scalable and is deployed in one replica
	t.Run("then a ProvisionManager StatefulSet has 1 ready replica", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("hatest-provmanager-provisionmanager-statefulset", 1))
	})
}

func assertConfigIsHealthy(t *testing.T, proxy *kubeproxy.HTTPProxy, p *core.Pod) {
	configProxy := proxy.NewSecureClient("contrail", p.Name, 8082)
	req, err := configProxy.NewRequest(http.MethodGet, "/projects", nil)
	assert.NoError(t, err)
	var res *http.Response
	err = k8swait.Poll(retryInterval, time.Second*20, func() (done bool, err error) {
		res, err = configProxy.Do(req)
		if err == nil {
			return true, nil
		}
		t.Log(err)
		return false, nil
	})
	require.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode, string(body))
}

func getHACluster(namespace, nodeLabel string) *contrail.Manager {
	trueVal := true

	cassandras := []*contrail.Cassandra{{
		ObjectMeta: meta.ObjectMeta{
			Name:      "hatest-cassandra",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: contrail.CassandraSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
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
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
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
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
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

	controls := []*contrail.Control{{
		ObjectMeta: meta.ObjectMeta{
			Name:      "hatest-control",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1", "control_role": "master"},
		},
		Spec: contrail.ControlSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
			},
			ServiceConfiguration: contrail.ControlConfiguration{
				CassandraInstance: "hatest-cassandra",
				Containers: []*contrail.Container{
					{Name: "control", Image: "registry:5000/contrail-nightly/contrail-controller-control-control:" + versionMap["cemVersion"]},
					{Name: "dns", Image: "registry:5000/contrail-nightly/contrail-controller-control-dns:" + versionMap["cemVersion"]},
					{Name: "named", Image: "registry:5000/contrail-nightly/contrail-controller-control-named:" + versionMap["cemVersion"]},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
					{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
				},
			},
		},
	}}

	webui := &contrail.Webui{
		ObjectMeta: meta.ObjectMeta{
			Name:      "hatest-webui",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: contrail.WebuiSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
			},
			ServiceConfiguration: contrail.WebuiConfiguration{
				CassandraInstance: "hatest-cassandra",
				Containers: []*contrail.Container{
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
					{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + versionMap["redis"]},
					{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + versionMap["cemVersion"]},
					{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + versionMap["cemVersion"]},
				},
			},
		},
	}

	oneVal := int32(1)
	provisionManager := &contrail.ProvisionManager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "hatest-provmanager",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: contrail.ProvisionManagerSpec{
			CommonConfiguration: contrail.PodConfiguration{
				Replicas: &oneVal,
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
			CommonConfiguration: contrail.ManagerConfiguration{
				HostNetwork: &trueVal,
				Tolerations: []core.Toleration{
					{
						Effect:   core.TaintEffectNoSchedule,
						Operator: core.TolerationOpExists,
					},
				},
				NodeSelector: map[string]string{nodeLabel: ""},
			},
			Services: contrail.Services{
				Cassandras:       cassandras,
				Zookeepers:       zookeepers,
				Controls:         controls,
				Config:           config,
				Webui:            webui,
				Rabbitmq:         rabbitmq,
				ProvisionManager: provisionManager,
			},
		},
	}
}

func labelKeyToSelector(key string) string {
	return key + "="
}

func requirePodsHaveUpdatedImages(t *testing.T, f *test.Framework, namespace string, log logger.Logger) {
	t.Run("then Zookeeper has updated image", func(t *testing.T) {
		t.Parallel()
		zkContainerImage := "registry:5000/common-docker-third-party/contrail/zookeeper:" + targetVersionMap["zookeeper"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=zookeeper", zkContainerImage, "zookeeper")
		require.NoError(t, err)
	})

	t.Run("then Rabbitmq has updated image", func(t *testing.T) {
		t.Parallel()
		rmqContainerImage := "registry:5000/common-docker-third-party/contrail/rabbitmq:" + targetVersionMap["rabbitmq"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=rabbitmq", rmqContainerImage, "rabbitmq")
		require.NoError(t, err)
	})

	t.Run("then Cassandra has updated image", func(t *testing.T) {
		t.Parallel()
		csContainerImage := "registry:5000/common-docker-third-party/contrail/cassandra:" + targetVersionMap["cassandra"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=cassandra", csContainerImage, "cassandra")
		require.NoError(t, err)
	})

	t.Run("then Control has updated image", func(t *testing.T) {
		t.Parallel()
		controlContainerImage := "registry:5000/contrail-nightly/contrail-controller-control-control:" + targetVersionMap["cemVersion"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=control", controlContainerImage, "control")
		require.NoError(t, err)
	})

	t.Run("then Config has updated image", func(t *testing.T) {
		t.Parallel()
		apiContainerImage := "registry:5000/contrail-nightly/contrail-controller-config-api:" + targetVersionMap["cemVersion"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=config", apiContainerImage, "api")
		require.NoError(t, err)
	})

	t.Run("then Webui has updated image", func(t *testing.T) {
		t.Parallel()
		webuijobContainerImage := "registry:5000/contrail-nightly/contrail-controller-webui-job:" + targetVersionMap["cemVersion"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=webui", webuijobContainerImage, "webuijob")
		require.NoError(t, err)
	})

	t.Run("then ProvisionManager has updated image", func(t *testing.T) {
		t.Parallel()
		pmContainerImage := "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + targetVersionMap["contrail-operator-provisioner"]
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "contrail_manager=provisionmanager", pmContainerImage, "provisioner")
		require.NoError(t, err)
	})
}

func updateManagerImages(t *testing.T, f *test.Framework, instance *contrail.Manager) {
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

	controlNodeContainer := utils.GetContainerFromList("control", instance.Spec.Services.Controls[0].Spec.ServiceConfiguration.Containers)
	controlDNSContainer := utils.GetContainerFromList("dns", instance.Spec.Services.Controls[0].Spec.ServiceConfiguration.Containers)
	controlNamedContainer := utils.GetContainerFromList("named", instance.Spec.Services.Controls[0].Spec.ServiceConfiguration.Containers)
	controlNodeContainer.Image = "registry:5000/contrail-nightly/contrail-controller-control-control:" + targetVersionMap["cemVersion"]
	controlDNSContainer.Image = "registry:5000/contrail-nightly/contrail-controller-control-dns:" + targetVersionMap["cemVersion"]
	controlNamedContainer.Image = "registry:5000/contrail-nightly/contrail-controller-control-named:" + targetVersionMap["cemVersion"]

	pmContainer := utils.GetContainerFromList("provisioner", instance.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.Containers)
	pmContainer.Image = "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + targetVersionMap["contrail-operator-provisioner"]

	err := f.Client.Update(context.TODO(), instance)
	require.NoError(t, err)
}
