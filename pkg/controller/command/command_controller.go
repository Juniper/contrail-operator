package command

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

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_command")

// Add creates a new Command Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	kubernetes := k8s.New(mgr.GetClient(), mgr.GetScheme())
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), kubernetes)
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	log.Info("Test adding Command")
	// Create a new controller
	c, err := controller.New("command-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource Command
	err = c.Watch(&source.Kind{Type: &contrail.Command{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Deployment and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Command{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Postgres and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Postgres{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	log.Info("Test Watch Command3", "Err", err)
	return err
}

// blank assignment to verify that ReconcileCommand implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCommand{}

// ReconcileCommand reconciles a Command object
type ReconcileCommand struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// NewReconciler is used to create command reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes) *ReconcileCommand {
	return &ReconcileCommand{client: client, scheme: scheme, kubernetes: kubernetes}
}

// Reconcile reads that state of the cluster for a Command object and makes changes based on the state read
func (r *ReconcileCommand) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Info("Test reconcile Command")
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Command")
	instanceType := "command"
	// Fetch the Command command
	command := &contrail.Command{}
	if err := r.client.Get(context.Background(), request.NamespacedName, command); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !command.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	commandConfigName := command.Name + "-command-configmap"
	if err := r.configMap(commandConfigName, "command", command).ensureCommandConfigExist(); err != nil {
		return reconcile.Result{}, err
	}

	psql, err := r.getPostgres(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubernetes.Owner(command).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}

	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}

	configVolumeName := request.Name + "-" + instanceType + "-volume"
	deployment := newDeployment(
		request.Name+"-"+instanceType+"-deployment",
		request.Namespace,
		configVolumeName,
		command.Spec.ServiceConfiguration.Containers,
	)

	contrail.AddVolumesToIntendedDeployments(deployment, map[string]string{commandConfigName: configVolumeName})

	if _, err := controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		_, err := command.PrepareIntendedDeployment(deployment,
			&command.Spec.CommonConfiguration, request, r.scheme)
		return err
	}); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(command, deployment)
}

func (r *ReconcileCommand) getPostgres(command *contrail.Command) (*contrail.Postgres, error) {
	psql := &contrail.Postgres{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.PostgresInstance,
	}, psql)

	return psql, err
}

func newDeployment(name, namespace, configVolumeName string, containers map[string]*contrail.Container) *apps.Deployment {
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
						Image:           getImage(containers, "api"),
						Command:         getCommand(containers, "api"),
						ReadinessProbe: &core.Probe{
							Handler: core.Handler{
								HTTPGet: &core.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 9091},
								},
							},
						},
						//TODO: Command should support CA certificates
						VolumeMounts: []core.VolumeMount{{
							Name:      configVolumeName,
							MountPath: "/etc/contrail",
						}},
					}},
					InitContainers: []core.Container{{
						Name:            "command-init",
						ImagePullPolicy: core.PullAlways,
						Image:           getImage(containers, "init"),
						Command:         getCommand(containers, "init"),
						VolumeMounts: []core.VolumeMount{{
							Name:      configVolumeName,
							MountPath: "/etc/contrail",
						}},
					}},
					DNSPolicy: core.DNSClusterFirst,
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
		},
	}
}

func getImage(containers map[string]*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"init": "localhost:5000/contrail-command",
		"api":  "localhost:5000/contrail-command",
	}

	c, ok := containers[containerName]
	if ok == false || c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(containers map[string]*contrail.Container, containerName string) []string {
	var defaultContainersCommand = map[string][]string{
		"init": []string{"bash", "/etc/contrail/bootstrap.sh"},
		"api":  []string{"bash", "/etc/contrail/entrypoint.sh"},
	}

	c, ok := containers[containerName]
	if ok == false || c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}

func (r *ReconcileCommand) updateStatus(
	command *contrail.Command,
	deployment *apps.Deployment,
) error {
	command.Status.Active = false
	intendentReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		intendentReplicas = *deployment.Spec.Replicas
	}

	if deployment.Status.ReadyReplicas == intendentReplicas {
		command.Status.Active = true
	}

	return r.client.Status().Update(context.Background(), command)
}
