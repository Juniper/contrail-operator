package v1alpha1

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	mRand "math/rand"

	yaml "gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Juniper/contrail-operator/pkg/certificates"
)

var src = mRand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type MonitorConfig struct {
	APIServerList  []string          `yaml:"apiServerList,omitempty"`
	Encryption     MonitorEncryption `yaml:"encryption,omitempty"`
	NodeType       string            `yaml:"nodeType,omitempty"`
	Interval       int64             `yaml:"interval,omitempty"`
	Hostname       string            `yaml:"hostname,omitempty"`
	InCluster      *bool             `yaml:"inCluster,omitempty"`
	KubeConfigPath string            `yaml:"kubeConfigPath,omitempty"`
	NodeName       string            `yaml:"nodeName,omitempty"`
	Namespace      string            `yaml:"namespace,omitempty"`
	PodName        string            `yaml:"podName,omitempty"`
}

type MonitorEncryption struct {
	CA       *string `yaml:"ca,omitempty"`
	Cert     *string `yaml:"cert,omitempty"`
	Key      *string `yaml:"key,omitempty"`
	Insecure bool    `yaml:"insecure,omitempty"`
}

// Container defines name, image and command.
// +k8s:openapi-gen=true
type Container struct {
	Name    string   `json:"name,omitempty"`
	Image   string   `json:"image,omitempty"`
	Command []string `json:"command,omitempty"`
}

// ServiceStatus provides information on the current status of the service.
// +k8s:openapi-gen=true
type ServiceStatus struct {
	Name    *string `json:"name,omitempty"`
	Active  *bool   `json:"active,omitempty"`
	Created *bool   `json:"created,omitempty"`
}

// ActiveStatus signals the current status
type ActiveStatus struct {
	Active *bool `json:"active,omitempty"`
}

// ServiceInstance is the interface to manage instances.
type ServiceInstance interface {
	Get(client.Client, reconcile.Request) error
	Update(client.Client) error
	Create(client.Client) error
	Delete(client.Client) error
}

// PodConfiguration is the common services struct.
// +k8s:openapi-gen=true
type PodConfiguration struct {
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	// Host networking requested for this pod. Use the host's network namespace.
	// If this option is set, the ports that will be used must be specified.
	// Default to false.
	// +k8s:conversion-gen=false
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified.
	// +optional
	// +patchMergeKey=ip
	// +patchStrategy=merge
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty" patchStrategy:"merge" patchMergeKey:"ip" protobuf:"bytes,23,rep,name=hostAliases"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
}

//GetReplicas is used to get number of desired pods.
func (cc *PodConfiguration) GetReplicas() int32 {
	if cc.Replicas != nil {
		return *cc.Replicas
	}
	return int32(1)
}

func (ss *ServiceStatus) ready() bool {
	if ss == nil {
		return false
	}
	if ss.Active == nil {
		return false
	}

	return *ss.Active

}

func CreateAccount(accountName string, namespace string, client client.Client, scheme *runtime.Scheme, owner v1.Object) error {

	serviceAccountName := "serviceaccount-" + accountName
	clusterRoleName := "clusterrole-" + accountName
	clusterRoleBindingName := "clusterrolebinding-" + accountName
	secretName := "secret-" + accountName

	existingServiceAccount := &corev1.ServiceAccount{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: serviceAccountName, Namespace: namespace}, existingServiceAccount)
	if err != nil && k8serrors.IsNotFound(err) {
		serviceAccount := &corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		}
		err = controllerutil.SetControllerReference(owner, serviceAccount, scheme)
		if err != nil {
			return err
		}
		if err = client.Create(context.TODO(), serviceAccount); err != nil && !k8serrors.IsAlreadyExists(err) {
			return err
		}
	}

	existingSecret := &corev1.Secret{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, existingSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
				Annotations: map[string]string{
					"kubernetes.io/service-account.name": serviceAccountName,
				},
			},
			Type: corev1.SecretType("kubernetes.io/service-account-token"),
		}
		err = controllerutil.SetControllerReference(owner, secret, scheme)
		if err != nil {
			return err
		}
		if err = client.Create(context.TODO(), secret); err != nil {
			return err
		}
	}

	existingClusterRole := &rbacv1.ClusterRole{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName}, existingClusterRole)
	if err != nil && k8serrors.IsNotFound(err) {
		clusterRole := &rbacv1.ClusterRole{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac/v1",
				Kind:       "ClusterRole",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterRoleName,
				Namespace: namespace,
			},
			Rules: []rbacv1.PolicyRule{{
				Verbs: []string{
					"*",
				},
				APIGroups: []string{
					"*",
				},
				Resources: []string{
					"*",
				},
			}},
		}
		if err = client.Create(context.TODO(), clusterRole); err != nil {
			return err
		}
	}

	existingClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleBindingName}, existingClusterRoleBinding)
	if err != nil && k8serrors.IsNotFound(err) {
		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac/v1",
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterRoleBindingName,
				Namespace: namespace,
			},
			Subjects: []rbacv1.Subject{{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			}},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     clusterRoleName,
			},
		}
		if err = client.Create(context.TODO(), clusterRoleBinding); err != nil {
			return err
		}
	}
	return nil
}

