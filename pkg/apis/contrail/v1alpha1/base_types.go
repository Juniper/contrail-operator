package v1alpha1

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mRand "math/rand"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	if err != nil && errors.IsNotFound(err) {
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
		if err = client.Create(context.TODO(), serviceAccount); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	existingSecret := &corev1.Secret{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, existingSecret)
	if err != nil && errors.IsNotFound(err) {
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
	if err != nil && errors.IsNotFound(err) {
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
		err = controllerutil.SetControllerReference(owner, clusterRole, scheme)
		if err != nil {
			return err
		}
		if err = client.Create(context.TODO(), clusterRole); err != nil {
			return err
		}
	}

	existingClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleBindingName}, existingClusterRoleBinding)
	if err != nil && errors.IsNotFound(err) {
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
		err = controllerutil.SetControllerReference(owner, clusterRoleBinding, scheme)
		if err != nil {
			return err
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

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
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
		if errors.IsNotFound(err) {
			configMap.SetName(configMapName)
			configMap.SetNamespace(request.Namespace)
			configMap.SetLabels(map[string]string{"contrail_manager": instanceType,
				instanceType: request.Name})
			configMap.Data = make(map[string]string)
			if err = controllerutil.SetControllerReference(object, configMap, scheme); err != nil {
				return nil, err
			}
			if err = client.Create(context.TODO(), configMap); err != nil && !errors.IsAlreadyExists(err) {
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
		if errors.IsNotFound(err) {
			secret.SetName(secretName)
			secret.SetNamespace(request.Namespace)
			secret.SetLabels(map[string]string{"contrail_manager": instanceType,
				instanceType: request.Name})
			var data = make(map[string][]byte)
			secret.Data = data
			if err = controllerutil.SetControllerReference(object, secret, scheme); err != nil {
				return nil, err
			}
			if err = client.Create(context.TODO(), secret); err != nil && !errors.IsAlreadyExists(err) {
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
		if errors.IsNotFound(err) {
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
		if errors.IsNotFound(err) {
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
							var annotationMap = make(map[string]string)
							if getHostname {
								var hostname string
								if pod.Spec.HostNetwork {
									hostname = pod.Spec.NodeName
								} else {
									hostname = pod.Spec.Hostname
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
							podList.Items[idx].SetAnnotations(annotationMap)
							(&podList.Items[idx]).SetAnnotations(annotationMap)
							foundPod := &corev1.Pod{}
							err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: podList.Items[idx].Name, Namespace: podList.Items[idx].Namespace}, foundPod)
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
func NewCassandraClusterConfiguration(name string, namespace string, client client.Client) (*CassandraClusterConfiguration, error) {
	var cassandraNodes []string
	var cassandraCluster CassandraClusterConfiguration
	var port string
	var cqlPort string
	var jmxPort string
	cassandraInstance := &Cassandra{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cassandraInstance)
	if err != nil {
		return &cassandraCluster, err
	}
	for _, ip := range cassandraInstance.Status.Nodes {
		cassandraNodes = append(cassandraNodes, ip)
	}
	cassandraConfigInterface := cassandraInstance.ConfigurationParameters()
	cassandraConfig := cassandraConfigInterface.(CassandraConfiguration)
	port = strconv.Itoa(*cassandraConfig.Port)
	cqlPort = strconv.Itoa(*cassandraConfig.CqlPort)
	jmxPort = strconv.Itoa(*cassandraConfig.JmxLocalPort)
	sort.SliceStable(cassandraNodes, func(i, j int) bool { return cassandraNodes[i] < cassandraNodes[j] })
	serverListCommaSeparated := strings.Join(cassandraNodes, ":"+port+",")
	serverListCommaSeparated = serverListCommaSeparated + ":" + port
	endpoint := cassandraInstance.Status.ClusterIP + ":" + port
	serverListCQLSpaceSeparated := strings.Join(cassandraNodes, ":"+cqlPort+" ")
	serverListCQLSpaceSeparated = serverListCQLSpaceSeparated + ":" + cqlPort
	serverListJMXCommaSeparated := strings.Join(cassandraNodes, ":"+jmxPort+",")
	serverListJMXCommaSeparated = serverListJMXCommaSeparated + ":" + jmxPort
	serverListJMXSpaceSeparated := strings.Join(cassandraNodes, ":"+jmxPort+" ")
	serverListJMXSpaceSeparated = serverListJMXSpaceSeparated + ":" + jmxPort
	serverListCommanSeparatedQuoted := strings.Join(cassandraNodes, "','")
	serverListCommanSeparatedQuoted = "'" + serverListCommanSeparatedQuoted + "'"
	cassandraCluster = CassandraClusterConfiguration{
		Port:                            port,
		CQLPort:                         cqlPort,
		JMXPort:                         jmxPort,
		Endpoint:                        endpoint,
		ServerListCQLSpaceSeparated:     serverListCQLSpaceSeparated,
		ServerListJMXCommaSeparated:     serverListJMXCommaSeparated,
		ServerListJMXSpaceSeparated:     serverListJMXSpaceSeparated,
		ServerListCommanSeparatedQuoted: serverListCommanSeparatedQuoted,
	}
	return &cassandraCluster, nil
}

// NewControlClusterConfiguration gets a struct containing various representations of Control nodes string.
func NewControlClusterConfiguration(name string, role string, namespace string, myclient client.Client) (*ControlClusterConfiguration, error) {
	var controlNodes []string
	var controlCluster ControlClusterConfiguration
	var bgpPort string
	var dnsPort string
	var xmppPort string
	var dnsIntrospectPort string
	var controlConfigInterface interface{}
	if name != "" {
		controlInstance := &Control{}
		err := myclient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, controlInstance)
		if err != nil {
			return &controlCluster, err
		}
		for _, ip := range controlInstance.Status.Nodes {
			controlNodes = append(controlNodes, ip)
		}
		controlConfigInterface = controlInstance.ConfigurationParameters()
	}
	if role != "" {
		labelSelector := labels.SelectorFromSet(map[string]string{"control_role": role})
		listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
		controlList := &ControlList{}
		err := myclient.List(context.TODO(), controlList, listOps)
		if err != nil {
			return &controlCluster, err
		}
		if len(controlList.Items) > 0 {
			for _, ip := range controlList.Items[0].Status.Nodes {
				controlNodes = append(controlNodes, ip)
			}
		} else {
			return &controlCluster, err
		}
		controlConfigInterface = controlList.Items[0].ConfigurationParameters()
	}
	controlConfig := controlConfigInterface.(ControlConfiguration)
	bgpPort = strconv.Itoa(*controlConfig.BGPPort)
	dnsPort = strconv.Itoa(*controlConfig.DNSPort)
	xmppPort = strconv.Itoa(*controlConfig.XMPPPort)
	dnsIntrospectPort = strconv.Itoa(*controlConfig.DNSIntrospectPort)
	sort.SliceStable(controlNodes, func(i, j int) bool { return controlNodes[i] < controlNodes[j] })
	serverListXMPPCommaSeparated := strings.Join(controlNodes, ":"+xmppPort+",")
	serverListXMPPCommaSeparated = serverListXMPPCommaSeparated + ":" + xmppPort
	serverListXMPPSpaceSeparated := strings.Join(controlNodes, ":"+xmppPort+" ")
	serverListXMPPSpaceSeparated = serverListXMPPSpaceSeparated + ":" + xmppPort
	serverListDNSCommaSeparated := strings.Join(controlNodes, ":"+dnsPort+",")
	serverListDNSCommaSeparated = serverListDNSCommaSeparated + ":" + dnsPort
	serverListDNSSpaceSeparated := strings.Join(controlNodes, ":"+dnsPort+" ")
	serverListDNSSpaceSeparated = serverListDNSSpaceSeparated + ":" + dnsPort
	serverListCommanSeparatedQuoted := strings.Join(controlNodes, "','")
	serverListCommanSeparatedQuoted = "'" + serverListCommanSeparatedQuoted + "'"
	controlCluster = ControlClusterConfiguration{
		BGPPort:                         bgpPort,
		DNSPort:                         dnsPort,
		DNSIntrospectPort:               dnsIntrospectPort,
		ServerListXMPPCommaSeparated:    serverListXMPPCommaSeparated,
		ServerListXMPPSpaceSeparated:    serverListXMPPSpaceSeparated,
		ServerListDNSCommaSeparated:     serverListDNSCommaSeparated,
		ServerListDNSSpaceSeparated:     serverListDNSSpaceSeparated,
		ServerListCommanSeparatedQuoted: serverListCommanSeparatedQuoted,
	}

	return &controlCluster, nil
}

// NewZookeeperClusterConfiguration gets a struct containing various representations of Zookeeper nodes string.
func NewZookeeperClusterConfiguration(name string, namespace string, client client.Client) (*ZookeeperClusterConfiguration, error) {
	var zookeeperNodes []string
	var zookeeperCluster ZookeeperClusterConfiguration
	var port string
	zookeeperInstance := &Zookeeper{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, zookeeperInstance)
	if err != nil {
		return &zookeeperCluster, err
	}
	for _, ip := range zookeeperInstance.Status.Nodes {
		zookeeperNodes = append(zookeeperNodes, ip)

	}
	zookeeperConfig := zookeeperInstance.ConfigurationParameters()
	port = strconv.Itoa(*zookeeperConfig.ClientPort)
	sort.SliceStable(zookeeperNodes, func(i, j int) bool { return zookeeperNodes[i] < zookeeperNodes[j] })
	serverListCommaSeparated := strings.Join(zookeeperNodes, ":"+port+",")
	serverListCommaSeparated = serverListCommaSeparated + ":" + port
	serverListSpaceSeparated := strings.Join(zookeeperNodes, ":"+port+" ")
	serverListSpaceSeparated = serverListSpaceSeparated + ":" + port
	zookeeperCluster = ZookeeperClusterConfiguration{
		ClientPort:               port,
		ServerListCommaSeparated: serverListCommaSeparated,
		ServerListSpaceSeparated: serverListSpaceSeparated,
	}

	return &zookeeperCluster, nil
}

// NewRabbitmqClusterConfiguration gets a struct containing various representations of Rabbitmq nodes string.
func NewRabbitmqClusterConfiguration(name string, namespace string, myclient client.Client) (*RabbitmqClusterConfiguration, error) {
	var rabbitmqNodes []string
	var rabbitmqCluster RabbitmqClusterConfiguration
	var port string
	var sslPort string
	secret := ""
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	rabbitmqList := &RabbitmqList{}
	err := myclient.List(context.TODO(), rabbitmqList, listOps)
	if err != nil {
		return &rabbitmqCluster, err
	}
	if len(rabbitmqList.Items) > 0 {
		for _, ip := range rabbitmqList.Items[0].Status.Nodes {
			rabbitmqNodes = append(rabbitmqNodes, ip)
		}
		rabbitmqConfigInterface := rabbitmqList.Items[0].ConfigurationParameters()
		rabbitmqConfig := rabbitmqConfigInterface.(RabbitmqConfiguration)
		port = strconv.Itoa(*rabbitmqConfig.Port)
		sslPort = strconv.Itoa(*rabbitmqConfig.SSLPort)
		secret = rabbitmqList.Items[0].Status.Secret
	}
	sort.SliceStable(rabbitmqNodes, func(i, j int) bool { return rabbitmqNodes[i] < rabbitmqNodes[j] })
	serverListCommaSeparated := strings.Join(rabbitmqNodes, ":"+port+",")
	serverListCommaSeparated = serverListCommaSeparated + ":" + port
	serverListSpaceSeparated := strings.Join(rabbitmqNodes, ":"+port+" ")
	serverListSpaceSeparated = serverListSpaceSeparated + ":" + port
	serverListCommaSeparatedWithoutPort := strings.Join(rabbitmqNodes, ",")

	serverListCommaSeparatedSSL := strings.Join(rabbitmqNodes, ":"+sslPort+",")
	serverListCommaSeparatedSSL = serverListCommaSeparatedSSL + ":" + sslPort
	serverListSpaceSeparatedSSL := strings.Join(rabbitmqNodes, ":"+sslPort+" ")
	serverListSpaceSeparatedSSL = serverListSpaceSeparatedSSL + ":" + sslPort

	rabbitmqCluster = RabbitmqClusterConfiguration{
		Port:                                port,
		SSLPort:                             sslPort,
		ServerListCommaSeparated:            serverListCommaSeparated,
		ServerListSpaceSeparated:            serverListSpaceSeparated,
		ServerListCommaSeparatedSSL:         serverListCommaSeparatedSSL,
		ServerListSpaceSeparatedSSL:         serverListSpaceSeparatedSSL,
		ServerListCommaSeparatedWithoutPort: serverListCommaSeparatedWithoutPort,
		Secret:                              secret,
	}
	return &rabbitmqCluster, nil
}

// NewConfigClusterConfiguration gets a struct containing various representations of Config nodes string.
func NewConfigClusterConfiguration(name string, namespace string, myclient client.Client) (*ConfigClusterConfiguration, error) {
	var configNodes []string
	var configCluster ConfigClusterConfiguration
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	configList := &ConfigList{}
	err := myclient.List(context.TODO(), configList, listOps)
	if err != nil {
		return &configCluster, err
	}

	var apiServerPort string
	var collectorServerPort string
	var analyticsServerPort string
	var redisServerPort string
	var authMode AuthenticationMode

	if len(configList.Items) > 0 {
		for _, ip := range configList.Items[0].Status.Nodes {
			configNodes = append(configNodes, ip)
		}
		configConfigInterface := configList.Items[0].ConfigurationParameters()
		configConfig := configConfigInterface.(ConfigConfiguration)
		authMode = configConfig.AuthMode
		apiServerPort = strconv.Itoa(*configConfig.APIPort)
		analyticsServerPort = strconv.Itoa(*configConfig.AnalyticsPort)
		collectorServerPort = strconv.Itoa(*configConfig.CollectorPort)
		redisServerPort = strconv.Itoa(*configConfig.RedisPort)
	}
	sort.SliceStable(configNodes, func(i, j int) bool { return configNodes[i] < configNodes[j] })
	apiServerListQuotedCommaSeparated := strings.Join(configNodes, "','")
	firstAPIServer := configNodes[0]
	apiServerListQuotedCommaSeparated = "'" + apiServerListQuotedCommaSeparated + "'"
	analyticsServerListQuotedCommaSeparated := strings.Join(configNodes, "','")
	analyticsServerListQuotedCommaSeparated = "'" + analyticsServerListQuotedCommaSeparated + "'"
	apiServerListCommaSeparated := strings.Join(configNodes, ",")
	apiServerListSpaceSeparated := strings.Join(configNodes, ":"+apiServerPort+" ")
	apiServerListSpaceSeparated = apiServerListSpaceSeparated + ":" + apiServerPort
	analyticsServerListSpaceSeparated := strings.Join(configNodes, ":"+analyticsServerPort+" ")
	analyticsServerListSpaceSeparated = analyticsServerListSpaceSeparated + ":" + analyticsServerPort
	collectorServerListSpaceSeparated := strings.Join(configNodes, ":"+collectorServerPort+" ")
	collectorServerListSpaceSeparated = collectorServerListSpaceSeparated + ":" + collectorServerPort
	configCluster = ConfigClusterConfiguration{
		APIServerPort:                           apiServerPort,
		APIServerListQuotedCommaSeparated:       apiServerListQuotedCommaSeparated,
		APIServerListCommaSeparated:             apiServerListCommaSeparated,
		APIServerListSpaceSeparated:             apiServerListSpaceSeparated,
		AnalyticsServerPort:                     analyticsServerPort,
		AnalyticsServerListSpaceSeparated:       analyticsServerListSpaceSeparated,
		AnalyticsServerListQuotedCommaSeparated: analyticsServerListQuotedCommaSeparated,
		CollectorPort:                           collectorServerPort,
		CollectorServerListSpaceSeparated:       collectorServerListSpaceSeparated,
		FirstAPIServer:                          firstAPIServer,
		RedisPort:                               redisServerPort,
		AuthMode:                                authMode,
	}
	return &configCluster, nil
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

// ConfigClusterConfiguration defines all configuration knobs used to write the config file.
type ConfigClusterConfiguration struct {
	APIServerPort                           string
	APIServerListSpaceSeparated             string
	APIServerListQuotedCommaSeparated       string
	APIServerListCommaSeparated             string
	AnalyticsServerPort                     string
	AnalyticsServerListSpaceSeparated       string
	AnalyticsServerListQuotedCommaSeparated string
	CollectorServerListSpaceSeparated       string
	CollectorPort                           string
	FirstAPIServer                          string
	RedisPort                               string
	AuthMode                                AuthenticationMode
}

// ControlClusterConfiguration defines all configuration knobs used to write the config file.
type ControlClusterConfiguration struct {
	BGPPort                         string
	DNSPort                         string
	DNSIntrospectPort               string
	ServerListXMPPCommaSeparated    string
	ServerListXMPPSpaceSeparated    string
	ServerListDNSCommaSeparated     string
	ServerListDNSSpaceSeparated     string
	ServerListCommanSeparatedQuoted string
}

// ZookeeperClusterConfiguration defines all configuration knobs used to write the config file.
type ZookeeperClusterConfiguration struct {
	ClientPort               string
	ServerPort               string
	ElectionPort             string
	ServerListCommaSeparated string
	ServerListSpaceSeparated string
}

// RabbitmqClusterConfiguration defines all configuration knobs used to write the config file.
type RabbitmqClusterConfiguration struct {
	Port                                string
	SSLPort                             string
	ServerListCommaSeparated            string
	ServerListSpaceSeparated            string
	ServerListCommaSeparatedSSL         string
	ServerListSpaceSeparatedSSL         string
	ServerListCommaSeparatedWithoutPort string
	Secret                              string
}

// CassandraClusterConfiguration defines all configuration knobs used to write the config file.
type CassandraClusterConfiguration struct {
	Port                            string
	CQLPort                         string
	JMXPort                         string
	Endpoint                        string
	ServerListCQLSpaceSeparated     string
	ServerListJMXCommaSeparated     string
	ServerListJMXSpaceSeparated     string
	ServerListCommanSeparatedQuoted string
}

// VrouterClusterConfiguration defines all configuration knobs used to write the config file.
type VrouterClusterConfiguration struct {
	PhysicalInterface string
	Gateway           string
	MetaDataSecret    string
}
