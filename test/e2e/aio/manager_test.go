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

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/operator-framework/operator-sdk/pkg/test"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/apis"
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/test/logger"
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

var targetVersionMap = map[string]string{
	"rabbitmq":    "3.7.17",
	"cassandra":   "3.11.4",
	"zookeeper":   "3.5.6",
	"config":      "2002-latest",
	"control":     "2002-latest",
	"kubemanager": "2002-latest",
}

func ManagerCluster(t *testing.T) {

	initialVersionMap := map[string]string{
		"rabbitmq":                      "3.7.16",
		"cassandra":                     "3.11.3",
		"zookeeper":                     "3.5.5",
		"config":                        cemRelease,
		"control":                       cemRelease,
		"kubemanager":                   cemRelease,
		"contrail-operator-provisioner": buildTag,
		"contrail-statusmonitor":        buildTag,
	}

	f := test.Global
	ctx := test.NewContext(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(v1alpha1.SchemeBuilder.AddToScheme, &v1alpha1.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := test.AddToFrameworkScheme(corev1.AddToScheme, &corev1.PersistentVolumeList{}); err != nil {
		t.Fatalf("Failed to add core framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}

	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetOperatorNamespace()
	if err != nil {
		t.Fatal(err)
	}

	var replicas int32 = 1
	var hostNetwork = true
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
	err = f.Client.DeleteAllOf(context.TODO(), &corev1.PersistentVolume{})
	if err != nil {
		t.Fatal(err)
	}
	err = contrailwait.Contrail{
		Namespace:     namespace,
		Timeout:       5 * time.Minute,
		RetryInterval: retryInterval,
		Client:        f.Client,
		Logger:        log,
	}.ForManagerDeletion(manager.Name)
	if err != nil {
		t.Fatal(err)
	}
}

func getManager(namespace string, replicas int32, hostNetwork bool, versionMap map[string]string) v1alpha1.Manager {
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
			CommonConfiguration: v1alpha1.ManagerConfiguration{
				HostNetwork:      &hostNetwork,
				ImagePullSecrets: []string{"contrail-nightly"},
				NodeSelector:     map[string]string{"node-role.juniper.net/contrail": ""},
			},
			Services: v1alpha1.Services{
				Rabbitmq: &v1alpha1.Rabbitmq{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rabbitmq1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.RabbitmqSpec{
						ServiceConfiguration: v1alpha1.RabbitmqConfiguration{
							Containers: []*v1alpha1.Container{
								{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:" + versionMap["rabbitmq"]},
								{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
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
						ServiceConfiguration: v1alpha1.ZookeeperConfiguration{
							Containers: []*v1alpha1.Container{
								{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:" + versionMap["zookeeper"]},
								{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
								{Name: "conf-init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
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
						ServiceConfiguration: v1alpha1.CassandraConfiguration{
							Containers: []*v1alpha1.Container{
								{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
								{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
								{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
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
						ServiceConfiguration: v1alpha1.ConfigConfiguration{
							CassandraInstance: "cassandra1",
							ZookeeperInstance: "zookeeper1",

							Containers: []*v1alpha1.Container{
								{Name: "api", Image: "registry:5000/contrail-nightly/contrail-controller-config-api:" + versionMap["config"]},
								{Name: "devicemanager", Image: "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + versionMap["config"]},
								{Name: "dnsmasq", Image: "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + versionMap["config"]},
								{Name: "schematransformer", Image: "registry:5000/contrail-nightly/contrail-controller-config-schema:" + versionMap["config"]},
								{Name: "servicemonitor", Image: "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + versionMap["config"]},
								{Name: "analyticsapi", Image: "registry:5000/contrail-nightly/contrail-analytics-api:" + versionMap["config"]},
								{Name: "collector", Image: "registry:5000/contrail-nightly/contrail-analytics-collector:" + versionMap["config"]},
								{Name: "queryengine", Image: "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + versionMap["config"]},
								{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:4.0.2"},
								{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
								{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
								{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
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
						ServiceConfiguration: v1alpha1.ControlConfiguration{
							CassandraInstance: "cassandra1",
							Containers: []*v1alpha1.Container{
								{Name: "control", Image: "registry:5000/contrail-nightly/contrail-controller-control-control:" + versionMap["control"]},
								{Name: "dns", Image: "registry:5000/contrail-nightly/contrail-controller-control-dns:" + versionMap["control"]},
								{Name: "named", Image: "registry:5000/contrail-nightly/contrail-controller-control-named:" + versionMap["control"]},
								{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
								{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
							},
						},
					},
				}},
				ProvisionManager: &v1alpha1.ProvisionManagerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "provmanager1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},
					Spec: v1alpha1.ProvisionManagerServiceSpec{
						CommonConfiguration: v1alpha1.PodConfiguration{
							Replicas: &replicas,
						},
						ServiceConfiguration: v1alpha1.ProvisionmanagerManagerServiceConfiguration{
							ProvisionManagerConfiguration: v1alpha1.ProvisionManagerConfiguration{
								Containers: []*v1alpha1.Container{
									{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
									{Name: "provisioner", Image: "registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + versionMap["contrail-operator-provisioner"]},
								},
							},
						},
					},
				},
				Kubemanagers: []*v1alpha1.KubemanagerService{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubemanager1",
						Namespace: namespace,
						Labels:    map[string]string{"contrail_cluster": "cluster1"},
					},

					Spec: v1alpha1.KubemanagerServiceSpec{
						ServiceConfiguration: v1alpha1.KubemanagerManagerServiceConfiguration{
							CassandraInstance: "cassandra1",
							ZookeeperInstance: "zookeeper1",
							KubemanagerConfiguration: v1alpha1.KubemanagerConfiguration{
								Containers: []*v1alpha1.Container{
									{Name: "kubemanager", Image: "registry:5000/contrail-nightly/contrail-kubernetes-kube-manager:" + versionMap["kubemanager"]},
									{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
									{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
								},
							},
						},
					},
				}},
			},
		},
	}
}

func upgradeZookeeper(t *testing.T, f *test.Framework, ctx *test.TestCtx, namespace, name string) error {
	instance := &v1alpha1.Manager{}
	err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return err
	}
	zkContainer := utils.GetContainerFromList("zookeeper", instance.Spec.Services.Zookeepers[0].Spec.ServiceConfiguration.Containers)
	zkContainer.Image = "registry:5000/common-docker-third-party/contrail/zookeeper:" + targetVersionMap["zookeeper"]
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