func StatusMonitorConfig(hostname string, configNodeList []string, podIP, nodeType, nodeName, namespace, podName string) (string, error) {
	cert := "/etc/certificates/server-" + podIP + ".crt"
	key := "/etc/certificates/server-key-" + podIP + ".pem"
	ca := certificates.SignerCAFilepath
	inCluster := true
	monitorConfig := MonitorConfig{
		APIServerList: configNodeList,
		Encryption: MonitorEncryption{
			CA:       &ca,
			Cert:     &cert,
			Key:      &key,
			Insecure: true,
		},
		NodeType:  nodeType,
		Hostname:  hostname,
		Interval:  10,
		InCluster: &inCluster,
		NodeName:  nodeName,
		Namespace: namespace,
		PodName:   podName,
	}

	monitorYaml, err := yaml.Marshal(monitorConfig)
	if err != nil {
		return "", err
	}
	return string(monitorYaml), nil
}

// SetPodsToReady sets the status label of a POD to ready.
func SetPodsToReady(podList *corev1.PodList, client client.Client) error {
	for _, pod := range podList.Items {
		podObject := &corev1.Pod{}
		if err := client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, podObject); err != nil {
			return err
		}
		pod.ObjectMeta.Labels["status"] = "ready"
		if err := client.Update(context.TODO(), &pod); err != nil {
			return err
		}
	}
	return nil
}

// +k8s:deepcopy-gen=false
type podAltIPsRetriver func(pod corev1.Pod) []string

// +k8s:deepcopy-gen=false
type PodAlternativeIPs struct {
	// Function which operate over pod object
	// to retrieve additional IP addresses used
	// by this pod.
	Retriever podAltIPsRetriver
	// ServiceIP through which pod can be reached.
	ServiceIP string
}

// PodsCertSubjects iterates over passed list of pods and for every pod prepares certificate subject
// which can be later used for generating certificate for given pod.
func PodsCertSubjects(podList *corev1.PodList, hostNetwork *bool, podAltIPs PodAlternativeIPs) []certificates.CertificateSubject {
	var pods []certificates.CertificateSubject
	useNodeName := true
	if hostNetwork != nil {
		useNodeName = *hostNetwork
	}
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodFailed {
			continue
		}
		var hostname string
		if useNodeName {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		var alternativeIPs []string
		if podAltIPs.ServiceIP != "" {
			alternativeIPs = append(alternativeIPs, podAltIPs.ServiceIP)
		}
		if podAltIPs.Retriever != nil {
			if altIPs := podAltIPs.Retriever(pod); len(altIPs) > 0 {
				alternativeIPs = append(alternativeIPs, altIPs...)
			}
		}
		podInfo := certificates.NewSubject(pod.Name, hostname, pod.Status.PodIP, alternativeIPs)
		pods = append(pods, podInfo)
	}
	return pods
}

// CreateConfigMap creates a config map based on the instance type.
func CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request,
	instanceType string,
	object v1.Object) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: request.Namespace}, configMap)
	// TODO: Bug. If config map exists without labels and references, they won't be updated
	if err != nil {
		if k8serrors.IsNotFound(err) {
			configMap.SetName(configMapName)
			configMap.SetNamespace(request.Namespace)
			configMap.SetLabels(map[string]string{"contrail_manager": instanceType,
				instanceType: request.Name})
			configMap.Data = make(map[string]string)
			if err = controllerutil.SetControllerReference(object, configMap, scheme); err != nil {
				return nil, err
			}
			if err = client.Create(context.TODO(), configMap); err != nil && !k8serrors.IsAlreadyExists(err) {
				return nil, err
			}
		}
	}
	return configMap, nil
}

