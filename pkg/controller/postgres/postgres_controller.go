package postgres

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
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
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
	"github.com/Juniper/contrail-operator/pkg/localvolume"
)

var log = logf.Log.WithName("controller_postgres")

const defaultPostgresStoragePath = "/mnt/postgres"

// Add creates a new Postgres Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePostgres{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme()),
		volumes:    localvolume.New(mgr.GetClient()),
	}
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

	// Watch for changes to secondary resource StatefulSet and requeue the owner Postgres
	err = c.Watch(&source.Kind{Type: &apps.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Postgres{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &core.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Postgres{},
	})

	return err
}

// blank assignment to verify that ReconcilePostgres implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePostgres{}

// NewReconciler is used to create a new ReconcilePostgres
func NewReconciler(client client.Client, scheme *runtime.Scheme) *ReconcilePostgres {
	return &ReconcilePostgres{
		client:     client,
		scheme:     scheme,
		kubernetes: k8s.New(client, scheme),
		volumes:    localvolume.New(client),
	}
}

// ReconcilePostgres reconciles a Postgres object
type ReconcilePostgres struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	volumes    localvolume.Volumes
}

// Reconcile reads that state of the cluster for a Postgres object and makes changes based on the state read
// and what is in the Postgres.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePostgres) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Postgres")

	// Fetch the Postgres postgres
	postgres := &contrail.Postgres{}
	err := r.client.Get(context.TODO(), request.NamespacedName, postgres)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if !postgres.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	if err := r.ensureLabelExists(postgres); err != nil {
		return reconcile.Result{}, err
	}

	postgresService, err := r.ensureServicesExist(postgres)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = contrail.CreateAccount("postgres", request.Namespace, r.client, r.scheme, postgres); err != nil {
		return reconcile.Result{}, err
	}

	postgresPods, err := r.listPostgresPods(postgres)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list Postgres pods: %v", err)
	}

	leaderClusterIP := postgresService.Spec.ClusterIP

	if err := r.ensureCertificatesExist(postgres, postgresPods, leaderClusterIP); err != nil {
		return reconcile.Result{}, err
	}

	if len(postgresPods.Items) > 0 {
		err = contrail.SetPodsToReady(postgresPods, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	if err = r.ensureLocalPVsExist(postgres); err != nil {
		return reconcile.Result{}, err
	}

	replicationPassSecretName := postgres.Name + "-postgres-replication-secret"
	if postgres.Spec.ServiceConfiguration.ReplicationPassSecretName != "" {
		replicationPassSecretName = postgres.Spec.ServiceConfiguration.ReplicationPassSecretName
	}

	if err = r.replicationPassSecret(replicationPassSecretName, "postgres", postgres).ensureExists(); err != nil {
		return reconcile.Result{}, err
	}
	serviceAccountName := "serviceaccount-postgres"
	rootPassSecretName := postgres.Spec.ServiceConfiguration.RootPassSecretName
	statefulSet, err := r.createOrUpdateSts(postgres, postgresService, replicationPassSecretName, rootPassSecretName, serviceAccountName)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.ensurePVCOwnershipExists(postgres); err != nil {
		return reconcile.Result{}, err
	}

	postgres.Status.Endpoint = leaderClusterIP
	postgres.Status.Active = false
	intendentReplicas := int32(1)
	if statefulSet.Spec.Replicas != nil {
		intendentReplicas = *statefulSet.Spec.Replicas
	}

	if statefulSet.Status.ReadyReplicas == intendentReplicas {
		postgres.Status.Active = true
	}

	return reconcile.Result{}, r.client.Status().Update(context.Background(), postgres)
}

func (r *ReconcilePostgres) ensureLabelExists(p *contrail.Postgres) error {
	if len(p.Labels) != 0 {
		return nil
	}

	p.Labels = contraillabel.New(contrail.PostgresInstanceType, p.Name)
	return r.client.Update(context.Background(), p)
}

func (r *ReconcilePostgres) ensureCertificatesExist(postgres *contrail.Postgres, pods *core.PodList, serviceIP string) error {
	hostNetwork := true
	if postgres.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *postgres.Spec.CommonConfiguration.HostNetwork
	}
	return certificates.NewCertificateWithServiceIP(r.client, r.scheme, postgres, pods, serviceIP, "postgres", hostNetwork).EnsureExistsAndIsSigned()
}

