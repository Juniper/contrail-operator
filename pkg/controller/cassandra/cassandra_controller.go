package cassandra

import (
	"context"
	"strconv"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
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
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_cassandra")
var err error

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.CassandraList{}
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
			list := &v1alpha1.CassandraList{}
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
			list := &v1alpha1.CassandraList{}
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
			list := &v1alpha1.CassandraList{}
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

// Add adds Cassandra controller to the manager.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCassandra{Client: mgr.GetClient(), Scheme: mgr.GetScheme(), Manager: mgr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller.

	c, err := controller.New("cassandra-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource Cassandra.
	if err = c.Watch(&source.Kind{Type: &v1alpha1.Cassandra{}},
		&handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to PODs.
	serviceMap := map[string]string{"contrail_manager": "cassandra"}
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
	predManagerSizeChange := utils.ManagerSizeChange(utils.CassandraGroupKind())
	if err = c.Watch(srcManager, managerHandler, predManagerSizeChange); err != nil {
		return err
	}

	srcSTS := &source.Kind{Type: &appsv1.StatefulSet{}}
	stsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Cassandra{},
	}
	stsPred := utils.STSStatusChange(utils.CassandraGroupKind())
	if err = c.Watch(srcSTS, stsHandler, stsPred); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCassandra implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileCassandra{}

// ReconcileCassandra reconciles a Cassandra object.
type ReconcileCassandra struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	Client  client.Client
	Scheme  *runtime.Scheme
	Manager manager.Manager
}

// Reconcile reconciles cassandra.
func (r *ReconcileCassandra) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Cassandra")
	instanceType := "cassandra"
	instance := &v1alpha1.Cassandra{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	// if not found we expect it a change in replicaset
	// and get the cassandra instance via label.
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	managerInstance, err := instance.OwnedByManager(r.Client, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if managerInstance != nil {
		if managerInstance.Spec.Services.Cassandras != nil {
			for _, cassandraManagerInstance := range managerInstance.Spec.Services.Cassandras {
				if cassandraManagerInstance.Name == request.Name {
					instance.Spec.CommonConfiguration = utils.MergeCommonConfiguration(
						managerInstance.Spec.CommonConfiguration,
						cassandraManagerInstance.Spec.CommonConfiguration)
					if err = r.Client.Update(context.TODO(), instance); err != nil {
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

	secretCertificates, err := instance.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{configMap.Name: request.Name + "-" + instanceType + "-volume"})
	instance.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	cassandraDefaultConfigurationInterface := instance.ConfigurationParameters()
	cassandraDefaultConfiguration := cassandraDefaultConfigurationInterface.(v1alpha1.CassandraConfiguration)

	storageResource := corev1.ResourceStorage
	diskSize, err := resource.ParseQuantity(cassandraDefaultConfiguration.Storage.Size)
	storageClassName := "local-storage"
	statefulSet.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pvc",
			Namespace: request.Namespace,
			Labels:    map[string]string{"contrail_manager": instanceType, instanceType: request.Name},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			StorageClassName: &storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{storageResource: diskSize},
			},
		},
	}}

	emptyVolume := corev1.Volume{
		Name: request.Name + "-keystore",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, emptyVolume)

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {

		if container.Name == "cassandra" {
			command := []string{"bash", "-c",
				"/docker-entrypoint.sh cassandra -f -Dcassandra.config=file:///mydata/${POD_IP}.yaml"}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/mydata",
			}
			volumeMountList = append(volumeMountList, volumeMount)

			volumeMount = corev1.VolumeMount{
				Name:      "pvc",
				MountPath: "/var/lib/cassandra",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-keystore",
				MountPath: "/etc/keystore",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			var jvmOpts string
			if instance.Spec.ServiceConfiguration.MinHeapSize != "" {
				jvmOpts = "-Xms" + instance.Spec.ServiceConfiguration.MinHeapSize
			}
			if instance.Spec.ServiceConfiguration.MaxHeapSize != "" {
				jvmOpts = jvmOpts + " -Xmx" + instance.Spec.ServiceConfiguration.MaxHeapSize
			}
			if jvmOpts != "" {
				envs := (&statefulSet.Spec.Template.Spec.Containers[idx]).Env
				envs = append(envs)
				jvmOptEnvVar := corev1.EnvVar{
					Name:  "JVM_OPTS",
					Value: jvmOpts,
				}
				envVars := statefulSet.Spec.Template.Spec.Containers[idx].Env
				envVars = append(envVars, jvmOptEnvVar)
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Env = envVars
			}

		}

	}
	initHostPathType := corev1.HostPathType("DirectoryOrCreate")
	initHostPathSource := &corev1.HostPathVolumeSource{
		Path: cassandraDefaultConfiguration.Storage.Path,
		Type: &initHostPathType,
	}
	initVolume := corev1.Volume{
		Name: request.Name + "-" + instanceType + "-init",
		VolumeSource: corev1.VolumeSource{
			HostPath: initHostPathSource,
		},
	}

	secret, err := instance.CreateSecret(request.Name+"-secret", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	_, KPok := secret.Data["keystorePassword"]
	_, TPok := secret.Data["truststorePassword"]
	if !KPok || !TPok {
		cassandraKeystorePassword := v1alpha1.RandomString(10)
		cassandraTruststorePassword := v1alpha1.RandomString(10)
		secretData := map[string][]byte{"keystorePassword": []byte(cassandraKeystorePassword), "truststorePassword": []byte(cassandraTruststorePassword)}
		secret.Data = secretData
	}
	r.Client.Update(context.TODO(), secret)

	cassandraKeystorePassword := string(secret.Data["keystorePassword"])
	cassandraTruststorePassword := string(secret.Data["truststorePassword"])

	statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, initVolume)
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		if container.Name == "init" {
			command := []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"}

			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image

			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-init",
				MountPath: cassandraDefaultConfiguration.Storage.Path,
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
		}
		if container.Name == "init2" {
			command := []string{"bash", "-c",
				"keytool -keystore /etc/keystore/server-truststore.jks -keypass " + cassandraKeystorePassword + " -storepass " + cassandraTruststorePassword + " -noprompt -alias CARoot -import -file /run/secrets/kubernetes.io/serviceaccount/ca.crt;" +
					"openssl pkcs12 -export -in /etc/certificates/server-${POD_IP}.crt -inkey /etc/certificates/server-key-${POD_IP}.pem -chain -CAfile /run/secrets/kubernetes.io/serviceaccount/ca.crt -password pass:" + cassandraTruststorePassword + " -name $(hostname -f) -out TmpFile;" +
					"keytool -importkeystore -deststorepass " + cassandraKeystorePassword + " -destkeypass " + cassandraKeystorePassword + " -destkeystore /etc/keystore/server-keystore.jks -deststoretype pkcs12 -srcstorepass " + cassandraTruststorePassword + " -srckeystore TmpFile -srcstoretype PKCS12 -alias $(hostname -f)"}
			//"keytool -importkeystore -deststorepass " + cassandraKeystorePassword + " -destkeypass " + cassandraKeystorePassword + " -destkeystore /etc/keystore/server-keystore.jks -deststoretype pkcs12 -srcstorepass " + cassandraTruststorePassword + " -srckeystore TmpFile -srcstoretype PKCS12 -alias $(hostname -f);while true;do sleep 10;done"}

			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image

			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-keystore",
				MountPath: "/etc/keystore",
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
		}

	}

	volumeBindingMode := storagev1.VolumeBindingMode("WaitForFirstConsumer")
	storageClass := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-storage",
			Namespace: instance.Namespace,
		},
		Provisioner:       "kubernetes.io/no-provisioner",
		VolumeBindingMode: &volumeBindingMode,
	}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: storageClass.Name}, storageClass)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(instance, storageClass, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.Client.Create(context.TODO(), storageClass)
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				return reconcile.Result{}, err
			}
		}
	}

	volumeMode := corev1.PersistentVolumeMode("Filesystem")
	nodeSelectorMatchExpressions := []corev1.NodeSelectorRequirement{}
	for k, v := range instance.Spec.CommonConfiguration.NodeSelector {
		valueList := []string{v}
		expression := corev1.NodeSelectorRequirement{
			Key:      k,
			Operator: corev1.NodeSelectorOperator("In"),
			Values:   valueList,
		}
		nodeSelectorMatchExpressions = append(nodeSelectorMatchExpressions, expression)
	}
	nodeSelectorTerm := corev1.NodeSelector{
		NodeSelectorTerms: []corev1.NodeSelectorTerm{{
			MatchExpressions: nodeSelectorMatchExpressions,
		}},
	}
	volumeNodeAffinity := corev1.VolumeNodeAffinity{
		Required: &nodeSelectorTerm,
	}
	if err != nil {
		return reconcile.Result{}, err
	}
	localVolumeSource := corev1.LocalVolumeSource{
		Path: cassandraDefaultConfiguration.Storage.Path,
	}

	replicasInt := int(*instance.Spec.CommonConfiguration.Replicas)
	for i := 0; i < replicasInt; i++ {
		pv := &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.Name + "-pv-" + strconv.Itoa(i),
				Namespace: instance.Namespace,
			},
			Spec: corev1.PersistentVolumeSpec{
				Capacity:   corev1.ResourceList{storageResource: diskSize},
				VolumeMode: &volumeMode,
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimPolicy("Delete"),
				StorageClassName:              "local-storage",
				NodeAffinity:                  &volumeNodeAffinity,
				PersistentVolumeSource: corev1.PersistentVolumeSource{
					Local: &localVolumeSource,
				},
			},
		}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pv.Name, Namespace: request.Namespace}, pv)
		if err != nil && errors.IsNotFound(err) {
			if err = controllerutil.SetControllerReference(instance, pv, r.Scheme); err != nil {
				return reconcile.Result{}, err
			}
			if err = r.Client.Create(context.TODO(), pv); err != nil && !errors.IsAlreadyExists(err) {
				return reconcile.Result{}, err
			}
		}
	}

	if err = instance.CreateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	if err = instance.UpdateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client, "rolling"); err != nil {
		return reconcile.Result{}, err
	}

	podIPList, podIPMap, err := instance.PodIPListAndIPMapFromInstance(instanceType, request, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPList.Items) > 0 {
		if err = instance.InstanceConfiguration(request,
			podIPList,
			r.Client); err != nil {
			return reconcile.Result{}, err
		}
		hostNetwork := true
		if instance.Spec.CommonConfiguration.HostNetwork != nil {
			hostNetwork = *instance.Spec.CommonConfiguration.HostNetwork
		}
		if err = v1alpha1.CreateAndSignCsr(r.Client, request, r.Scheme, instance, r.Manager.GetConfig(), podIPList, hostNetwork); err != nil {
			return reconcile.Result{}, err
		}
		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		if err = instance.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": instanceType, instanceType: request.Name})
		listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
		pvcList := &corev1.PersistentVolumeClaimList{}
		err := r.Client.List(context.TODO(), pvcList, listOps)
		if err != nil {
			return reconcile.Result{}, err
		}
		for _, pvc := range pvcList.Items {
			if err = controllerutil.SetControllerReference(instance, &pvc, r.Scheme); err != nil {
				return reconcile.Result{}, err
			}
			if err = r.Client.Update(context.TODO(), &pvc); err != nil {
				return reconcile.Result{}, err
			}
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