// CreateSecret creates a secret based on the instance type.
func CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request,
	instanceType string,
	object v1.Object) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: request.Namespace}, secret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			secret.SetName(secretName)
			secret.SetNamespace(request.Namespace)
			secret.SetLabels(map[string]string{"contrail_manager": instanceType,
				instanceType: request.Name})
			var data = make(map[string][]byte)
			secret.Data = data
			if err = controllerutil.SetControllerReference(object, secret, scheme); err != nil {
				return nil, err
			}
			if err = client.Create(context.TODO(), secret); err != nil && !k8serrors.IsAlreadyExists(err) {
				return nil, err
			}
		}
	}
	return secret, nil
}

// CurrentConfigMapExists checks if a current configuration exists and returns it.
func CurrentConfigMapExists(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (corev1.ConfigMap, bool) {
	configMapExists := false
	configMap := &corev1.ConfigMap{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: request.Namespace}, configMap); err == nil {
		if len(configMap.Data) > 0 {
			configMapExists = true
		}
	}
	return *configMap, configMapExists
}

// PrepareSTS prepares the intended podList.
func PrepareSTS(sts *appsv1.StatefulSet,
	commonConfiguration *PodConfiguration,
	instanceType string,
	request reconcile.Request,
	scheme *runtime.Scheme,
	object v1.Object,
	client client.Client,
	waitForInit bool) error {
	SetSTSCommonConfiguration(sts, commonConfiguration)
	if waitForInit {
		sts.Spec.PodManagementPolicy = appsv1.PodManagementPolicyType("Parallel")
	} else {
		sts.Spec.PodManagementPolicy = appsv1.PodManagementPolicyType("OrderedReady")
	}
	sts.SetName(request.Name + "-" + instanceType + "-statefulset")
	sts.SetNamespace(request.Namespace)
	sts.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	sts.Spec.Selector.MatchLabels = map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name}
	sts.Spec.Template.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	if err := controllerutil.SetControllerReference(object, sts, scheme); err != nil {
		return err
	}
	return nil
}

// SetDeploymentCommonConfiguration takes common configuration parameters
// and applies it to the deployment.
func SetDeploymentCommonConfiguration(deployment *appsv1.Deployment,
	commonConfiguration *PodConfiguration) *appsv1.Deployment {
	deployment.Spec.Replicas = commonConfiguration.Replicas
	if len(commonConfiguration.Tolerations) > 0 {
		deployment.Spec.Template.Spec.Tolerations = commonConfiguration.Tolerations
	}
	if len(commonConfiguration.NodeSelector) > 0 {
		deployment.Spec.Template.Spec.NodeSelector = commonConfiguration.NodeSelector
	}
	if commonConfiguration.HostNetwork != nil {
		deployment.Spec.Template.Spec.HostNetwork = *commonConfiguration.HostNetwork
	} else {
		deployment.Spec.Template.Spec.HostNetwork = false
	}

	if len(commonConfiguration.HostAliases) > 0 {
		deployment.Spec.Template.Spec.HostAliases = commonConfiguration.HostAliases
	}

	if len(commonConfiguration.ImagePullSecrets) > 0 {
		imagePullSecretList := []corev1.LocalObjectReference{}
		for _, imagePullSecretName := range commonConfiguration.ImagePullSecrets {
			imagePullSecret := corev1.LocalObjectReference{
				Name: imagePullSecretName,
			}
			imagePullSecretList = append(imagePullSecretList, imagePullSecret)
		}
		deployment.Spec.Template.Spec.ImagePullSecrets = imagePullSecretList
	}
	return deployment
}

// SetSTSCommonConfiguration takes common configuration parameters
// and applies it to the pod.
func SetSTSCommonConfiguration(sts *appsv1.StatefulSet,
	commonConfiguration *PodConfiguration) {
	sts.Spec.Replicas = commonConfiguration.Replicas
	if len(commonConfiguration.Tolerations) > 0 {
		sts.Spec.Template.Spec.Tolerations = commonConfiguration.Tolerations
	}
	if len(commonConfiguration.NodeSelector) > 0 {
		sts.Spec.Template.Spec.NodeSelector = commonConfiguration.NodeSelector
	}
	if commonConfiguration.HostNetwork != nil {
		sts.Spec.Template.Spec.HostNetwork = *commonConfiguration.HostNetwork
	} else {
		sts.Spec.Template.Spec.HostNetwork = false
	}

	if len(commonConfiguration.HostAliases) > 0 {
		sts.Spec.Template.Spec.HostAliases = commonConfiguration.HostAliases
	}

	if len(commonConfiguration.ImagePullSecrets) > 0 {
		imagePullSecretList := []corev1.LocalObjectReference{}
		for _, imagePullSecretName := range commonConfiguration.ImagePullSecrets {
			imagePullSecret := corev1.LocalObjectReference{
				Name: imagePullSecretName,
			}
			imagePullSecretList = append(imagePullSecretList, imagePullSecret)
		}
		sts.Spec.Template.Spec.ImagePullSecrets = imagePullSecretList
	}
}

