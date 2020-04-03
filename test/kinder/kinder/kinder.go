package kinder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
)

// KindCluster configures the Kind Cluster
type KindCluster struct {
	Name           string
	Nodes          *int
	KubeConfig     string
	KubeConfigFile string
	Stop           chan bool
	Running        bool
	Provider       *cluster.Provider
}

// Start creates and starts a Kind Cluster
func (k *KindCluster) Start() {
	var (
		err error
	)
	nodeCount := 1
	if k.Nodes != nil {
		nodeCount = *k.Nodes
	}

	nodeList := []v1alpha4.Node{}
	for i := 0; i < nodeCount; i++ {
		node := v1alpha4.Node{
			Role: v1alpha4.NodeRole("control-plane"),
		}
		nodeList = append(nodeList, node)
	}
	config := &v1alpha4.Cluster{
		Nodes: nodeList,
	}
	clusterCreateOption := cluster.CreateWithV1Alpha4Config(config)

	newProvider := cluster.NewProvider()
	kubeConfFile, err := ioutil.TempFile(".", "kubeconf")
	if err != nil {
		log.Fatal(err)
	}

	k.KubeConfigFile = kubeConfFile.Name()

	k.Stop = make(chan bool)
	go func() {
		log.Println("creating a cluster")
		if err := newProvider.Create(k.Name, clusterCreateOption); err != nil {
			log.Fatal(err)
		}
	}()
	err = retry(20, 2*time.Second, func() (err error) {
		k.KubeConfig, err = newProvider.KubeConfig(k.Name, false)
		return
	})
	if err != nil {
		log.Fatal(err)
	}
	newProvider.ExportKubeConfig(k.Name, k.KubeConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	k.Provider = newProvider
	k.Running = true
	<-k.Stop
}

// Delete deletes a KindCluster
func (k *KindCluster) Delete(delay time.Duration) {
	log.Printf("Waiting %d msec before deleting cluster", delay)
	time.Sleep(time.Duration(delay))

	err := k.Provider.Delete(k.Name, k.KubeConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(k.KubeConfigFile)
	close(k.Stop)
}

// WaitForReady waits for the Kind Cluster to become ready
func (k *KindCluster) WaitForReady() {
	err := retry(40, 2*time.Second, func() (err error) {
		if !k.Running {
			err = fmt.Errorf("not running")
		}
		return
	})
	if err != nil {
		log.Fatal(err)
	}
	config, _ := clientcmd.BuildConfigFromFlags("", k.KubeConfigFile)
	clientset, _ := kubernetes.NewForConfig(config)
	var nodes *corev1.NodeList
	log.Printf("Waiting for %d master nodes\n", *k.Nodes)
	err = retry(200, 10*time.Second, func() (err error) {
		nodes, _ = clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if len(nodes.Items) != *k.Nodes {
			log.Printf("Expecting %d nodes, have %d nodes", *k.Nodes, len(nodes.Items))
			err = fmt.Errorf("not enough nodes")
		}
		return

	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for master nodes to be \"Ready\"")
	err = retry(200, 10*time.Second, func() (err error) {
		conditionMap := make(map[string]corev1.NodeConditionType)
		nodes, _ = clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			conditionMap[node.GetName()] = "NotReady"
			conditions := node.Status.Conditions
			for _, condition := range conditions {
				if condition.Type == corev1.NodeReady {
					conditionMap[node.GetName()] = "Ready"
				}
			}
		}
		for k, v := range conditionMap {
			log.Printf("Node %s condition %s\n", k, v)
			if v == "NotReady" {
				err = fmt.Errorf("not all nodes ready")
			}
		}
		return
	})
	if err != nil {
		log.Fatal(err)
	}
}

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