func (r *ReconcilePostgres) listPostgresPods(postgres *contrail.Postgres) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(postgres.Labels)
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}

func (r *ReconcilePostgres) ensureServicesExist(postgres *contrail.Postgres) (*core.Service, error) {
	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      postgres.Name,
			Namespace: postgres.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, service, func() error {
		service.ObjectMeta.Labels = postgres.Labels
		service.Spec.Type = core.ServiceTypeClusterIP
		nodePort := int32(0)
		listenPort := int32(postgres.Spec.ServiceConfiguration.ListenPort)
		for i, p := range service.Spec.Ports {
			if p.Port == listenPort {
				nodePort = service.Spec.Ports[i].NodePort
			}
		}
		service.Spec.Ports = []core.ServicePort{
			{Port: listenPort, Protocol: "TCP", NodePort: nodePort},
		}
		return controllerutil.SetControllerReference(postgres, service, r.scheme)
	})

	if err != nil {
		return nil, err
	}

	serviceRepl := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      postgres.Name + "-replica",
			Namespace: postgres.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serviceRepl, func() error {
		labels := copyStringMap(postgres.Labels)
		labels["role"] = "replica"
		serviceRepl.ObjectMeta.Labels = labels
		serviceRepl.Spec.Selector = labels
		serviceRepl.Spec.Type = core.ServiceTypeClusterIP
		nodePort := int32(0)
		listenPort := int32(postgres.Spec.ServiceConfiguration.ListenPort)
		for i, p := range serviceRepl.Spec.Ports {
			if p.Port == listenPort {
				nodePort = serviceRepl.Spec.Ports[i].NodePort
			}
		}
		serviceRepl.Spec.Ports = []core.ServicePort{
			{Port: listenPort, Protocol: "TCP", NodePort: nodePort},
		}
		return controllerutil.SetControllerReference(postgres, serviceRepl, r.scheme)
	})

	return service, err
}