// AddVolumesToIntendedSTS adds volumes to a deployment.
func AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	volumeList := sts.Spec.Template.Spec.Volumes
	for configMapName, volumeName := range volumeConfigMapMap {
		volume := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName,
					},
				},
			},
		}
		volumeList = append(volumeList, volume)
	}
	sts.Spec.Template.Spec.Volumes = volumeList
}

// AddSecretVolumesToIntendedSTS adds volumes to a deployment.
func AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeSecretMap map[string]string) {
	volumeList := sts.Spec.Template.Spec.Volumes
	for secretName, volumeName := range volumeSecretMap {
		volume := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		}
		volumeList = append(volumeList, volume)
	}
	sts.Spec.Template.Spec.Volumes = volumeList
}

// AddSecretVolumesToIntendedDS adds volumes to a deployment.
func AddSecretVolumesToIntendedDS(ds *appsv1.DaemonSet, volumeSecretMap map[string]string) {
	volumeList := ds.Spec.Template.Spec.Volumes
	for secretName, volumeName := range volumeSecretMap {
		volume := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		}
		volumeList = append(volumeList, volume)
	}
	ds.Spec.Template.Spec.Volumes = volumeList
}

// CreateSTS creates the STS.
func CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	foundSTS := &appsv1.StatefulSet{}
	err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-statefulset", Namespace: request.Namespace}, foundSTS)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			sts.Spec.Template.ObjectMeta.Labels["version"] = "1"
			if err = reconcileClient.Create(context.TODO(), sts); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateSTS updates the STS.
func UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	currentSTS := &appsv1.StatefulSet{}
	err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-statefulset", Namespace: request.Namespace}, currentSTS)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	replicasChanged := false
	if *sts.Spec.Replicas != *currentSTS.Spec.Replicas {
		replicasChanged = true
	}
	imagesChanged := false
	for _, intendedContainer := range sts.Spec.Template.Spec.Containers {
		for _, currentContainer := range currentSTS.Spec.Template.Spec.Containers {
			if intendedContainer.Name == currentContainer.Name {
				if intendedContainer.Image != currentContainer.Image {
					imagesChanged = true
				}
			}
		}
	}
	if imagesChanged || replicasChanged {
		if strategy == "deleteFirst" {
			versionInt, _ := strconv.Atoi(currentSTS.Spec.Template.ObjectMeta.Labels["version"])
			newVersion := versionInt + 1
			sts.Spec.Template.ObjectMeta.Labels["version"] = strconv.Itoa(newVersion)
		} else {
			sts.Spec.Template.ObjectMeta.Labels["version"] = currentSTS.Spec.Template.ObjectMeta.Labels["version"]
		}
		if err = reconcileClient.Update(context.TODO(), sts); err != nil {
			return err
		}
	}
	return nil
}

// SetInstanceActive sets the instance to active.
func SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request, object runtime.Object) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace: request.Namespace},
		sts); err != nil {
		return err
	}
	active := false
	if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
		active = true
	}

	*activeStatus = active
	if err := client.Status().Update(context.TODO(), object); err != nil {
		return err
	}
	return nil
}

