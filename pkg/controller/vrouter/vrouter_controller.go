package vrouter

import (
	"context"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_vrouter")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.VrouterList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.MetaNew.GetNamespace()}
			list := &v1alpha1.VrouterList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.MetaNew.GetNamespace(),
					}})
				}
			}
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.VrouterList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},

		GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.VrouterList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
	}
	return appHandler
}

// Add creates a new Vrouter Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVrouter{Client: mgr.GetClient(), Scheme: mgr.GetScheme(), Manager: mgr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller.
	c, err := controller.New("vrouter-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Vrouter.
	if err = c.Watch(&source.Kind{Type: &v1alpha1.Vrouter{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to PODs.
	serviceMap := map[string]string{"contrail_manager": "vrouter"}
	srcPod := &source.Kind{Type: &corev1.Pod{}}
	podHandler := resourceHandler(mgr.GetClient())
	predInitStatus := utils.PodInitStatusChange(serviceMap)
	predPodIPChange := utils.PodIPChange(serviceMap)
	predInitRunning := utils.PodInitRunning(serviceMap)

	if err = c.Watch(srcPod, podHandler, predPodIPChange); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitStatus); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitRunning); err != nil {
		return err
	}

	srcManager := &source.Kind{Type: &v1alpha1.Manager{}}
	managerHandler := resourceHandler(mgr.GetClient())
	predManagerSizeChange := utils.ManagerSizeChange(utils.VrouterGroupKind())
	if err = c.Watch(srcManager, managerHandler, predManagerSizeChange); err != nil {
		return err
	}

	srcConfig := &source.Kind{Type: &v1alpha1.Config{}}
	configHandler := resourceHandler(mgr.GetClient())
	predConfigSizeChange := utils.ConfigActiveChange()
	if err = c.Watch(srcConfig, configHandler, predConfigSizeChange); err != nil {
		return err
	}

	srcControl := &source.Kind{Type: &v1alpha1.Control{}}
	controlHandler := resourceHandler(mgr.GetClient())
	predControlSizeChange := utils.ControlActiveChange()
	if err = c.Watch(srcControl, controlHandler, predControlSizeChange); err != nil {
		return err
	}

	srcDS := &source.Kind{Type: &appsv1.DaemonSet{}}
	dsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Vrouter{},
	}
	dsPred := utils.DSStatusChange(utils.VrouterGroupKind())
	if err = c.Watch(srcDS, dsHandler, dsPred); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileVrouter implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileVrouter{}

// ReconcileVrouter reconciles a Vrouter object.
type ReconcileVrouter struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	Client  client.Client
	Scheme  *runtime.Scheme
	Manager manager.Manager
}

