package patroni

import (
	"context"
	"encoding/json"
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrailcertificates "github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
	"github.com/Juniper/contrail-operator/pkg/localvolume"
)

var log = logf.Log.WithName("controller_patroni")

const defaultPatroniStoragePath = "/mnt/patroni"

// Add creates a new Patroni Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePatroni{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme()),
		volumes:    localvolume.New(mgr.GetClient()),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("patroni-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Patroni
	err = c.Watch(&source.Kind{Type: &contrail.Patroni{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSet and requeue the owner Patroni
	err = c.Watch(&source.Kind{Type: &apps.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Patroni{},
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
		volumes:    localvolume.New(client),
	}
}

// ReconcilePatroni reconciles a Patroni object
type ReconcilePatroni struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	volumes    localvolume.Volumes
}

// Reconcile reads that state of the cluster for a Patroni object and makes changes based on the state read
// and what is in the Patroni.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePatroni) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Patroni")

	// Fetch the Patroni patroni
	patroni := &contrail.Patroni{}
	err := r.client.Get(context.TODO(), request.NamespacedName, patroni)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if err := r.ensureLabelExists(patroni); err != nil {
		return reconcile.Result{}, err
	}

	patroniService, err := r.ensureServicesExist(patroni)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = contrail.CreateAccount("patroni", request.Namespace, r.client, r.scheme, patroni); err != nil {
		return reconcile.Result{}, err
	}

	patroniPods, err := r.listPatroniPods(patroni.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list command pods: %v", err)
	}

	if err := r.ensureCertificatesExist(patroni, patroniPods, patroniService.Spec.ClusterIP); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.ensureLocalPVsExist(patroni); err != nil {
		return reconcile.Result{}, err
	}

	// TODO - implement master election

	credentialsSecretName := patroni.Name + "-patroni-credentials-secret"
	if patroni.Spec.ServiceConfiguration.CredentialsSecretName != "" {
		credentialsSecretName = patroni.Spec.ServiceConfiguration.CredentialsSecretName
	}

	if err = r.credentialsSecret(credentialsSecretName, "patroni", patroni).ensureExists(); err != nil {
		return reconcile.Result{}, err
	}
	serviceAccountName := "serviceaccount-patroni"
	statefulSet, err := r.createOrUpdateSts(patroni, patroniService, credentialsSecretName, serviceAccountName)
	if err != nil {
		return reconcile.Result{}, err
	}

	patroni.Status.CredentialsSecretName = credentialsSecretName
	patroni.Status.ClusterIP = patroniService.Spec.ClusterIP
	patroni.Status.Active = false
	intendentReplicas := int32(1)
	if statefulSet.Spec.Replicas != nil {
		intendentReplicas = *statefulSet.Spec.Replicas
	}

	if statefulSet.Status.ReadyReplicas == intendentReplicas {
		patroni.Status.Active = true
	}

	return reconcile.Result{}, nil
}

func (r *ReconcilePatroni) ensureLabelExists(p *contrail.Patroni) error {
	if len(p.Labels) != 0 {
		return nil
	}

	p.Labels = contraillabel.New(contrail.PatroniInstanceType, p.Name)
	return r.client.Update(context.Background(), p)
}

func (r *ReconcilePatroni) ensureCertificatesExist(patroni *contrail.Patroni, pods *core.PodList, serviceIP string) error {
	hostNetwork := true
	if patroni.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *patroni.Spec.CommonConfiguration.HostNetwork
	}
	return contrailcertificates.NewCertificateWithServiceIP(r.client, r.scheme, patroni, pods, serviceIP, "patroni", hostNetwork).EnsureExistsAndIsSigned()
}

func (r *ReconcilePatroni) listPatroniPods(patroniName string) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(contraillabel.New(contrail.PatroniInstanceType, patroniName))
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}