// RandomString creates a random string of size
func RandomString(size int) string {
	b := make([]byte, size)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := size-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func getPodsHostname(c client.Reader, pod *corev1.Pod) (string, error) {
	if !pod.Spec.HostNetwork {
		return pod.Spec.Hostname, nil
	}
	n := corev1.Node{}
	if err := c.Get(context.Background(), types.NamespacedName{Name: pod.Spec.NodeName}, &n); err != nil {
		return "", err
	}

	for _, a := range n.Status.Addresses {
		if a.Type == corev1.NodeHostName {
			return a.Address, nil
		}
	}

	return "", errors.New("couldn't get pods hostname")
}

func getPodInitStatus(reconcileClient client.Client,
	podList *corev1.PodList,
	getHostname bool,
	getInterface bool,
	getMac bool,
	getPrefix bool,
	getGateway bool,
	waitForInit bool) (map[string]string, error) {
	var podNameIPMap = make(map[string]string)
	for idx, pod := range podList.Items {
		if !waitForInit && pod.Status.Phase != "Running" && pod.Status.Phase != "Pending" {
			return map[string]string{}, nil
		}
		if pod.Status.PodIP == "" {
			return map[string]string{}, nil
		}
		if getHostname || getInterface || getMac || getPrefix {
			for _, initStatus := range pod.Status.InitContainerStatuses {
				if initStatus.Name == "init" {
					if initStatus.State.Terminated == nil {
						if initStatus.State.Running != nil {
							annotationMap := pod.GetAnnotations()
							if annotationMap == nil {
								annotationMap = make(map[string]string)
							}
							if getHostname {

								hostname, err := getPodsHostname(reconcileClient, &pod)
								if err != nil {
									return map[string]string{}, err
								}
								annotationMap["hostname"] = hostname
							}
							if getInterface {
								command := []string{"/bin/sh", "-c", "ifconfig | sed -n '/addr:" + pod.Status.PodIP + "/{g;h;p};h;x'  |awk '{print $1}'"}
								physicalInterface, _, err := ExecToPodThroughAPI(command, "init", pod.Name, pod.Namespace, nil)
								if err != nil {
									return map[string]string{}, fmt.Errorf("failed getting interface")
								}
								annotationMap["physicalInterface"] = strings.Trim(physicalInterface, "\n")
							}
							if getGateway {
								command := []string{"/bin/sh", "-c", "ip route get 1.1.1.1 |grep -v cache |sed -e 's/.* via \\(.*\\) dev.*/\\1/'"}
								gateway, _, err := ExecToPodThroughAPI(command, "init", pod.Name, pod.Namespace, nil)
								if err != nil {
									return map[string]string{}, fmt.Errorf("failed getting gateway")
								}
								annotationMap["gateway"] = strings.Trim(gateway, "\n")
							}
							if getMac {
								command := []string{"/bin/sh", "-c", "ip address show | sed -n '/inet " + pod.Status.PodIP + "\\//{g;h;p};h;x' |awk '{print $2}'"}
								physicalInterfaceMac, _, err := ExecToPodThroughAPI(command, "init", pod.Name, pod.Namespace, nil)
								if err != nil {
									return map[string]string{}, fmt.Errorf("failed getting mac")
								}
								annotationMap["physicalInterfaceMac"] = strings.Trim(physicalInterfaceMac, "\n")
							}
							if getPrefix {
								command := []string{"/bin/sh", "-c", "ip addr sh |sed -n 's/.*" + pod.Status.PodIP + "\\/\\([^ ]*\\).*/\\1/p'"}
								prefixLength, _, err := ExecToPodThroughAPI(command, "init", pod.Name, pod.Namespace, nil)
								if err != nil {
									return map[string]string{}, fmt.Errorf("failed getting prefix")
								}
								annotationMap["prefixLength"] = strings.Trim(prefixLength, "\n")
							}

							if cidr, ok := pod.Annotations["dataSubnet"]; ok {
								if cidr != "" {
									command := []string{"/bin/sh", "-c", "ip r | grep " + cidr + " | awk -F' ' '{print $NF}'"}
									addr, _, err := ExecToPodThroughAPI(command, "init", pod.Name, pod.Namespace, nil)
									if err != nil {
										return map[string]string{}, fmt.Errorf("failed getting ip address from data subnet")
									}
									ip := strings.Trim(addr, "\n")
									if net.ParseIP(ip) != nil {
										annotationMap["dataSubnetIP"] = ip
									} else {
										return map[string]string{}, fmt.Errorf("no valid ip from data subnet")
									}
								}
							}

							podList.Items[idx].SetAnnotations(annotationMap)
							(&podList.Items[idx]).SetAnnotations(annotationMap)
							foundPod := &corev1.Pod{}
							err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)
							if err != nil {
								return map[string]string{}, err
							}
							foundPod.SetAnnotations(annotationMap)
							err = reconcileClient.Update(context.TODO(), foundPod)
							if err != nil {
								return map[string]string{}, err
							}
							podList.Items[idx] = *foundPod
						} else {
							return map[string]string{}, nil
						}
					}
				}
			}
		}
		podNameIPMap[pod.Name] = pod.Status.PodIP
	}
	return podNameIPMap, nil
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func PodIPListAndIPMapFromInstance(instanceType string,
	commonConfiguration *PodConfiguration,
	request reconcile.Request,
	reconcileClient client.Client,
	waitForInit bool,
	getHostname bool,
	getInterface bool,
	getMac bool,
	getPrefix bool,
	getGateway bool) (*corev1.PodList, map[string]string, error) {
	var podNameIPMap = make(map[string]string)
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
	podList := &corev1.PodList{}
	err := reconcileClient.List(context.TODO(), podList, listOps)
	if err != nil {
		return &corev1.PodList{}, map[string]string{}, err
	}
	if len(podList.Items) > 0 {
		if waitForInit {
			if int32(len(podList.Items)) == *commonConfiguration.Replicas {
				podNameIPMap, err = getPodInitStatus(reconcileClient, podList, getHostname, getInterface, getMac, getPrefix, getGateway, waitForInit)
				if err != nil {
					return podList, podNameIPMap, err
				}
			}
			if int32(len(podNameIPMap)) != *commonConfiguration.Replicas {
				return &corev1.PodList{}, map[string]string{}, nil
			}
		} else if len(podList.Items) > 0 {
			podNameIPMap, err = getPodInitStatus(reconcileClient, podList, getHostname, getInterface, getMac, getPrefix, getGateway, waitForInit)
			if err != nil {
				return podList, podNameIPMap, err
			}
		}
		return podList, podNameIPMap, nil
	}
	return &corev1.PodList{}, map[string]string{}, nil
}

// NewCassandraClusterConfiguration gets a struct containing various representations of Cassandra nodes string.
func NewCassandraClusterConfiguration(name string, namespace string, client client.Client) (CassandraClusterConfiguration, error) {
	var cassandraCluster CassandraClusterConfiguration
	var cassandraNodes []string
	cassandraInstance := &Cassandra{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cassandraInstance)
	if err != nil {
		return cassandraCluster, err
	}
	for _, ip := range cassandraInstance.Status.Nodes {
		cassandraNodes = append(cassandraNodes, ip)
	}
	sort.SliceStable(cassandraNodes, func(i, j int) bool { return cassandraNodes[i] < cassandraNodes[j] })
	cassandraConfig := cassandraInstance.ConfigurationParameters()
	endpoint := cassandraInstance.Status.ClusterIP + ":" + strconv.Itoa(*cassandraConfig.Port)
	cassandraCluster = CassandraClusterConfiguration{
		Port:         *cassandraConfig.Port,
		CQLPort:      *cassandraConfig.CqlPort,
		JMXPort:      *cassandraConfig.JmxLocalPort,
		ServerIPList: cassandraNodes,
		Endpoint:     endpoint,
	}
	return cassandraCluster, nil
}

// NewControlClusterConfiguration gets a struct containing various representations of Control nodes string.
func NewControlClusterConfiguration(name string, role string, namespace string, myclient client.Client) (ControlClusterConfiguration, error) {
	var controlNodes []string
	var controlCluster ControlClusterConfiguration
	var controlConfig ControlConfiguration
	if name != "" {
		controlInstance := &Control{}
		err := myclient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, controlInstance)
		if err != nil {
			return controlCluster, err
		}
		for _, ip := range controlInstance.Status.Nodes {
			controlNodes = append(controlNodes, ip)
		}
		controlConfig = controlInstance.ConfigurationParameters()
	}
	if role != "" {
		labelSelector := labels.SelectorFromSet(map[string]string{"control_role": role})
		listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
		controlList := &ControlList{}
		err := myclient.List(context.TODO(), controlList, listOps)
		if err != nil {
			return controlCluster, err
		}
		if len(controlList.Items) > 0 {
			for _, ip := range controlList.Items[0].Status.Nodes {
				controlNodes = append(controlNodes, ip)
			}
		} else {
			return controlCluster, err
		}
		controlConfig = controlList.Items[0].ConfigurationParameters()
	}
	sort.SliceStable(controlNodes, func(i, j int) bool { return controlNodes[i] < controlNodes[j] })
	controlCluster = ControlClusterConfiguration{
		XMPPPort:            *controlConfig.XMPPPort,
		BGPPort:             *controlConfig.BGPPort,
		DNSPort:             *controlConfig.DNSPort,
		DNSIntrospectPort:   *controlConfig.DNSIntrospectPort,
		ControlServerIPList: controlNodes,
	}

	return controlCluster, nil
}

