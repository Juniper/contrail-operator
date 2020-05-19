package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const debug = true

// Kubernetes is used to create and update meaningful objects
type Kubernetes struct {
	client client.Client
	scheme *runtime.Scheme
}

type object interface {
	GetName() string
	GetUID() types.UID
	GetOwnerReferences() []meta.OwnerReference
	SetOwnerReferences(references []meta.OwnerReference)
	runtime.Object
}

// New is used to create a new Kubernetes
func New(client client.Client, scheme *runtime.Scheme) *Kubernetes {
	return &Kubernetes{
		client: client,
		scheme: scheme,
	}
}

// Owner is used to create Owner object
func (k *Kubernetes) Owner(owner object) *Owner {
	return &Owner{owner: owner, client: k.client, scheme: k.scheme}
}

// ConfigMap is used to create ConfigMap object
func (k *Kubernetes) ConfigMap(name, ownerType string, owner v1.Object) *ConfigMap {
	return &ConfigMap{name: name, ownerType: ownerType, owner: owner, client: k.client, scheme: k.scheme}
}

// Secret is used to create Secret object
func (k *Kubernetes) Secret(name, ownerType string, owner v1.Object) *Secret {
	return &Secret{name: name, ownerType: ownerType, owner: owner, client: k.client, scheme: k.scheme}
}

// getClientConfig first tries to get a config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it returns the same.
func getClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if debug {
			fmt.Printf("Unable to create config. Error: %+v\n", err)
		}
		err1 := err
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}

// getDynamicClientFromConfig takes REST config and Create a dynamic client based on that return that client.
func getDynamicClientFromConfig(config *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("failed creating dynamic client. Error: %+v", err)
		return nil, err
	}
	return dynamicClient, nil
}

// GetDynamicClient first tries to get REST config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it creates a new Dynamic Client for the given config and returns the client.
func GetDynamicClient() (dynamic.Interface, error) {
	config, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	return getDynamicClientFromConfig(config)
}