func (r *ReconcilePatroni) ensureServicesExist(patroni *contrail.Patroni) (*core.Service, error) {
	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      patroni.Name + "-patroni-service",
			Namespace: patroni.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, service, func() error {
		service.ObjectMeta.Labels = contraillabel.New(patroni.Kind, patroni.Name)
		service.Spec.Type = core.ServiceTypeClusterIP
		service.Spec.Ports = []core.ServicePort{{
			Port: 5432,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: 5432,
			},
		}}
		return controllerutil.SetControllerReference(patroni, service, r.scheme)
	})

	if err != nil {
		return nil, err
	}

	endpoints := &core.Endpoints{
		ObjectMeta: meta.ObjectMeta{
			Name:      patroni.Name + "-patroni-service",
			Namespace: patroni.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, endpoints, func() error {
		endpoints.ObjectMeta.Labels = contraillabel.New(patroni.Kind, patroni.Name)
		return controllerutil.SetControllerReference(patroni, service, r.scheme)
	})

	if err != nil {
		return nil, err
	}

	serviceRepl := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      patroni.Name + "-patroni-service-replica",
			Namespace: patroni.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serviceRepl, func() error {
		labels := contraillabel.New(patroni.Kind, patroni.Name)
		labels["role"] = "replica"
		serviceRepl.ObjectMeta.Labels = labels
		serviceRepl.Spec.Selector = labels
		serviceRepl.Spec.Type = core.ServiceTypeClusterIP
		serviceRepl.Spec.Ports = []core.ServicePort{{
			Port: 5432,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: 5432,
			},
		}}
		return controllerutil.SetControllerReference(patroni, serviceRepl, r.scheme)
	})

	return service, err
}