// NewZookeeperClusterConfiguration gets a struct containing various representations of Zookeeper nodes string.
func NewZookeeperClusterConfiguration(name string, namespace string, client client.Client) (ZookeeperClusterConfiguration, error) {
	var zookeeperNodes []string
	var zookeeperCluster ZookeeperClusterConfiguration
	zookeeperInstance := &Zookeeper{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, zookeeperInstance)
	if err != nil {
		return zookeeperCluster, err
	}
	for _, ip := range zookeeperInstance.Status.Nodes {
		zookeeperNodes = append(zookeeperNodes, ip)

	}
	zookeeperConfig := zookeeperInstance.ConfigurationParameters()
	sort.SliceStable(zookeeperNodes, func(i, j int) bool { return zookeeperNodes[i] < zookeeperNodes[j] })
	zookeeperCluster = ZookeeperClusterConfiguration{
		ClientPort:   *zookeeperConfig.ClientPort,
		ServerIPList: zookeeperNodes,
	}

	return zookeeperCluster, nil
}

// NewRabbitmqClusterConfiguration gets a struct containing various representations of Rabbitmq nodes string.
func NewRabbitmqClusterConfiguration(name string, namespace string, myclient client.Client) (RabbitmqClusterConfiguration, error) {
	var rabbitmqNodes []string
	var rabbitmqCluster RabbitmqClusterConfiguration
	secret := ""
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	rabbitmqList := &RabbitmqList{}
	err := myclient.List(context.TODO(), rabbitmqList, listOps)
	if err != nil {
		return rabbitmqCluster, err
	}
	var rabbitmqConfig RabbitmqConfiguration
	if len(rabbitmqList.Items) > 0 {
		for _, ip := range rabbitmqList.Items[0].Status.Nodes {
			rabbitmqNodes = append(rabbitmqNodes, ip)
		}
		rabbitmqConfig = rabbitmqList.Items[0].ConfigurationParameters()
		secret = rabbitmqList.Items[0].Status.Secret
	}
	sort.SliceStable(rabbitmqNodes, func(i, j int) bool { return rabbitmqNodes[i] < rabbitmqNodes[j] })
	rabbitmqCluster = RabbitmqClusterConfiguration{
		Port:         *rabbitmqConfig.Port,
		SSLPort:      *rabbitmqConfig.SSLPort,
		ServerIPList: rabbitmqNodes,
		Secret:       secret,
	}
	return rabbitmqCluster, nil
}

