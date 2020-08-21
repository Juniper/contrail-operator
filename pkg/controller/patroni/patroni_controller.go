package patroni

import (
	"context"

	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
)

var log = logf.Log.WithName("controller_patroni")

// Add creates a new Patroni Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePatroni{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("patroni-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Patroni
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Patroni{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Patroni
	err = c.Watch(&source.Kind{Type: &core.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.Patroni{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePatroni implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePatroni{}

// NewReconciler is used to create a new ReconcilePatroni
func NewReconciler(client client.Client, scheme *runtime.Scheme) *ReconcilePatroni {
	return &ReconcilePatroni{
		client:     client,
		scheme:     scheme,
		kubernetes: k8s.New(client, scheme),
	}
}

// ReconcilePatroni reconciles a Patroni object
type ReconcilePatroni struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Patroni object and makes changes based on the state read
// and what is in the Patroni.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePatroni) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Patroni")

	// Fetch the Patroni instance
	instance := &contrailv1alpha1.Patroni{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	if err := r.ensureLabelExists(instance); err != nil {
		return reconcile.Result{}, err
	}

	err = r.ensureServicesExist(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.ensureServiceAccountExists(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.ensureRoleExists(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.ensureRoleBindingExists(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// TODO - create Service, service account, role and role binding
	// TODO - create STS -
	// TODO - create PVs and PVCs
	// TODO - create secret with passwords
	// TODO - create certs, mount them to sts and point their location in patroni container
	// TODO - implement master election
	//
	// REQUIREMENTS:
	// - host networking
	// - master exposed as ClusterIP service
	// - patroni pods have to have label: "patroni": instance.Name for Service selector

	statefulSet, err := r.createOrUpdateSts(request, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	instance.Status.Active = false
	intendentReplicas := int32(1)
	if statefulSet.Spec.Replicas != nil {
		intendentReplicas = *statefulSet.Spec.Replicas
	}

	if statefulSet.Status.ReadyReplicas == intendentReplicas {
		instance.Status.Active = true
	}

	return reconcile.Result{}, nil
}

func (r *ReconcilePatroni) ensureLabelExists(p *contrailv1alpha1.Patroni) error {
	if len(p.Labels) != 0 {
		return nil
	}

	p.Labels = contraillabel.New(contrailv1alpha1.PatroniInstanceType, p.Name)
	return r.client.Update(context.Background(), p)
}

func (r *ReconcilePatroni) ensureServicesExist(instance *contrailv1alpha1.Patroni) error {
	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-service",
			Namespace: instance.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, service, func() error {
		service.ObjectMeta.Labels = map[string]string{"app": "patroni", "cluster-manager": "patroni", "patroni": instance.Name}
		service.Spec = core.ServiceSpec{
			Type: core.ServiceTypeClusterIP,
			Ports: []core.ServicePort{{
				Port: 5432,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 5432,
				},
			}},
		}
		return controllerutil.SetControllerReference(instance, service, r.scheme)
	})

	if err != nil {
		return err
	}

	endpoints := &core.Endpoints{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-service",
			Namespace: instance.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, endpoints, func() error {
		endpoints.ObjectMeta.Labels = map[string]string{"app": "patroni", "cluster-manager": "patroni", "patroni": instance.Name}
		return controllerutil.SetControllerReference(instance, service, r.scheme)
	})

	if err != nil {
		return err
	}

	serviceRepl := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-service-replica",
			Namespace: instance.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serviceRepl, func() error {
		serviceRepl.ObjectMeta.Labels = map[string]string{"app": "patroni", "role": "replica", "cluster-manager": "patroni", "patroni": instance.Name}
		serviceRepl.Spec = core.ServiceSpec{
			Selector: map[string]string{
				"app":              "patroni",
				"contrail_manager": "patroni",
				"patroni":          instance.Name,
				"role":             "replica",
			},
			Type: core.ServiceTypeClusterIP,
			Ports: []core.ServicePort{{
				Port: 5432,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 5432,
				},
			}},
		}
		return controllerutil.SetControllerReference(instance, serviceRepl, r.scheme)
	})

	return err
}

func (r *ReconcilePatroni) ensureServiceAccountExists(instance *contrailv1alpha1.Patroni) error {
	serviceAccount := &core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-service-account",
			Namespace: instance.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serviceAccount, func() error {
		serviceAccount.ObjectMeta.Labels = map[string]string{"app": "patroni", "cluster-manager": "patroni", "patroni": instance.Name}
		return controllerutil.SetControllerReference(instance, serviceAccount, r.scheme)
	})

	return err
}

func (r *ReconcilePatroni) ensureRoleExists(instance *contrailv1alpha1.Patroni) error {
	role := &rbac.Role{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-role",
			Namespace: instance.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, role, func() error {
		role.ObjectMeta.Labels = map[string]string{"app": "patroni", "cluster-manager": "patroni", "patroni": instance.Name}
		role.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs: []string{
					"create",
					"get",
					"list",
					"patch",
					"update",
					"watch",
					"delete",
					"deletecollection",
				},
				Resources: []string{"configmaps"},
			},
			{
				APIGroups: []string{""},
				Verbs: []string{
					"get",
					"patch",
					"update",
					"create",
					"list",
					"watch",
					"delete",
					"deletecollection",
				},
				Resources: []string{"endpoints"},
			},
			{
				APIGroups: []string{""},
				Verbs: []string{
					"get",
					"list",
					"patch",
					"update",
					"watch",
				},
				Resources: []string{"pods"},
			},
			{
				APIGroups: []string{""},
				Verbs: []string{
					"create",
				},
				Resources: []string{"services"},
			},
		}
		return controllerutil.SetControllerReference(instance, role, r.scheme)
	})

	return err
}
func (r *ReconcilePatroni) ensureRoleBindingExists(instance *contrailv1alpha1.Patroni) error {
	rb := &rbac.RoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Name:      instance.Name + "-patroni-role-binding",
			Namespace: instance.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, rb, func() error {
		rb.ObjectMeta.Labels = map[string]string{"app": "patroni", "cluster-manager": "patroni", "patroni": instance.Name}
		rb.Subjects = []rbac.Subject{{
			Kind: "ServiceAccount",
			Name: instance.Name + "-patroni-service-account",
		}}
		rb.RoleRef = rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     instance.Name + "-patroni-role",
		}
		return controllerutil.SetControllerReference(instance, rb, r.scheme)
	})

	return err
}