func (r *ReconcilePatroni) createOrUpdateSts(patroni *contrail.Patroni, service *core.Service, passSecretName, serviceAccountName string) (*apps.StatefulSet, error) {

	statefulSet := &apps.StatefulSet{}
	statefulSet.Namespace = patroni.Namespace
	statefulSet.Name = patroni.Name + "-statefulset"
	csrSignerCaVolumeName := patroni.Name + "-csr-signer-ca"
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
		Name: "PATRONI_SCOPE",
		Value: patroni.Name,
	}

	var namespaceEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_NAMESPACE",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}

	var labelsEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_LABELS",
		Value: contraillabel.AsString("patroni", patroni.Name),
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
					Name: passSecretName,
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
					Name: passSecretName,
				},
				Key: "superuser-password",
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

	var podContainers = []core.Container{
		{
			Name:  "patroni",
			Image: getImage(patroni.Spec.ServiceConfiguration.Containers, "patroni"),
			Env: []core.EnvVar{
				nameEnv,
				scopeEnv,
				podIPEnv,
				namespaceEnv,
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
			},
			ImagePullPolicy: "Always",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "pgdata",
					ReadOnly:  false,
					MountPath: "/var/lib/postgresql/data",
					SubPath:   "postgres",
				},
				{
					Name:      patroni.Name + "-secret-certificates",
					MountPath: "/var/lib/ssl_certificates",
				},
				{
					Name:      csrSignerCaVolumeName,
					MountPath: contrailcertificates.SignerCAMountPath,
				},
			},
		},
	}

	var podInitContainers = []core.Container{
		{
			Name:            "init",
			Image:           getImage(patroni.Spec.ServiceConfiguration.Containers, "init"),
			ImagePullPolicy: "Always",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "patroni-storage-init",
					ReadOnly:  false,
					MountPath: "/mnt/",
				},
			},
		},
	}

	storagePath := patroni.Spec.ServiceConfiguration.Storage.Path
	if storagePath == "" {
		storagePath = defaultPatroniStoragePath
	}
	initHostPathType := core.HostPathDirectoryOrCreate
	var podSpec = core.PodSpec{
		Affinity: &core.Affinity{
			PodAntiAffinity: &core.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
					LabelSelector: &meta.LabelSelector{MatchLabels: patroni.Labels},
					TopologyKey:   "kubernetes.io/hostname",
				}},
			},
		},
		InitContainers:     podInitContainers,
		Containers:         podContainers,
		HostNetwork:        true,
		NodeSelector:       patroni.Spec.CommonConfiguration.NodeSelector,
		ServiceAccountName: serviceAccountName,
		Tolerations:        patroni.Spec.CommonConfiguration.Tolerations,
		Volumes: []core.Volume{
			{
				Name: "patroni-storage-init",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: storagePath,
						Type: &initHostPathType,
					},
				},
			},
			{
				Name: patroni.Name + "-secret-certificates",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: patroni.Name + "-secret-certificates",
					},
				},
			},
			{
				Name: csrSignerCaVolumeName,
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: contrailcertificates.SignerCAConfigMapName,
						},
					},
				},
			},
		}}

	var stsSelector = meta.LabelSelector{
		MatchLabels: contraillabel.New("patroni", patroni.Name),
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: contraillabel.New("patroni", patroni.Name),
		},
		Spec: podSpec,
	}

	statefulSet.Spec = apps.StatefulSetSpec{
		Selector:    &stsSelector,
		ServiceName: service.Name,
		Replicas:    patroni.Spec.CommonConfiguration.Replicas,
		Template:    stsTemplate,
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {

		contrail.SetSTSCommonConfiguration(statefulSet, &patroni.Spec.CommonConfiguration)

		var patroniGroupId int64 = 0
		statefulSet.Spec.Template.Spec.SecurityContext = &core.PodSecurityContext{}
		statefulSet.Spec.Template.Spec.SecurityContext.FSGroup = &patroniGroupId
		statefulSet.Spec.Template.Spec.SecurityContext.RunAsGroup = &patroniGroupId
		statefulSet.Spec.Template.Spec.SecurityContext.RunAsUser = &patroniGroupId

		storageClassName := "local-storage"
		statefulSet.Spec.VolumeClaimTemplates = []core.PersistentVolumeClaim{
			{
				ObjectMeta: meta.ObjectMeta{
					Name:      "pgdata",
					Namespace: patroni.Namespace,
					Labels:    patroni.Labels,
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

		statefulSet.Spec.Selector = &meta.LabelSelector{MatchLabels: patroni.Labels}
		return controllerutil.SetControllerReference(patroni, statefulSet, r.scheme)
	})
	return statefulSet, err
}

func (r *ReconcilePatroni) ensureLocalPVsExist(patroni *contrail.Patroni) error {
	path := patroni.Spec.ServiceConfiguration.Storage.Path
	size := patroni.Spec.ServiceConfiguration.Storage.Size
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
		path = defaultPatroniStoragePath
	}

	for i := int32(0); i < patroni.Spec.CommonConfiguration.GetReplicas(); i++ {
		name := fmt.Sprintf("%v-%v-patroni-data-%v", patroni.Name, patroni.Namespace, i)
		nodeSelectors := patroni.Spec.CommonConfiguration.NodeSelector
		lv, err := r.volumes.New(name, path, storage, patroni.Labels, nodeSelectors)
		if err != nil {
			return err
		}
		if err := lv.EnsureExists(); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcilePatroni) ensurePVCOwnershipExists(patroni *contrail.Patroni) error {
	listOps := &client.ListOptions{Namespace: patroni.Namespace, LabelSelector: labels.SelectorFromSet(patroni.Labels)}
	pvcList := &core.PersistentVolumeClaimList{}
	if err := r.client.List(context.TODO(), pvcList, listOps); err != nil {
		return err
	}
	for _, pvc := range pvcList.Items {
		if err := r.kubernetes.Owner(patroni).EnsureOwns(&pvc); err != nil {
			return err
		}
	}
	return nil
}

func getImage(containers []*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"patroni": "localhost:5000/patroni",
	}
	c := utils.GetContainerFromList(containerName, containers)
	if c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

