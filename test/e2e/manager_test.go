// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	goctx "context"

	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apis "github.com/Juniper/contrail-operator/pkg/apis"

	"github.com/operator-framework/operator-sdk/pkg/test"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrailwait "github.com/Juniper/contrail-operator/test/wait"
)

/*
func TestRabbitmq(t *testing.T) {
	rabbitmqList := &v1alpha1.RabbitmqList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Rabbitmq",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
	}
	err := test.AddToFrameworkScheme(apis.AddToScheme, rabbitmqList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("rabbitmq-group", func(t *testing.T) {
		t.Run("Cluster", RabbitmqCluster)
	})
}
*/

func TestManager(t *testing.T) {
	managerList := &v1alpha1.ManagerList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Manager",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
	}
	err := test.AddToFrameworkScheme(apis.AddToScheme, managerList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("manager-group", func(t *testing.T) {
		t.Run("Cluster", ManagerCluster)
	})
}

var initialVersionMap = map[string]string{
	"rabbitmq":    "3.7.16",
	"cassandra":   "3.11.3",
	"zookeeper":   "3.5.4-beta",
	"config":      "1912-latest",
	"control":     "1912-latest",
	"kubemanager": "1912-latest",
}

var targetVersionMap = map[string]string{
	"rabbitmq":    "3.7.17",
	"cassandra":   "3.11.4",
	"zookeeper":   "3.5.5",
	"config":      "2002-latest",
	"control":     "2002-latest",
	"kubemanager": "2002-latest",
}