func (r *ReconcilePatroni) createOrUpdateSts(request reconcile.Request, instance *contrailv1alpha1.Patroni) (*apps.StatefulSet, error) {

	statefulSet := GetSTS(instance, "patroni-service-account")

	statefulSet.Namespace = request.Namespace
	statefulSet.Name = request.Name + "-statefulset"

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {

		contrailv1alpha1.SetSTSCommonConfiguration(statefulSet, &instance.Spec.CommonConfiguration)

		var patroniGroupId int64 = 0
		statefulSet.Spec.Template.Spec.SecurityContext = &core.PodSecurityContext{}
		statefulSet.Spec.Template.Spec.SecurityContext.FSGroup = &patroniGroupId
		statefulSet.Spec.Template.Spec.SecurityContext.RunAsGroup = &patroniGroupId
		statefulSet.Spec.Template.Spec.SecurityContext.RunAsUser = &patroniGroupId

		storageClassName := "local-storage"
		statefulSet.Spec.VolumeClaimTemplates = []core.PersistentVolumeClaim{
			{
				ObjectMeta: meta.ObjectMeta{
					Name:      "storage-device",
					Namespace: request.Namespace,
					Labels:    instance.Labels,
				},
				Spec: core.PersistentVolumeClaimSpec{
					AccessModes: []core.PersistentVolumeAccessMode{
						core.ReadWriteOnce,
					},
					StorageClassName: &storageClassName,
					Resources: core.ResourceRequirements{
						Requests: map[core.ResourceName]resource.Quantity{
							core.ResourceStorage: resource.MustParse("5Gi"),
						},
					},
				},
			},
		}

		statefulSet.Spec.Selector = &meta.LabelSelector{MatchLabels: instance.Labels}
		return controllerutil.SetControllerReference(instance, statefulSet, r.scheme)
	})
	return statefulSet, err
}

type containerGenerator struct {
	containersSpec []*contrailv1alpha1.Container
	device         string
}