func (r *ReconcilePostgres) createOrUpdateSts(postgres *contrail.Postgres, service *core.Service,
	replicationPassSecretName, rootPassSecretName, serviceAccountName string) (*apps.StatefulSet, error) {

	statefulSet := &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      postgres.Name + "-statefulset",
			Namespace: postgres.Namespace,
		},
	}

	storagePath := postgres.Spec.ServiceConfiguration.Storage.Path
	if storagePath == "" {
		storagePath = defaultPostgresStoragePath
	}

	var (
		initHostPathType            = core.HostPathDirectoryOrCreate
		postgresUID           int64 = 999
		labelsMountPermission int32 = 0644
		csrSignerCaVolumeName       = postgres.Name + "-csr-signer-ca"
		storageClassName            = "local-storage"
	)

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {
		statefulSet.Labels = postgres.Labels
		contrail.SetSTSCommonConfiguration(statefulSet, &postgres.Spec.CommonConfiguration)
		statefulSet.Spec.Selector = &meta.LabelSelector{MatchLabels: postgres.Labels}
		statefulSet.Spec.ServiceName = service.Name
		statefulSet.Spec.Replicas = postgres.Spec.CommonConfiguration.Replicas
		statefulSet.Spec.Template.Labels = postgres.Labels
		statefulSet.Spec.Template.Spec.Affinity = &core.Affinity{
			PodAntiAffinity: &core.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
					LabelSelector: &meta.LabelSelector{MatchLabels: postgres.Labels},
					TopologyKey:   "kubernetes.io/hostname",
				}},
			},
		}
		statefulSet.Spec.Template.Spec.InitContainers = r.initContainers(postgres)
		statefulSet.Spec.Template.Spec.Containers = r.containers(postgres, rootPassSecretName, replicationPassSecretName, csrSignerCaVolumeName)
		statefulSet.Spec.Template.Spec.ServiceAccountName = serviceAccountName
		statefulSet.Spec.Template.Spec.SecurityContext = &core.PodSecurityContext{
			RunAsUser:          &postgresUID,
			FSGroup:            &postgresUID,
			SupplementalGroups: []int64{999, 1000},
		}
		statefulSet.Spec.Template.Spec.Volumes = []core.Volume{
			{
				Name: "postgres-storage-init",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: storagePath,
						Type: &initHostPathType,
					},
				},
			},
			{
				Name: postgres.Name + "-secret-certificates",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: postgres.Name + "-secret-certificates",
					},
				},
			},
			{
				Name: csrSignerCaVolumeName,
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: certificates.SignerCAConfigMapName,
						},
					},
				},
			},
			{
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
		}
		statefulSet.Spec.VolumeClaimTemplates = []core.PersistentVolumeClaim{
			{
				ObjectMeta: meta.ObjectMeta{
					Name:      "pgdata",
					Namespace: postgres.Namespace,
					Labels:    postgres.Labels,
				},
				Spec: core.PersistentVolumeClaimSpec{
					AccessModes: []core.PersistentVolumeAccessMode{
						core.ReadWriteOnce,
					},
					Selector: &meta.LabelSelector{
						MatchLabels: postgres.Labels,
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

		return controllerutil.SetControllerReference(postgres, statefulSet, r.scheme)
	})
	return statefulSet, err
}

func (r *ReconcilePostgres) initContainers(postgres *contrail.Postgres) []core.Container {
	return []core.Container{
		{
			Name:            "wait-for-ready-conf",
			Image:           getImage(postgres.Spec.ServiceConfiguration.Containers, "wait-for-ready-conf"),
			Command:         getCommand(postgres.Spec.ServiceConfiguration.Containers, "wait-for-ready-conf"),
			ImagePullPolicy: core.PullAlways,
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "status",
					MountPath: "/tmp/podinfo",
				},
			},
		},
		{
			Name:            "init",
			Image:           getImage(postgres.Spec.ServiceConfiguration.Containers, "init"),
			Command:         []string{"/bin/sh", "-c", "if [[ -d /mnt/postgres/postgres ]]; then chmod 0750 /mnt/postgres/postgres; fi"},
			ImagePullPolicy: "Always",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "postgres-storage-init",
					ReadOnly:  false,
					MountPath: "/mnt/",
				},
			},
		},
	}
}

func (r *ReconcilePostgres) containers(postgres *contrail.Postgres, rootPassSecretName, replicationPassSecretName, csrSignerCaVolumeName string) []core.Container {
	return []core.Container{
		{
			Name:  "patroni",
			Image: getImage(postgres.Spec.ServiceConfiguration.Containers, "patroni"),
			ReadinessProbe: &core.Probe{
				Handler: core.Handler{
					HTTPGet: &core.HTTPGetAction{
						Scheme: core.URISchemeHTTP,
						Path:   "/readiness",
						Port:   intstr.IntOrString{IntVal: 8008},
					},
				},
			},
			Env:             r.containerEnv(postgres.Name, rootPassSecretName, replicationPassSecretName),
			ImagePullPolicy: "Always",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "pgdata",
					ReadOnly:  false,
					MountPath: "/var/lib/postgresql/data",
					SubPath:   "postgres",
				},
				{
					Name:      postgres.Name + "-secret-certificates",
					MountPath: "/var/lib/ssl_certificates",
				},
				{
					Name:      csrSignerCaVolumeName,
					MountPath: certificates.SignerCAMountPath,
				},
			},
		},
	}
}