// NewConfigClusterConfiguration gets a struct containing various representations of Config nodes string.
func NewConfigClusterConfiguration(name string, namespace string, myclient client.Client) (ConfigClusterConfiguration, error) {
	var configNodes []string
	var configCluster ConfigClusterConfiguration
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	configList := &ConfigList{}
	err := myclient.List(context.TODO(), configList, listOps)
	if err != nil {
		return configCluster, err
	}

	var authMode AuthenticationMode
	var apiServerPort, analyticsPort, collectorPort, redisPort int

	if len(configList.Items) > 0 {
		for _, ip := range configList.Items[0].Status.Nodes {
			configNodes = append(configNodes, ip)
		}
		configConfig := configList.Items[0].ConfigurationParameters()
		authMode = configConfig.AuthMode
		apiServerPort = *configConfig.APIPort
		analyticsPort = *configConfig.AnalyticsPort
		collectorPort = *configConfig.CollectorPort
		redisPort = *configConfig.RedisPort
	}
	sort.SliceStable(configNodes, func(i, j int) bool { return configNodes[i] < configNodes[j] })
	configCluster = ConfigClusterConfiguration{
		APIServerPort:         apiServerPort,
		APIServerIPList:       configNodes,
		AnalyticsServerPort:   analyticsPort,
		AnalyticsServerIPList: configNodes,
		CollectorPort:         collectorPort,
		CollectorServerIPList: configNodes,
		RedisPort:             redisPort,
		AuthMode:              authMode,
	}
	return configCluster, nil
}

// WebUIClusterConfiguration defines all configuration knobs used to write the config file.
type WebUIClusterConfiguration struct {
	AdminUsername string
	AdminPassword string
}

// CommandClusterConfiguration defines all configuration knobs used to write the config file.
type CommandClusterConfiguration struct {
	AdminUsername string
	AdminPassword string
}

// ConfigClusterConfiguration  stores all information about service's endpoints
// under the Contrail Config
type ConfigClusterConfiguration struct {
	APIServerPort         int                `json:"apiServerPort,omitempty"`
	APIServerIPList       []string           `json:"apiServerIPList,omitempty"`
	AnalyticsServerPort   int                `json:"analyticsServerPort,omitempty"`
	AnalyticsServerIPList []string           `json:"analyticsServerIPList,omitempty"`
	CollectorPort         int                `json:"collectorPort,omitempty"`
	CollectorServerIPList []string           `json:"collectorServerIPList,omitempty"`
	RedisPort             int                `json:"redisPort,omitempty"`
	AuthMode              AuthenticationMode `json:"authMode,omitempty"`
}

