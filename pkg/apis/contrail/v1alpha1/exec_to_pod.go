package v1alpha1

import (
	"bytes"
	"fmt"
	"io"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	// "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const debug = false

// GetDynamicClientFromConfig takes REST config and Create a dynamic client based on that return that client.
func GetDynamicClientFromConfig(config *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

// GetDynamicClient first tries to get REST config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it creates a new Dynamic Client for the given config and returns the client.
func GetDynamicClient() (dynamic.Interface, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return GetDynamicClientFromConfig(config)
}

// GetRESTClient first tries to get a config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it
func GetRESTClient() (*rest.RESTClient, error) {
	config, err := config.GetConfig()
	if err != nil {
		return &rest.RESTClient{}, err
	}

	return rest.RESTClientFor(config)
}

// ExecToPodThroughAPI uninterractively exec to the pod with the command specified.
// :param string command: list of the str which specify the command.
// :param string pod_name: Pod name
// :param string namespace: namespace of the Pod.
// :param io.Reader stdin: Standerd Input if necessary, otherwise `nil`
// :return: string: Output of the command. (STDOUT)
//          string: Errors. (STDERR)
//           error: If any error has occurred otherwise `nil`
func ExecToPodThroughAPI(command []string, containerName, podName, namespace string, stdin io.Reader) (string, string, error) {
	config, err := config.GetConfig()
	if err != nil {
		return "", "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", "", err
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := core_v1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&core_v1.PodExecOptions{
		//Command:   strings.Fields(command),
		Command:   command,
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	if debug {
		fmt.Println("Request URL:", req.URL().String())
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil
}
