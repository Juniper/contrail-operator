package provisionmanager

import (
	"context"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const RequiredAnnotationsKey = "managed_by"

var log = logf.Log.WithName("controller_provisionmanager")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
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
			list := &v1alpha1.ProvisionManagerList{}
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
			list := &v1alpha1.ProvisionManagerList{}
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
			list := &v1alpha1.ProvisionManagerList{}
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

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ProvisionManager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProvisionManager{Client: mgr.GetClient(), Scheme: mgr.GetScheme(), Manager: mgr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("provisionmanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ProvisionManager
	err = c.Watch(&source.Kind{Type: &v1alpha1.ProvisionManager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to PODs
	serviceMap := map[string]string{"contrail_manager": "provisionmanager"}
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

	srcConfig := &source.Kind{Type: &v1alpha1.Config{}}
	configHandler := resourceHandler(mgr.GetClient())
	predConfigActiveChange := utils.ConfigActiveChange()
	if err = c.Watch(srcConfig, configHandler, predConfigActiveChange); err != nil {
		return err
	}

	srcControl := &source.Kind{Type: &v1alpha1.Control{}}
	controlHandler := resourceHandler(mgr.GetClient())
	predControlActiveChange := utils.ControlActiveChange()
	if err = c.Watch(srcControl, controlHandler, predControlActiveChange); err != nil {
		return err
	}

	srcVrouter := &source.Kind{Type: &v1alpha1.Vrouter{}}
	vrouterHandler := resourceHandler(mgr.GetClient())
	predVrouterActiveChange := utils.VrouterActiveChange()
	if err = c.Watch(srcVrouter, vrouterHandler, predVrouterActiveChange); err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileProvisionManager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileProvisionManager{}

// ReconcileProvisionManager reconciles a ProvisionManager object
type ReconcileProvisionManager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client  client.Client
	Scheme  *runtime.Scheme
	Manager manager.Manager
}

func (r *ReconcileProvisionManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling  ProvisionManager")
	instanceType := "provisionmanager"
	instance := &v1alpha1.ProvisionManager{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	configInstance := &v1alpha1.Config{}
	configActive := configInstance.IsActive(instance.Labels["contrail_cluster"], request.Namespace, r.Client)
	if !configActive {
		return reconcile.Result{}, nil
	}

	configMapConfigNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-confignodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapControlNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-controlnodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapVrouterNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-vrouternodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapAnalyticsNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-analyticsnodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapDatabaseNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-databasenodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapAPIServer, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-apiserver", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapKeystoneAuthConf, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-keystoneauth", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapGlobalVrouterConf, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-globalvrouter", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	secretCertificates, err := instance.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	statefulSet.Spec.Template.Annotations = map[string]string{RequiredAnnotationsKey: request.Name + "-" + instanceType}

	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{
		configMapConfigNodes.Name:          request.Name + "-" + instanceType + "-confignodes-volume",
		configMapControlNodes.Name:         request.Name + "-" + instanceType + "-controlnodes-volume",
		configMapVrouterNodes.Name:         request.Name + "-" + instanceType + "-vrouternodes-volume",
		configMapAnalyticsNodes.Name:       request.Name + "-" + instanceType + "-analyticsnodes-volume",
		configMapDatabaseNodes.Name:        request.Name + "-" + instanceType + "-databasenodes-volume",
		configMapAPIServer.Name:            request.Name + "-" + instanceType + "-apiserver-volume",
		configMapKeystoneAuthConf.Name:     request.Name + "-" + instanceType + "-keystoneauth-volume",
		configMapGlobalVrouterConf.Name:    request.Name + "-" + instanceType + "-globalvrouter-volume",
		certificates.SignerCAConfigMapName: csrSignerCaVolumeName,
	})
	instance.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {
		if container.Name == "provisioner" {
			command := []string{"sh", "-c",
				`/app/contrail-provisioner/contrail-provisioner-image.binary \
					-controlNodes /etc/provision/control/controlnodes.yaml \
					-configNodes /etc/provision/config/confignodes.yaml \
					-analyticsNodes /etc/provision/analytics/analyticsnodes.yaml \
					-vrouterNodes /etc/provision/vrouter/vrouternodes.yaml \
					-databaseNodes /etc/provision/database/databasenodes.yaml \
					-apiserver /etc/provision/apiserver/apiserver-${POD_IP}.yaml \
					-keystoneAuthConf /etc/provision/keystone/keystone-auth-${POD_IP}.yaml \
					-globalVrouterConf /etc/provision/globalvrouter/globalvrouter.json \
					-mode watch`,
			}
			instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
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
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-confignodes-volume",
				MountPath: "/etc/provision/config",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-analyticsnodes-volume",
				MountPath: "/etc/provision/analytics",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-controlnodes-volume",
				MountPath: "/etc/provision/control",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-vrouternodes-volume",
				MountPath: "/etc/provision/vrouter",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-databasenodes-volume",
				MountPath: "/etc/provision/database",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-apiserver-volume",
				MountPath: "/etc/provision/apiserver",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-keystoneauth-volume",
				MountPath: "/etc/provision/keystone",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-globalvrouter-volume",
				MountPath: "/etc/provision/globalvrouter",
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

	// Configure InitContainers
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
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instanceContainer.Image
		if instanceContainer.Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instanceContainer.Command
		}
	}

	if err = instance.CreateSTS(statefulSet, instanceType, request, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	strategy := "deleteFirst"
	if err = instance.UpdateSTS(statefulSet, instanceType, request, r.Client, strategy); err != nil {
		return reconcile.Result{}, err
	}

	podIPList, podIPMap, err := instance.PodIPListAndIPMapFromInstance(request, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPMap) > 0 {
		if err = instance.InstanceConfiguration(request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		if err := r.ensureCertificatesExist(instance, podIPList, instanceType); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
	}
	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err = instance.SetInstanceActive(r.Client, instance.Status.Active, statefulSet, request); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileProvisionManager) ensureCertificatesExist(provision *v1alpha1.ProvisionManager, pods *corev1.PodList, instanceType string) error {
	subjects := provision.PodsCertSubjects(pods)
	crt := certificates.NewCertificate(r.Client, r.Scheme, provision, subjects, instanceType)
	return crt.EnsureExistsAndIsSigned()
}