func (r *ReconcilePostgres) containerEnv(name, rootPassSecretName, replicationPassSecretName string) []core.EnvVar {
	var podIPEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var nameEnv = core.EnvVar{
		Name: "PATRONI_NAME",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	}

	var scopeEnv = core.EnvVar{
		Name:  "PATRONI_SCOPE",
		Value: name,
	}

	var namespaceEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_NAMESPACE",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}

	var scopeLabelEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_SCOPE_LABEL",
		Value: "postgres",
	}

	var labelsEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_LABELS",
		Value: contraillabel.AsString("postgres", name),
	}

	var postgresListenAddressEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_LISTEN",
		Value: "0.0.0.0:5432",
	}

	var restApiListenAddressEnv = core.EnvVar{
		Name:  "PATRONI_RESTAPI_LISTEN",
		Value: "0.0.0.0:8008",
	}

	var replicationUserEnv = core.EnvVar{
		Name:  "PATRONI_REPLICATION_USERNAME",
		Value: "standby",
	}

	var replicationPassEnv = core.EnvVar{
		Name: "PATRONI_REPLICATION_PASSWORD",
		ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{
					Name: replicationPassSecretName,
				},
				Key: "replication-password",
			},
		},
	}

	var superuserEnv = core.EnvVar{
		Name:  "PATRONI_SUPERUSER_USERNAME",
		Value: "root",
	}

	var postgresDBEnv = core.EnvVar{
		Name:  "POSTGRES_DB",
		Value: "contrail_test",
	}

	var superuserPassEnv = core.EnvVar{
		Name: "PATRONI_SUPERUSER_PASSWORD",
		ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{
					Name: rootPassSecretName,
				},
				Key: "password",
			},
		},
	}

	var endpointsEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_USE_ENDPOINTS",
		Value: "true",
	}

	var dataDirEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_DATA_DIR",
		Value: "/var/lib/postgresql/data/postgres",
	}

	var pgpassEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_PGPASS",
		Value: "/tmp/pgpass",
	}

	return []core.EnvVar{
		nameEnv,
		scopeEnv,
		podIPEnv,
		namespaceEnv,
		scopeLabelEnv,
		labelsEnv,
		endpointsEnv,
		replicationUserEnv,
		replicationPassEnv,
		superuserEnv,
		superuserPassEnv,
		dataDirEnv,
		postgresDBEnv,
		postgresListenAddressEnv,
		restApiListenAddressEnv,
		pgpassEnv,
	}
}

func (r *ReconcilePostgres) ensureLocalPVsExist(postgres *contrail.Postgres) error {
	path := postgres.Spec.ServiceConfiguration.Storage.Path
	size := postgres.Spec.ServiceConfiguration.Storage.Size
	var storage resource.Quantity
	var err error
	if size == "" {
		storage = resource.MustParse("5Gi")
	} else {
		storage, err = resource.ParseQuantity(size)
		if err != nil {
			return err
		}
	}

	if path == "" {
		path = defaultPostgresStoragePath
	}

	for i := int32(0); i < postgres.Spec.CommonConfiguration.GetReplicas(); i++ {
		name := fmt.Sprintf("%v-%v-postgres-data-%v", postgres.Name, postgres.Namespace, i)
		nodeSelectors := postgres.Spec.CommonConfiguration.NodeSelector
		lv, err := r.volumes.New(name, path, storage, postgres.Labels, nodeSelectors)
		if err != nil {
			return err
		}
		if err := lv.EnsureExists(); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcilePostgres) ensurePVCOwnershipExists(postgres *contrail.Postgres) error {
	listOps := &client.ListOptions{Namespace: postgres.Namespace, LabelSelector: labels.SelectorFromSet(postgres.Labels)}
	pvcList := &core.PersistentVolumeClaimList{}
	if err := r.client.List(context.TODO(), pvcList, listOps); err != nil {
		return err
	}
	for _, pvc := range pvcList.Items {
		if err := controllerutil.SetControllerReference(postgres, &pvc, r.scheme); err != nil {
			return err
		}
		if err := r.client.Update(context.TODO(), &pvc); err != nil {
			return err
		}
	}
	return nil
}

func getImage(containers []*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"patroni":             "localhost:5000/patroni:1.6.5.logical",
		"init":                "localhost:5000/busybox:1.31",
		"wait-for-ready-conf": "localhost:5000/busybox:1.31",
	}
	c := utils.GetContainerFromList(containerName, containers)
	if c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(containers []*contrail.Container, containerName string) []string {
	var defaultContainersCommand = map[string][]string{
		"wait-for-ready-conf": {"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
	}

	c := utils.GetContainerFromList(containerName, containers)
	if c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}

func copyStringMap(m map[string]string) map[string]string {
	newMap := make(map[string]string)

	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}