func ManagerCluster(t *testing.T) {
	t.Parallel()
	f := test.Global
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()

	if err := test.AddToFrameworkScheme(v1alpha1.SchemeBuilder.AddToScheme, &v1alpha1.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}

	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	var replicas int32 = 1
	var hostNetwork = false
	manager := getManager(namespace, replicas, hostNetwork, initialVersionMap)

	err = f.Client.Create(goctx.TODO(), &manager, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}
	err = waitForZookeeper(t, f, ctx, namespace, "zookeeper1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForCassandra(t, f, ctx, namespace, "cassandra1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForRabbitmq(t, f, ctx, namespace, "rabbitmq1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForConfig(t, f, ctx, namespace, "config1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForControl(t, f, ctx, namespace, "control1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForKubemanager(t, f, ctx, namespace, "kubemanager1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForManager(t, f, ctx, namespace, "cluster1", retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = upgradeZookeeper(t, f, ctx, namespace, "cluster1"); err != nil {
		t.Fatal(err)
	}

	err = zookeeperVersion(t, f, ctx, namespace, "zookeeper1")
	if err != nil {
		t.Fatal(err)
	}

	pp := meta.DeletePropagationForeground
	err = f.Client.Delete(context.TODO(), &manager, &client.DeleteOptions{
		PropagationPolicy: &pp,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = contrailwait.Contrail{
		Namespace:     namespace,
		Timeout:       5 * time.Minute,
		RetryInterval: retryInterval,
		Client:        f.Client,
	}.ForManagerDeletion(manager.Name)
	if err != nil {
		t.Fatal(err)
	}
}

func getManager(namespace string, replicas int32, hostNetwork bool, versionMap map[string]string) v1alpha1.Manager {
	create := true
	return v1alpha1.Manager{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Manager",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: namespace,
		},
		Spec: v1alpha1.ManagerSpec{
			CommonConfiguration: v1alpha1.CommonConfiguration{
				Replicas:         &replicas,
				HostNetwork:      &hostNetwork,
				ImagePullSecrets: []string{"contrail-nightly"},
			},
			Services: v1alpha1.Services{
				Rabbitmq: &v1alpha1.Rabbitmq{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rabbitmq1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.RabbitmqSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.RabbitmqConfiguration{
							Containers: map[string]*v1alpha1.Container{
								"rabbitmq": &v1alpha1.Container{Image: "registry:5000/rabbitmq:" + versionMap["rabbitmq"]},
								"init":     &v1alpha1.Container{Image: "registry:5000/busybox"},
							},
						},
					},
				},
				Zookeepers: []*v1alpha1.Zookeeper{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "zookeeper1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.ZookeeperSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.ZookeeperConfiguration{
							Containers: map[string]*v1alpha1.Container{
								"zookeeper": &v1alpha1.Container{Image: "registry:5000/zookeeper:" + versionMap["zookeeper"]},
								"init":      &v1alpha1.Container{Image: "registry:5000/busybox"},
							},
						},
					},
				}},
				Cassandras: []*v1alpha1.Cassandra{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cassandra1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.CassandraSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.CassandraConfiguration{
							Containers: map[string]*v1alpha1.Container{
								"cassandra": &v1alpha1.Container{Image: "registry:5000/cassandra:" + versionMap["cassandra"]},
								"init":      &v1alpha1.Container{Image: "registry:5000/busybox"},
								"init2":     &v1alpha1.Container{Image: "registry:5000/cassandra:" + versionMap["cassandra"]},
							},
						},
					},
				}},
				Config: &v1alpha1.Config{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "config1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.ConfigSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.ConfigConfiguration{
							CassandraInstance: "cassandra1",
							ZookeeperInstance: "zookeeper1",

							Containers: map[string]*v1alpha1.Container{
								"api":               &v1alpha1.Container{Image: "registry:5000/contrail-controller-config-api:" + versionMap["config"]},
								"devicemanager":     &v1alpha1.Container{Image: "registry:5000/contrail-controller-config-devicemgr:dev-5"}, // Using custom dev version until required changes are merged upstream
								"dnsmasq":           &v1alpha1.Container{Image: "registry:5000/contrail-controller-config-dnsmasq:dev"},     // Using custom dev version until required changes are merged upstream
								"schematransformer": &v1alpha1.Container{Image: "registry:5000/contrail-controller-config-schema:" + versionMap["config"]},
								"servicemonitor":    &v1alpha1.Container{Image: "registry:5000/contrail-controller-config-svcmonitor:" + versionMap["config"]},
								"analyticsapi":      &v1alpha1.Container{Image: "registry:5000/contrail-analytics-api:" + versionMap["config"]},
								"collector":         &v1alpha1.Container{Image: "registry:5000/contrail-analytics-collector:" + versionMap["config"]},
								"queryengine":       &v1alpha1.Container{Image: "registry:5000/contrail-analytics-query-engine:" + versionMap["config"]},
								"redis":             &v1alpha1.Container{Image: "registry:5000/redis:4.0.2"},
								"nodeinit":          &v1alpha1.Container{Image: "registry:5000/contrail-node-init:" + versionMap["config"]},
								"init":              &v1alpha1.Container{Image: "registry:5000/python:alpine"},
								"init2":             &v1alpha1.Container{Image: "registry:5000/busybox"},
							},
						},
					},
				},
				Controls: []*v1alpha1.Control{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "control1",
						Namespace: namespace,
						Labels: map[string]string{
							"contrail_cluster": "cluster1",
							"control_role":     "master",
						},
					},
					Spec: v1alpha1.ControlSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.ControlConfiguration{
							CassandraInstance: "cassandra1",
							ZookeeperInstance: "zookeeper1",
							Containers: map[string]*v1alpha1.Container{
								"control":       &v1alpha1.Container{Image: "registry:5000/contrail-controller-control-control:" + versionMap["control"]},
								"dns":           &v1alpha1.Container{Image: "registry:5000/contrail-controller-control-dns:" + versionMap["control"]},
								"named":         &v1alpha1.Container{Image: "registry:5000/contrail-controller-control-named:" + versionMap["control"]},
								"statusmonitor": &v1alpha1.Container{Image: "registry:5000/contrail-statusmonitor:debug"},
								"nodeinit":      &v1alpha1.Container{Image: "registry:5000/contrail-node-init:" + versionMap["control"]},
								"init":          &v1alpha1.Container{Image: "registry:5000/python:alpine"},
							},
						},
					},
				}},
				Kubemanagers: []*v1alpha1.Kubemanager{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubemanager1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.KubemanagerSpec{
						CommonConfiguration: v1alpha1.CommonConfiguration{
							Create:       &create,
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: v1alpha1.KubemanagerConfiguration{
							CassandraInstance: "cassandra1",
							ZookeeperInstance: "zookeeper1",

							Containers: map[string]*v1alpha1.Container{
								"kubemanager": &v1alpha1.Container{Image: "registry:5000/contrail-kubernetes-kube-manager:" + versionMap["kubemanager"]},
								"nodeinit":    &v1alpha1.Container{Image: "registry:5000/contrail-node-init:" + versionMap["kubemanager"]},
								"init":        &v1alpha1.Container{Image: "registry:5000/busybox"},
							},
						},
					},
				}},
			},
		},
	}
}

func RabbitmqCluster(t *testing.T) {
	t.Parallel()
	f := test.Global
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	var replicas int32 = 1
	rabbitmq := &v1alpha1.Rabbitmq{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rabbitmq1",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: v1alpha1.RabbitmqSpec{
			CommonConfiguration: v1alpha1.CommonConfiguration{
				Replicas:     &replicas,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: v1alpha1.RabbitmqConfiguration{
				Containers: map[string]*v1alpha1.Container{
					"rabbitmq": &v1alpha1.Container{Image: "registry:5000/rabbitmq:3.7"},
					"init":     &v1alpha1.Container{Image: "registry:5000/busybox"},
				},
			},
		},
	}

	err = f.Client.Create(goctx.TODO(), rabbitmq, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "rabbitmq1-rabbitmq-deployment", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}
}