// Reconcile reads that state of the cluster for a Vrouter object and makes changes based on the state read
// and what is in the Vrouter.Spec.
func (r *ReconcileVrouter) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Vrouter")
	instanceType := "vrouter"
	instance := &v1alpha1.Vrouter{}
	controlInstance := v1alpha1.Control{}
	configInstance := v1alpha1.Config{}

	if err := r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	configActive := configInstance.IsActive(instance.Labels["contrail_cluster"],
		request.Namespace, r.Client)

	controlActive := controlInstance.IsActive(instance.Spec.ServiceConfiguration.ControlInstance,
		request.Namespace, r.Client)

	if !configActive || !controlActive {
		return reconcile.Result{}, nil
	}

	managerInstance, err := instance.OwnedByManager(r.Client, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if managerInstance != nil {
		if managerInstance.Spec.Services.Vrouters != nil {
			for _, vrouterManagerInstance := range managerInstance.Spec.Services.Vrouters {
				if vrouterManagerInstance.Name == request.Name {
					instance.Spec.CommonConfiguration = utils.MergeCommonConfiguration(
						managerInstance.Spec.CommonConfiguration,
						vrouterManagerInstance.Spec.CommonConfiguration)
					err = r.Client.Update(context.TODO(), instance)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
		}
	}

	configMap, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-1", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	secretCertificates, err := instance.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	daemonSet := GetDaemonset()
	if err = instance.PrepareDaemonSet(daemonSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	instance.AddVolumesToIntendedDS(daemonSet, map[string]string{
		configMap.Name:                          request.Name + "-" + instanceType + "-volume",
		cacertificates.CsrSignerCAConfigMapName: csrSignerCaVolumeName,
	})
	instance.AddSecretVolumesToIntendedDS(daemonSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	var serviceAccountName string
	if instance.Spec.ServiceConfiguration.ServiceAccount != "" {
		serviceAccountName = instance.Spec.ServiceConfiguration.ServiceAccount
	} else {
		serviceAccountName = "contrail-service-account-cni"
	}

	var clusterRoleName string
	if instance.Spec.ServiceConfiguration.ClusterRole != "" {
		clusterRoleName = instance.Spec.ServiceConfiguration.ClusterRole
	} else {
		clusterRoleName = "contrail-cluster-role-cni"
	}

	var clusterRoleBindingName string
	if instance.Spec.ServiceConfiguration.ClusterRoleBinding != "" {
		clusterRoleBindingName = instance.Spec.ServiceConfiguration.ClusterRoleBinding
	} else {
		clusterRoleBindingName = "contrail-cluster-role-binding-cni"
	}

	existingServiceAccount := &corev1.ServiceAccount{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: serviceAccountName, Namespace: instance.Namespace}, existingServiceAccount)
	if err != nil && errors.IsNotFound(err) {
		serviceAccount := &corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceAccountName,
				Namespace: instance.Namespace,
			},
		}
		err = controllerutil.SetControllerReference(instance, serviceAccount, r.Scheme)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = r.Client.Create(context.TODO(), serviceAccount); err != nil && !errors.IsAlreadyExists(err) {
			return reconcile.Result{}, err
		}
	}

	existingClusterRole := &rbacv1.ClusterRole{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName}, existingClusterRole)
	if err != nil && errors.IsNotFound(err) {
		clusterRole := &rbacv1.ClusterRole{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac/v1",
				Kind:       "ClusterRole",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterRoleName,
				Namespace: instance.Namespace,
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
		err = controllerutil.SetControllerReference(instance, clusterRole, r.Scheme)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = r.Client.Create(context.TODO(), clusterRole); err != nil {
			return reconcile.Result{}, err
		}
	}

	existingClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleBindingName}, existingClusterRoleBinding)
	if err != nil && errors.IsNotFound(err) {
		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac/v1",
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterRoleBindingName,
				Namespace: instance.Namespace,
			},
			Subjects: []rbacv1.Subject{{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: instance.Namespace,
			}},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     clusterRoleName,
			},
		}
		err = controllerutil.SetControllerReference(instance, clusterRoleBinding, r.Scheme)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = r.Client.Create(context.TODO(), clusterRoleBinding); err != nil {
			return reconcile.Result{}, err
		}
	}

	daemonSet.Spec.Template.Spec.ServiceAccountName = serviceAccountName
	nodemgr := true

	if _, ok := instance.Spec.ServiceConfiguration.Containers["nodemanager"]; !ok {
		nodemgr = false
	}
	for idx, container := range daemonSet.Spec.Template.Spec.Containers {
		if container.Name == "vrouteragent" {
			//command := []string{"bash", "-c",
			//	"/usr/bin/contrail-vrouter-agent --config_file /etc/mycontrail/vrouter.${POD_IP}"}
			command := []string{"bash", "-c",
				"/entrypoint.sh /usr/bin/contrail-vrouter-agent --config_file /etc/mycontrail/vrouter.${POD_IP}"}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      csrSignerCaVolumeName,
				MountPath: cacertificates.CsrSignerCAMountPath,
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&daemonSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
			(&daemonSet.Spec.Template.Spec.Containers[idx]).EnvFrom = []corev1.EnvFromSource{{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: request.Name + "-" + instanceType + "-configmap-1",
					},
				},
			}}
		}
		if container.Name == "vroutercni" {
			command := []string{"sh", "-c",
				"cp loopback /host/opt_cni_bin/ && cp /contrailcni_client /host/opt_cni_bin/contrail-k8s-cni && mkdir /host/etc_cni/net.d ||true && mkdir /var/run/contrail || true && mkdir -p /var/lib/contrail/ports/vm || true && cp /etc/mycontrail/10-contrail.conf /host/etc_cni/net.d/10-contrail.conf && /contrailcni_server"}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}

			volumeMountList := []corev1.VolumeMount{}
			if len((&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&daemonSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
			(&daemonSet.Spec.Template.Spec.Containers[idx]).EnvFrom = []corev1.EnvFromSource{{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: request.Name + "-" + instanceType + "-configmap-1",
					},
				},
			}}
		}
		if container.Name == "nodemanager" {
			if nodemgr {
				command := []string{"bash", "-c",
					"bash /etc/mycontrail/provision.sh.${POD_IP} add; /usr/bin/python /usr/bin/contrail-nodemgr --nodetype=contrail-vrouter"}
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
					(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = command
				} else {
					(&daemonSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
				}

				volumeMountList := []corev1.VolumeMount{}
				if len((&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
					volumeMountList = (&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
				}
				volumeMount := corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				}
				volumeMountList = append(volumeMountList, volumeMount)
				volumeMount = corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				}
				volumeMountList = append(volumeMountList, volumeMount)
				volumeMount = corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				}
				volumeMountList = append(volumeMountList, volumeMount)
				(&daemonSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
				(&daemonSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
			}
		}
	}

	if !nodemgr {
		for idx, container := range daemonSet.Spec.Template.Spec.Containers {
			if container.Name == "nodemanager" {
				daemonSet.Spec.Template.Spec.Containers = utils.RemoveIndex(daemonSet.Spec.Template.Spec.Containers, idx)
			}
		}
	}

	ubuntu := v1alpha1.UBUNTU
	for idx, container := range daemonSet.Spec.Template.Spec.InitContainers {
		(&daemonSet.Spec.Template.Spec.InitContainers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		if instance.Spec.ServiceConfiguration.Containers[container.Name].Command != nil {
			(&daemonSet.Spec.Template.Spec.InitContainers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
		}
		if container.Name == "vrouterkernelinit" {
			(&daemonSet.Spec.Template.Spec.InitContainers[idx]).EnvFrom = []corev1.EnvFromSource{{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: request.Name + "-" + instanceType + "-configmap-1",
					},
				},
			}}
			if instance.Spec.ServiceConfiguration.Distribution != nil || instance.Spec.ServiceConfiguration.Distribution == &ubuntu {
				(&daemonSet.Spec.Template.Spec.InitContainers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
			}
		}

	}
	if err = instance.CreateDS(daemonSet, &instance.Spec.CommonConfiguration, instanceType, request,
		r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	if err = instance.UpdateDS(daemonSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	getPhysicalInterface := false
	if instance.Spec.ServiceConfiguration.PhysicalInterface == "" {
		getPhysicalInterface = true
	}
	getGateway := false
	if instance.Spec.ServiceConfiguration.Gateway == "" {
		getGateway = true
	}
	podIPList, podIPMap, err := instance.PodIPListAndIPMapFromInstance(instanceType, request, r.Client, getPhysicalInterface, true, true, getGateway)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPMap) > 0 {
		if err = instance.InstanceConfiguration(request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		hostNetwork := true
		if instance.Spec.CommonConfiguration.HostNetwork != nil {
			hostNetwork = *instance.Spec.CommonConfiguration.HostNetwork
		}
		if err = certificates.CreateAndSignCsr(r.Client, request, r.Scheme, instance, r.Manager.GetConfig(), podIPList, hostNetwork); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
		}
	}

	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err = instance.SetInstanceActive(r.Client, instance.Status.Active, daemonSet, request, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
