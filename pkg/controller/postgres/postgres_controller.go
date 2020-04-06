package postgres

import (
	"context"
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

var log = logf.Log.WithName("controller_postgres")

// Add creates a new Postgres Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	client := mgr.GetClient()
	scheme := mgr.GetScheme()
	claims := volumeclaims.New(client, scheme)
	config := mgr.GetConfig()
	return &ReconcilePostgres{client: client, scheme: scheme, claims: claims, config: config}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("postgres-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Postgres
	err = c.Watch(&source.Kind{Type: &contrail.Postgres{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Postgres
	err = c.Watch(&source.Kind{Type: &core.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Postgres{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePostgres implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePostgres{}

// ReconcilePostgres reconciles a Postgres object
type ReconcilePostgres struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	claims volumeclaims.PersistentVolumeClaims
	config *rest.Config
}

// Reconcile reads that state of the cluster for a Postgres object and makes changes based on the state read
// and what is in the Postgres.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePostgres) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Postgres")

	// Fetch the Postgres instance
	instance := &contrail.Postgres{}
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

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	postgresPods, err := r.listPostgresPods(instance.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list postgres pods: %v", err)
	}

	if err := r.ensureCertificatesExist(instance, postgresPods); err != nil {
		return reconcile.Result{}, err
	}

	if len(postgresPods.Items) > 0 {
		err = contrail.SetPodsToReady(postgresPods, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	claimName := types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      instance.Name + "-pv-claim",
	}
	claim := r.claims.New(claimName, instance)
	claim.SetStoragePath(instance.Spec.Storage.Path)
	if instance.Spec.Storage.Size != "" {
		var quantity resource.Quantity
		quantity, err = instance.Spec.Storage.SizeAsQuantity()
		if err != nil {
			return reconcile.Result{}, err
		}
		claim.SetStorageSize(quantity)
	}
	claim.SetNodeSelector(map[string]string{"node-role.kubernetes.io/master": ""})
	if err = claim.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance, claimName.Name)

	// Set Postgres instance as the owner and controller
	if err = controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &core.Pod{}
	err = r.client.Get(context.Background(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.Background(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(instance, pod)
}

func (r *ReconcilePostgres) updateStatus(
	postgres *contrail.Postgres, pod *core.Pod,
) error {
	err := r.client.Get(context.Background(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, pod)
	if err != nil {
		return err
	}

	if len(pod.Status.ContainerStatuses) != 0 {
		postgres.Status.Active = pod.Status.ContainerStatuses[0].Ready
		postgres.Status.Node = pod.Status.PodIP + ":5432"
	} else {
		postgres.Status.Active = false
	}
	return r.client.Status().Update(context.Background(), postgres)
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *contrail.Postgres, claimName string) *core.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	db := "contrail_test"
	var labelsMountPermission int32 = 0644
	return &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: core.PodSpec{
			HostNetwork:  true,
			NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			DNSPolicy:    core.DNSClusterFirst,
			InitContainers: []core.Container{
				{
					Name:            "wait-for-ready-conf",
					Image:           getImage(cr.Spec.Containers, "wait-for-ready-conf"),
					Command:         getCommand(cr.Spec.Containers, "wait-for-ready-conf"),
					ImagePullPolicy: core.PullAlways,
					VolumeMounts: []core.VolumeMount{{
						Name:      "status",
						MountPath: "/tmp/podinfo",
					}},
				},
			},
			Containers: []core.Container{
				{
					Image:           getImage(cr.Spec.Containers, "postgres"),
					Command:         getCommand(cr.Spec.Containers, "postgres"),
					Name:            "postgres",
					ImagePullPolicy: core.PullAlways,
					ReadinessProbe: &core.Probe{
						Handler: core.Handler{
							Exec: &core.ExecAction{
								Command: []string{"pg_isready", "-h", "localhost", "-U", "root", "-d", db},
							},
						},
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      cr.Name + "-volume",
							MountPath: "/var/lib/postgresql/data",
							SubPath:   "postgres",
						},
						{
							Name:      cr.Name + "-secret-certificates",
							MountPath: "/var/lib/postgresql/ssl",
						},
					},
					Env: []core.EnvVar{
						{Name: "POSTGRES_USER", Value: "root"},
						{Name: "POSTGRES_PASSWORD", Value: "contrail123"},
						{Name: "POSTGRES_DB", Value: db},
						{
							Name: "MY_POD_IP",
							ValueFrom: &core.EnvVarSource{
								FieldRef: &core.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
					},
				},
			},
			Volumes: []core.Volume{
				core.Volume{
					Name: cr.Name + "-volume",
					VolumeSource: core.VolumeSource{
						PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
						},
					},
				},
				core.Volume{
					Name: cr.Name + "-secret-certificates",
					VolumeSource: core.VolumeSource{
						Secret: &core.SecretVolumeSource{
							SecretName: cr.Name + "-secret-certificates",
						},
					},
				},
				core.Volume{
					Name: "status",
					VolumeSource: core.VolumeSource{
						DownwardAPI: &core.DownwardAPIVolumeSource{
							Items: []core.DownwardAPIVolumeFile{
								{
									FieldRef: &core.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "metadata.labels",
									},
									Path: "pod_labels",
								},
							},
							DefaultMode: &labelsMountPermission,
						},
					},
				},
			},
			Tolerations: []core.Toleration{
				{Operator: "Exists", Effect: "NoSchedule"},
				{Operator: "Exists", Effect: "NoExecute"},
			},
		},
	}

}

func (r *ReconcilePostgres) ensureCertificatesExist(postgres *contrail.Postgres, pods *core.PodList) error {
	hostNetwork := true
	if postgres.Spec.HostNetwork != nil {
		hostNetwork = *postgres.Spec.HostNetwork
	}
	return certificates.New(r.client, r.scheme, postgres, r.config, pods, "postgres", hostNetwork).EnsureExistsAndIsSigned()
}

func (r *ReconcilePostgres) listPostgresPods(postgresName string) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"postgres": postgresName})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}

func getImage(containers map[string]*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"postgres":            "localhost:5000/contrail-command",
		"wait-for-ready-conf": "localhost:5000/busybox",
	}

	c, ok := containers[containerName]
	if !ok || c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(containers map[string]*contrail.Container, containerName string) []string {
	var defaultContainersCommand = map[string][]string{
		"postgres":            {"/bin/bash", "-c", "docker-entrypoint.sh -c wal_level=logical -c ssl=on -c ssl_cert_file=/var/lib/postgresql/ssl/server-${MY_POD_IP}.crt -c ssl_key_file=/var/lib/postgresql/ssl/server-key-${MY_POD_IP}.pem"},
		"wait-for-ready-conf": {"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
	}

	c, ok := containers[containerName]
	if !ok || c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}
