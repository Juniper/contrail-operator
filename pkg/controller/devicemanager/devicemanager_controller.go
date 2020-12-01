package devicemanager

import (
	"context"
	"reflect"

	"github.com/Juniper/contrail-operator/pkg/certificates"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/Juniper/contrail-operator/pkg/controller/utils"

	"github.com/Juniper/contrail-operator/pkg/k8s"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_devicemanager")

// Add creates a new Devicemanager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDevicemanager{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Manager:    mgr,
		Kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme()),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("devicemanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Devicemanager.
	if err = c.Watch(&source.Kind{Type: &v1alpha1.Devicemanager{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Devicemanager
	if err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Devicemanager{},
	}); err != nil {
		return err
	}

	srcSTS := &source.Kind{Type: &appsv1.StatefulSet{}}
	stsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Devicemanager{},
	}
	if err = c.Watch(srcSTS, stsHandler); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDevicemanager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDevicemanager{}

// ReconcileDevicemanager reconciles a Devicemanager object
type ReconcileDevicemanager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client     client.Client
	Scheme     *runtime.Scheme
	Manager    manager.Manager
	Kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Devicemanager object and makes changes based on the state read
// and what is in the Devicemanager.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDevicemanager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Devicemanager")
	instanceType := "devicemanager"

	// Fetch the Devicemanager instance
	devicemanager := &contrailv1alpha1.Devicemanager{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, devicemanager)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !devicemanager.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Checking dependencies are up")

	cassandraInstance := &v1alpha1.Cassandra{}
	zookeeperInstance := &v1alpha1.Zookeeper{}
	rabbitmqInstance := &v1alpha1.Rabbitmq{}
	configInstance := &v1alpha1.Config{}

	cassandraActive := cassandraInstance.IsActive(devicemanager.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	zookeeperActive := zookeeperInstance.IsActive(devicemanager.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, r.Client)
	configActive := configInstance.IsActive(devicemanager.Labels["contrail_cluster"],
		request.Namespace, r.Client)
	rabbitmqActive := rabbitmqInstance.IsActive(devicemanager.Labels["contrail_cluster"],
		request.Namespace, r.Client)

	if !cassandraActive || !rabbitmqActive || !zookeeperActive || !configActive {
		reqLogger.Info("Skipping: Waiting for dependencies")
		return reconcile.Result{}, nil
	}
	// Define a new Pod object
	//pod := newPodForCR(instance)

	// TODO figure out if we need it. We don't use  SetControllerReference in config
	// Set Devicemanager instance as the owner and controller
	//if err := controllerutil.SetControllerReference(devicemanager, pod, r.Scheme); err != nil {
	//	return reconcile.Result{}, err
	//}

	// TODO check if need this
	// Check if this Pod already exists
	//found := &corev1.Pod{}
	//err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	//	err = r.Client.Create(context.TODO(), pod)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//
	//	// Pod created successfully - don't requeue
	//	return reconcile.Result{}, nil
	//} else if err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//// Pod already exists - don't requeue
	//reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	//return reconcile.Result{}, nil

	servicePortsMap := map[int32]string{
		int32(v1alpha1.ConfigDeviceManagerIntrospectPort): "introspect",
	}
	devicemanagerService := r.Kubernetes.Service(request.Name+"-"+instanceType, corev1.ServiceTypeClusterIP,
		servicePortsMap, instanceType, devicemanager)

	if err := devicemanagerService.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	currentConfigMap, currentConfigExists := devicemanager.CurrentConfigMapExists(
		request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)

	reqLogger.Info("Creating config map")

	configMap, err := devicemanager.CreateConfigMap(request.Name+"-"+instanceType+"-configmap",
		r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Creating secrets")

	secretCertificates, err := devicemanager.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	trueVal := true
	statefulSet.Spec.Template.Spec.ShareProcessNamespace = &trueVal
	if err = devicemanager.PrepareSTS(statefulSet, &devicemanager.Spec.CommonConfiguration, request, r.Scheme,
		r.Client); err != nil {
		return reconcile.Result{}, err
	}

	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	devicemanager.AddVolumesToIntendedSTS(statefulSet, map[string]string{
		configMap.Name:                     request.Name + "-" + instanceType + "-volume",
		certificates.SignerCAConfigMapName: csrSignerCaVolumeName,
	})
	devicemanager.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	reqLogger.Info("Creating account statusmonitor-devicemanager")
	if err = v1alpha1.CreateAccount("statusmonitor-devicemanager", request.Namespace, r.Client,
		r.Scheme, devicemanager); err != nil {
		return reconcile.Result{}, err
	}

	statefulSet.Spec.Template.Spec.ServiceAccountName = "serviceaccount-statusmonitor-devicemanager"
	statefulSet.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      instanceType,
						Operator: "In",
						Values:   []string{request.Name},
					}},
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}

	reqLogger.Info("Adding containers")
	for idx, container := range statefulSet.Spec.Template.Spec.Containers {
		reqLogger.Info("Adding container", "name", container.Name)
		switch container.Name {
		case "devicemanager":
			deviceManagerCommand := `/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/contrailconfigmaps/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini;
/usr/bin/rm -f /etc/contrail/contrail-keystone-auth.conf; ln -s /etc/contrailconfigmaps/contrail-keystone-auth.conf /etc/contrail/contrail-keystone-auth.conf;
/usr/bin/rm -f /etc/contrail/contrail-fabric-ansible.conf; ln -s /etc/contrailconfigmaps/contrail-fabric-ansible.conf.${POD_IP} /etc/contrail/contrail-fabric-ansible.conf;
/usr/bin/python /usr/bin/contrail-device-manager --conf_file /etc/contrailconfigmaps/devicemanager.${POD_IP} --conf_file /etc/contrail/contrail-keystone-auth.conf
`
			command := []string{"bash", "-c", deviceManagerCommand}
			instanceContainer := utils.GetContainerFromList(container.Name, devicemanager.Spec.ServiceConfiguration.Containers)
			if instanceContainer.Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instanceContainer.Command
			}
			(&statefulSet.Spec.Template.Spec.Containers[idx]).SecurityContext = &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{"SYS_PTRACE"},
				},
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/contrailconfigmaps",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: certificates.SignerCAMountPath,
				},
				corev1.VolumeMount{
					Name:      "tftp",
					MountPath: "/var/lib/tftp",
				},
				corev1.VolumeMount{
					Name:      "dnsmasq",
					MountPath: "/var/lib/dnsmasq",
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instanceContainer.Image

		case "dnsmasq":
			container := &statefulSet.Spec.Template.Spec.Containers[idx]
			container.Command = []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini;ln -s /etc/contrailconfigmaps/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini;" +
					"dnsmasq -k -p0 --conf-file=/etc/contrailconfigmaps/dnsmasq.${POD_IP}"}
			instanceContainer := utils.GetContainerFromList(container.Name, devicemanager.Spec.ServiceConfiguration.Containers)
			if instanceContainer.Command != nil {
				container.Command = instanceContainer.Command
			}
			container.SecurityContext = &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{"NET_ADMIN", "NET_RAW"},
				},
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/contrailconfigmaps",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: certificates.SignerCAMountPath,
				},
				corev1.VolumeMount{
					Name:      "tftp",
					MountPath: "/etc/tftp",
				},
				corev1.VolumeMount{
					Name:      "dnsmasq",
					MountPath: "/var/lib/dnsmasq",
				},
			)
			// DNSMasq container requires those variables to be set
			// TODO: Pass keystone credentials
			container.Env = append(container.Env, []corev1.EnvVar{
				{Name: "KEYSTONE_AUTH_ADMIN_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: devicemanager.Spec.ServiceConfiguration.KeystoneSecretName,
							},
							Key: "password",
						},
					},
				},
				{Name: "KEYSTONE_AUTH_ADMIN_USER", Value: "admin"},
				{Name: "KEYSTONE_AUTH_ADMIN_TENANT", Value: "admin"},
			}...)
			container.VolumeMounts = volumeMountList
			container.Image = instanceContainer.Image

		case "statusmonitor":
			instanceContainer := utils.GetContainerFromList(container.Name, devicemanager.Spec.ServiceConfiguration.Containers)
			command := []string{"sh", "-c",
				"/app/statusmonitor/contrail-statusmonitor-image.binary -config /etc/contrailconfigmaps/monitorconfig.${POD_IP}.yaml"}
			if instanceContainer.Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instanceContainer.Command
			}

			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/contrailconfigmaps",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      csrSignerCaVolumeName,
				MountPath: certificates.SignerCAMountPath,
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instanceContainer.Image
		}
	}

	reqLogger.Info("Adding init containers")
	// Configure InitContainers
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		instanceContainer := utils.GetContainerFromList(container.Name, devicemanager.Spec.ServiceConfiguration.Containers)
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instanceContainer.Image
		if instanceContainer.Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instanceContainer.Command
		}
	}

	configChanged := false
	if devicemanager.Status.ConfigChanged != nil {
		configChanged = *devicemanager.Status.ConfigChanged
	}

	if err = devicemanager.CreateSTS(statefulSet, instanceType, request, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	podIPList, podIPMap, err := devicemanager.PodIPListAndIPMapFromInstance(request, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPMap) > 0 {
		if err = devicemanager.InstanceConfiguration(request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err := r.ensureCertificatesExist(devicemanager, podIPList, instanceType); err != nil {
			return reconcile.Result{}, err
		}

		if err = devicemanager.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = devicemanager.WaitForPeerPods(request, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = devicemanager.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err = devicemanager.SetEndpointInStatus(r.Client, devicemanagerService.ClusterIP()); err != nil {
		return reconcile.Result{}, err
	}
	if currentConfigExists {
		newConfigMap := &corev1.ConfigMap{}
		_ = r.Client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-configmap", Namespace: request.Namespace}, newConfigMap)
		if !reflect.DeepEqual(currentConfigMap.Data, newConfigMap.Data) {
			configChanged = true
		} else {
			configChanged = false
		}
		devicemanager.Status.ConfigChanged = &configChanged
	}

	if devicemanager.Status.Active == nil {
		active := false
		devicemanager.Status.Active = &active
	}
	if err = devicemanager.SetInstanceActive(r.Client, devicemanager.Status.Active, statefulSet, request); err != nil {
		return reconcile.Result{}, err
	}
	if devicemanager.Status.ConfigChanged != nil {
		if *devicemanager.Status.ConfigChanged {
			return reconcile.Result{Requeue: true}, nil
		}
	}
	return reconcile.Result{}, nil

}

func (r *ReconcileDevicemanager) ensureCertificatesExist(devicemaanger *v1alpha1.Devicemanager, pods *corev1.PodList, instanceType string) error {
	subjects := devicemaanger.PodsCertSubjects(pods)
	crt := certificates.NewCertificate(r.Client, r.Scheme, devicemaanger, subjects, instanceType)
	return crt.EnsureExistsAndIsSigned()
}