func upgradeZookeeper(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string) error {
	instance := &v1alpha1.Manager{}
	err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return err
	}
	instance.Spec.Services.Zookeepers[0].Spec.ServiceConfiguration.Containers["zookeeper"].Image = "registry:5000/zookeeper:" + targetVersionMap["zookeeper"]
	err = f.Client.Update(goctx.TODO(), instance)
	if err != nil {
		return err
	}
	err = waitForZookeeper(t, f, ctx, namespace, "zookeeper1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForConfig(t, f, ctx, namespace, "config1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForControl(t, f, ctx, namespace, "control1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForKubemanager(t, f, ctx, namespace, "kubemanager1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	return nil
}
func zookeeperVersion(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string) error {
	instance := &v1alpha1.Zookeeper{}
	err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return err
	}
	t.Log("Got zookeeper")
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": "zookeeper",
		"zookeeper": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	podList := &corev1.PodList{}
	err = f.Client.List(goctx.TODO(), podList, listOps)
	if err != nil {
		return err
	}
	t.Log("Got zookeeper pods")
	for _, pod := range podList.Items {
		command := []string{"bash", "-c", "echo stats|nc " + pod.Status.PodIP + " " + instance.Status.Ports.ClientPort + "|grep 'Zookeeper version:' |sed -e 's/.*Zookeeper version: *\\(.*\\)-[^-]*,.*/\\1/'"}
		//command = "echo bla > /bla"
		output, stderr, err := v1alpha1.ExecToPodThroughAPI(command, "zookeeper", pod.Name, namespace, nil)
		if len(stderr) != 0 {
			fmt.Println("STDERR:", stderr)
		}
		if err != nil {
			fmt.Printf("Error occurred while `exec`ing to the Pod %q, namespace %q, command %q. Error: %+v\n", pod.Name, namespace, command, err)
		} else {
			fmt.Println("Output:")
			fmt.Println(output)
			if output != targetVersionMap["zookeeper"] {
				return err
			}
		}
	}
	return nil
}

func upgradeCassandra(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string) error {
	instance := &v1alpha1.Manager{}
	err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return err
	}

	instance.Spec.Services.Cassandras[0].Spec.ServiceConfiguration.Containers["cassandra"].Image = "registry:5000/cassandra:" + targetVersionMap["cassandra"]
	err = f.Client.Update(goctx.TODO(), instance)
	if err != nil {
		return err
	}
	err = waitForCassandra(t, f, ctx, namespace, "cassandra1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForConfig(t, f, ctx, namespace, "config1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}
	return nil
}

func imageUpgradeTest(t *testing.T, f *test.Framework, ctx *test.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	manager := v1alpha1.Manager{}
	f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, &manager)
	if err != nil {
		return fmt.Errorf("could not get manager: %v", err)
	}
	manager = getManager(namespace, 1, false, targetVersionMap)
	err = f.Client.Create(goctx.TODO(), &manager, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}
	err = waitForZookeeper(t, f, ctx, namespace, "zookeeper1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForCassandra(t, f, ctx, namespace, "cassandra1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForRabbitmq(t, f, ctx, namespace, "rabbitmq1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForConfig(t, f, ctx, namespace, "config1", 1, retryInterval, waitTimeout)
	if err != nil {
		t.Fatal(err)
	}
	return nil
}

func managerScaleTest(t *testing.T, f *test.Framework, ctx *test.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	manager := &v1alpha1.Manager{}
	f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, manager)
	if err != nil {
		return fmt.Errorf("could not get manager: %v", err)
	}
	var replicas int32 = 3
	manager.Spec.CommonConfiguration.Replicas = &replicas
	err = f.Client.Update(goctx.TODO(), manager)
	if err != nil {
		return fmt.Errorf("could not update manager: %v", err)
	}
	waitTimeout = time.Second * 120
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "rabbitmq1-rabbitmq-deployment", 3, retryInterval, waitTimeout)
	if err != nil {
		return fmt.Errorf("rabbitmq deployment is wrong: %v", err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "zookeeper1-zookeeper-deployment", 3, retryInterval, waitTimeout)
	if err != nil {
		return fmt.Errorf("zookeeper deployment is wrong: %v", err)
	}
	return nil
}

func waitForZookeeper(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Zookeeper{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s zookeeper\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s zookeeper\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Zoopkeeper %s available\n", name)
	return nil
}

func waitForCassandra(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Cassandra{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s cassandra\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s cassandra\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Cassandra %s available\n", name)
	return nil
}

func waitForRabbitmq(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Rabbitmq{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s rabbitmq\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s rabbitmq\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Rabbitmq %s available\n", name)
	return nil
}

func waitForConfig(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Config{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s config\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s config\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Config %s available\n", name)
	return nil
}

func waitForKubemanager(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Kubemanager{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s kubemanager\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s kubemanager\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Kubemanager %s available\n", name)
	return nil
}

func waitForControl(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, replicas int, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Control{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s control\n", name)
				return false, nil
			}
			return false, err
		}

		if instance.Status.Active != nil && *instance.Status.Active {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s control\n", name)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Control %s available\n", name)
	return nil
}

func waitForManager(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string, retryInterval, waitTimeout time.Duration) error {
	err := wait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
		instance := &v1alpha1.Manager{}
		err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s Manager\n", name)
				return false, nil
			}
			return false, err
		}
		if !instance.IsClusterReady() {
			t.Logf("Waiting for full availability of %s Manager\n", name)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	t.Logf("Manager %s available\n", name)
	return nil
}
