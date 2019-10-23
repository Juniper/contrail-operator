package contrailcommand

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrailv1alpha1 "atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
)

var log = logf.Log.WithName("controller_contrailcommand")

// Add creates a new ContrailCommand Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileContrailCommand{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("contrailcommand-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ContrailCommand
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.ContrailCommand{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner ContrailCommand
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.ContrailCommand{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileContrailCommand implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileContrailCommand{}

// ReconcileContrailCommand reconciles a ContrailCommand object
type ReconcileContrailCommand struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ContrailCommand object and makes changes based on the state read
func (r *ReconcileContrailCommand) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ContrailCommand")
	instanceType := "contrailcommand"
	// Fetch the ContrailCommand instance
	instance := &contrailv1alpha1.ContrailCommand{}
	err := r.client.Get(context.Background(), request.NamespacedName, instance)
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

	configMap, err := contrailv1alpha1.CreateConfigMap(request.Name+"-contrailcommand-configmap", r.client, r.scheme, request, "contrailcommand", instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = instance.InstanceConfiguration(request, r.client); err != nil {
		return reconcile.Result{}, err
	}

	deployment, err := instance.PrepareIntendedDeployment(getDeployment(), &instance.Spec.CommonConfiguration, request, r.scheme)
	if err != nil {
		return reconcile.Result{}, err
	}

	configVolumeName := request.Name + "-" + instanceType + "-volume"
	instance.AddVolumesToIntendedDeployments(deployment,
		map[string]string{configMap.Name: configVolumeName})

	volumeMountList := []corev1.VolumeMount{}
	volumeMount := corev1.VolumeMount{
		Name:      configVolumeName,
		MountPath: "/etc/contrail",
	}
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(volumeMountList, volumeMount)

	if _, err := controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func(existing runtime.Object) error {
		// TODO Handle update
		return nil
	}); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func getDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "command",
						ImagePullPolicy: corev1.PullAlways,
						Image:           "localhost:5000/contrail-command",
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 9091},
								},
							},
						},
					}},
					DNSPolicy: corev1.DNSClusterFirst,
				},
			},
		},
	}
}