// FillWithDefaultValues sets the default port values if they are set to the
// zero value
func (c *ConfigClusterConfiguration) FillWithDefaultValues() {
	if c.APIServerPort == 0 {
		c.APIServerPort = ConfigApiPort
	}
	if c.AnalyticsServerPort == 0 {
		c.AnalyticsServerPort = AnalyticsApiPort
	}
	if c.CollectorPort == 0 {
		c.CollectorPort = CollectorPort
	}
	if c.RedisPort == 0 {
		c.RedisPort = RedisServerPort
	}
	if c.AuthMode == "" {
		c.AuthMode = AuthenticationModeNoAuth
	}
}

// ControlClusterConfiguration stores all information about services' endpoints
// under the Contrail Control
type ControlClusterConfiguration struct {
	XMPPPort            int      `json:"xmppPort,omitempty"`
	BGPPort             int      `json:"bgpPort,omitempty"`
	DNSPort             int      `json:"dnsPort,omitempty"`
	DNSIntrospectPort   int      `json:"dnsIntrospectPort,omitempty"`
	ControlServerIPList []string `json:"controlServerIPList,omitempty"`
}

// FillWithDefaultValues sets the default port values if they are set to the
// zero value
func (c *ControlClusterConfiguration) FillWithDefaultValues() {
	if c.XMPPPort == 0 {
		c.XMPPPort = XmppServerPort
	}
	if c.BGPPort == 0 {
		c.BGPPort = BgpPort
	}
	if c.DNSPort == 0 {
		c.DNSPort = DnsServerPort
	}
	if c.DNSIntrospectPort == 0 {
		c.DNSIntrospectPort = DnsIntrospectPort
	}
}

// ZookeeperClusterConfiguration stores all information about Zookeeper's endpoints.
type ZookeeperClusterConfiguration struct {
	ClientPort   int      `json:"clientPort,omitempty"`
	ServerPort   int      `json:"serverPort,omitempty"`
	ElectionPort int      `json:"electionPort,omitempty"`
	ServerIPList []string `json:"serverIPList,omitempty"`
}

// FillWithDefaultValues fills Zookeeper config with default values
func (c *ZookeeperClusterConfiguration) FillWithDefaultValues() {
	if c.ClientPort == 0 {
		c.ClientPort = ZookeeperPort
	}
	if c.ElectionPort == 0 {
		c.ElectionPort = ZookeeperElectionPort
	}
	if c.ServerPort == 0 {
		c.ServerPort = ZookeeperServerPort
	}
}

// RabbitmqClusterConfiguration stores all information about Rabbitmq's endpoints.
type RabbitmqClusterConfiguration struct {
	Port         int      `json:"port,omitempty"`
	SSLPort      int      `json:"sslPort,omitempty"`
	ServerIPList []string `json:"serverIPList,omitempty"`
	Secret       string   `json:"secret,omitempty"`
}

// FillWithDefaultValues fills Rabbitmq config with default values
func (c *RabbitmqClusterConfiguration) FillWithDefaultValues() {
	if c.Port == 0 {
		c.Port = RabbitmqNodePort
	}
	if c.SSLPort == 0 {
		c.SSLPort = RabbitmqNodePortSSL
	}
}

// CassandraClusterConfiguration stores all information about Cassandra's endpoints.
type CassandraClusterConfiguration struct {
	Port         int      `json:"port,omitempty"`
	CQLPort      int      `json:"cqlPort,omitempty"`
	JMXPort      int      `json:"jmxPort,omitempty"`
	ServerIPList []string `json:"serverIPList,omitempty"`
	Endpoint     string   `json:"endpoint,omitempty"`
}

// FillWithDefaultValues fills Cassandra config with default values
func (c *CassandraClusterConfiguration) FillWithDefaultValues() {
	if c.CQLPort == 0 {
		c.CQLPort = CassandraCqlPort
	}
	if c.JMXPort == 0 {
		c.JMXPort = CassandraJmxLocalPort
	}
	if c.Port == 0 {
		c.Port = CassandraPort
	}
}

// VrouterClusterConfiguration defines all configuration knobs used to write the config file.
type VrouterClusterConfiguration struct {
	PhysicalInterface string
	Gateway           string
	MetaDataSecret    string
}
