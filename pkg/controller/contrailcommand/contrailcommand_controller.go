package contrailcommand

import (
	"context"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrail "atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
)

var log = logf.Log.WithName("controller_contrailcommand")

// Add creates a new ContrailCommand Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileContrailCommand{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("contrailcommand-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ContrailCommand
	err = c.Watch(&source.Kind{Type: &contrail.ContrailCommand{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner ContrailCommand
	err = c.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.ContrailCommand{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Postgres and requeue the owner ContrailCommand
	err = c.Watch(&source.Kind{Type: &contrail.Postgres{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.ContrailCommand{},
	})

	return err
}

// blank assignment to verify that ReconcileContrailCommand implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileContrailCommand{}

// ReconcileContrailCommand reconciles a ContrailCommand object
type ReconcileContrailCommand struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ContrailCommand object and makes changes based on the state read
func (r *ReconcileContrailCommand) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ContrailCommand")
	instanceType := "contrailcommand"
	// Fetch the ContrailCommand command
	command := &contrail.ContrailCommand{}
	if err := r.Client.Get(context.Background(), request.NamespacedName, command); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	configMap, err := contrail.CreateConfigMap(request.Name+"-contrailcommand-configmap", r.Client, r.Scheme, request, "contrailcommand", command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := command.InstanceConfiguration(request, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	postgres, err := r.ensurePostgres(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !postgres.Status.Active {
		return reconcile.Result{}, nil
	}

	configVolumeName := request.Name + "-" + instanceType + "-volume"
	deployment := newDeployment(
		request.Name+"-"+instanceType+"-deployment",
		request.Namespace,
		configVolumeName,
	)

	contrail.AddVolumesToIntendedDeployments(deployment, map[string]string{configMap.Name: configVolumeName})

	if _, err := controllerutil.CreateOrUpdate(context.Background(), r.Client, deployment, func(existing runtime.Object) error {
		deployment := existing.(*apps.Deployment)
		deployment, err := command.PrepareIntendedDeployment(deployment,
			&command.Spec.CommonConfiguration, request, r.Scheme)
		return err
	}); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(command, deployment, request)
}

func (r *ReconcileContrailCommand) ensurePostgres(command *contrail.ContrailCommand) (*contrail.Postgres, error) {
	postgres := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Name:      command.Name + "-db",
			Namespace: command.Namespace,
		},
	}

	// Set Command instance as the owner and controller
	if err := controllerutil.SetControllerReference(command, postgres, r.Scheme); err != nil {
		return nil, err
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.Client, postgres, func(existing runtime.Object) error {
		postgres = existing.(*contrail.Postgres)
		return nil
	})

	return postgres, err
}

func newDeployment(name, namespace, configVolumeName string) *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{},
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					Containers: []core.Container{{
						Name:            "command",
						ImagePullPolicy: core.PullAlways,
						Image:           "localhost:5000/contrail-command",
						ReadinessProbe: &core.Probe{
							Handler: core.Handler{
								HTTPGet: &core.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 9091},
								},
							},
						},
						VolumeMounts: []core.VolumeMount{{
							Name:      configVolumeName,
							MountPath: "/etc/contrail",
						}},
					}},
					DNSPolicy: core.DNSClusterFirst,
				},
			},
		},
	}
}

func (r *ReconcileContrailCommand) updateStatus(
	command *contrail.ContrailCommand,
	deployment *apps.Deployment,
	request reconcile.Request,
) error {
	err := r.Client.Get(context.Background(), types.NamespacedName{Name: deployment.Name, Namespace: request.Namespace}, deployment)
	if err != nil {
		return err
	}
	command.Status.Active = false
	intendentReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		intendentReplicas = *deployment.Spec.Replicas
	}

	if deployment.Status.ReadyReplicas == intendentReplicas {
		command.Status.Active = true
	}

	return r.Client.Status().Update(context.Background(), command)
}
