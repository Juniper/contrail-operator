package config

import (
	"context"
	"reflect"
	"strconv"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"

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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_config")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ConfigList{}
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
			list := &v1alpha1.ConfigList{}
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
			list := &v1alpha1.ConfigList{}
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
			list := &v1alpha1.ConfigList{}
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

// Add adds the Config controller to the manager.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConfig{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Manager: mgr,
		claims:  volumeclaims.New(mgr.GetClient(), mgr.GetScheme()),
	}
}
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller.
	c, err := controller.New("config-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Config.
	if err = c.Watch(&source.Kind{Type: &v1alpha1.Config{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	serviceMap := map[string]string{"contrail_manager": "config"}
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
	predManagerSizeChange := utils.ManagerSizeChange(utils.ConfigGroupKind())
	if err = c.Watch(srcManager, managerHandler, predManagerSizeChange); err != nil {
		return err
	}

	srcCassandra := &source.Kind{Type: &v1alpha1.Cassandra{}}
	cassandraHandler := resourceHandler(mgr.GetClient())
	predCassandraSizeChange := utils.CassandraActiveChange()
	if err = c.Watch(srcCassandra, cassandraHandler, predCassandraSizeChange); err != nil {
		return err
	}

	srcRabbitmq := &source.Kind{Type: &v1alpha1.Rabbitmq{}}
	rabbitmqHandler := resourceHandler(mgr.GetClient())
	predRabbitmqSizeChange := utils.RabbitmqActiveChange()
	if err = c.Watch(srcRabbitmq, rabbitmqHandler, predRabbitmqSizeChange); err != nil {
		return err
	}

	srcZookeeper := &source.Kind{Type: &v1alpha1.Zookeeper{}}
	zookeeperHandler := resourceHandler(mgr.GetClient())
	predZookeeperSizeChange := utils.ZookeeperActiveChange()
	if err = c.Watch(srcZookeeper, zookeeperHandler, predZookeeperSizeChange); err != nil {
		return err
	}

	srcSTS := &source.Kind{Type: &appsv1.StatefulSet{}}
	stsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Config{},
	}
	stsPred := utils.STSStatusChange(utils.ConfigGroupKind())
	if err = c.Watch(srcSTS, stsHandler, stsPred); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileConfig implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileConfig{}

// ReconcileConfig reconciles a Config object.
type ReconcileConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	Client  client.Client
	Scheme  *runtime.Scheme
	Manager manager.Manager
	claims  volumeclaims.PersistentVolumeClaims
}

// Reconcile reconciles Config.
func (r *ReconcileConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Config")
	instanceType := "config"
	config := &v1alpha1.Config{}
	cassandraInstance := &v1alpha1.Cassandra{}
	zookeeperInstance := &v1alpha1.Zookeeper{}
	rabbitmqInstance := &v1alpha1.Rabbitmq{}

	if err := r.Client.Get(context.TODO(), request.NamespacedName, config); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if !config.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	cassandraActive := cassandraInstance.IsActive(config.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	zookeeperActive := zookeeperInstance.IsActive(config.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, r.Client)
	rabbitmqActive := rabbitmqInstance.IsActive(config.Labels["contrail_cluster"],
		request.Namespace, r.Client)
	cassandraUpgrading := cassandraInstance.IsUpgrading(config.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	zookeeperUpgrading := zookeeperInstance.IsUpgrading(config.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	rabbitmqUpgrading := rabbitmqInstance.IsUpgrading(config.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)

	if !cassandraActive || !rabbitmqActive || !zookeeperActive || cassandraUpgrading || zookeeperUpgrading || rabbitmqUpgrading {
		return reconcile.Result{}, nil
	}

	managerInstance, err := config.OwnedByManager(r.Client, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if managerInstance != nil {
		if managerInstance.Spec.Services.Config != nil {
			configManagerInstance := managerInstance.Spec.Services.Config
			if configManagerInstance.Name == request.Name {
				config.Spec.CommonConfiguration = utils.MergeCommonConfiguration(
					managerInstance.Spec.CommonConfiguration,
					configManagerInstance.Spec.CommonConfiguration)
				err = r.Client.Update(context.TODO(), config)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	currentConfigMap, currentConfigExists := config.CurrentConfigMapExists(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)

	configMap, err := config.CreateConfigMap(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	secretCertificates, err := config.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	// DeviceManager pushes configuration to dnsmasq service and then needs to restart it by sending a signal.
	// Therefore those services needs to share a one process namespace
	// TODO: Move device manager and dnsmasq to a separate pod. They are separate services which requires
	// persistent volumes and capabilities
	trueVal := true
	statefulSet.Spec.Template.Spec.ShareProcessNamespace = &trueVal
	if err = config.PrepareSTS(statefulSet, &config.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	configDefaultConfigurationInterface := config.ConfigurationParameters()
	configDefaultConfiguration := configDefaultConfigurationInterface.(v1alpha1.ConfigConfiguration)
	var persistentVolumeClaimList []corev1.PersistentVolumeClaim
	for storageName, storage := range configDefaultConfiguration.Storages {
		storageResource := corev1.ResourceStorage
		diskSize, err := resource.ParseQuantity(storage.Size)
		if err != nil {
			return reconcile.Result{}, err
		}
		storageClassName := "local-storage"
		pvc := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      storageName,
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
		}
		persistentVolumeClaimList = append(persistentVolumeClaimList, pvc)

	}
	statefulSet.Spec.VolumeClaimTemplates = persistentVolumeClaimList

	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	config.AddVolumesToIntendedSTS(statefulSet, map[string]string{
		configMap.Name:                          request.Name + "-" + instanceType + "-volume",
		cacertificates.CsrSignerCAConfigMapName: csrSignerCaVolumeName,
	})
	config.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	configNodeMgr := true
	analyticsNodeMgr := true

	if _, ok := config.Spec.ServiceConfiguration.Containers["nodemanagerconfig"]; !ok {
		configNodeMgr = false
	}

	if _, ok := config.Spec.ServiceConfiguration.Containers["nodemanageranalytics"]; !ok {
		analyticsNodeMgr = false
	}

	if config.Spec.ServiceConfiguration.NodeManager != nil && !*config.Spec.ServiceConfiguration.NodeManager {
		for idx, container := range statefulSet.Spec.Template.Spec.Containers {
			if container.Name == "nodemanagerconfig" {
				statefulSet.Spec.Template.Spec.Containers = utils.RemoveIndex(statefulSet.Spec.Template.Spec.Containers, idx)
			}
		}
		for idx, container := range statefulSet.Spec.Template.Spec.Containers {
			if container.Name == "nodemanageranalytics" {
				statefulSet.Spec.Template.Spec.Containers = utils.RemoveIndex(statefulSet.Spec.Template.Spec.Containers, idx)
			}
		}
	}

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

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {

		switch container.Name {
		case "api":
			command := []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini; /usr/bin/python /usr/bin/contrail-api --conf_file /etc/mycontrail/api.${POD_IP} --conf_file /etc/mycontrail/contrail-keystone-auth.conf --worker_id 0"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "devicemanager":
			deviceManagerCommand := `/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini;
/usr/bin/rm -f /etc/contrail/contrail-keystone-auth.conf; ln -s /etc/mycontrail/contrail-keystone-auth.conf /etc/contrail/contrail-keystone-auth.conf;
/usr/bin/rm -f /etc/contrail/contrail-fabric-ansible.conf; ln -s /etc/mycontrail/contrail-fabric-ansible.conf.${POD_IP} /etc/contrail/contrail-fabric-ansible.conf;
/usr/bin/python /usr/bin/contrail-device-manager --conf_file /etc/mycontrail/devicemanager.${POD_IP} --conf_file /etc/contrail/contrail-keystone-auth.conf
`
			command := []string{"bash", "-c", deviceManagerCommand}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			(&statefulSet.Spec.Template.Spec.Containers[idx]).SecurityContext = &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{"SYS_PTRACE"},
				},
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
				corev1.VolumeMount{
					Name:      "tftp",
					MountPath: "/var/lib/tftp",
				},
				corev1.VolumeMount{
					Name:      "dnsmasq",
					MountPath: "/var/lib/dnsmasq",
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "dnsmasq":
			container := &statefulSet.Spec.Template.Spec.Containers[idx]
			container.Command = []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini;ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini;" +
					"dnsmasq -k -p0 --conf-file=/etc/mycontrail/dnsmasq.${POD_IP}"}

			if config.Spec.ServiceConfiguration.Containers[container.Name].Command != nil {
				container.Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			container.SecurityContext = &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{"NET_ADMIN", "NET_RAW"},
				},
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
				corev1.VolumeMount{
					Name:      "tftp",
					MountPath: "/etc/tftp",
				},
				corev1.VolumeMount{
					Name:      "dnsmasq",
					MountPath: "/var/lib/dnsmasq",
				},
			)
			// DNSMasq container requires those variables to be set
			// TODO: Pass keystone credentials
			container.Env = append(container.Env, []corev1.EnvVar{
				{Name: "KEYSTONE_AUTH_ADMIN_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: config.Spec.ServiceConfiguration.KeystoneSecretName,
							},
							Key: "password",
						},
					},
				},
				{Name: "KEYSTONE_AUTH_ADMIN_USER", Value: "admin"},
				{Name: "KEYSTONE_AUTH_ADMIN_TENANT", Value: "admin"},
			}...)
			container.VolumeMounts = volumeMountList
			container.Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "servicemonitor":
			command := []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini; /usr/bin/python /usr/bin/contrail-svc-monitor --conf_file /etc/mycontrail/servicemonitor.${POD_IP} --conf_file /etc/mycontrail/contrail-keystone-auth.conf"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "schematransformer":
			command := []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini; /usr/bin/python /usr/bin/contrail-schema --conf_file /etc/mycontrail/schematransformer.${POD_IP}  --conf_file /etc/mycontrail/contrail-keystone-auth.conf"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "analyticsapi":
			command := []string{"bash", "-c",
				"/usr/bin/rm -f /etc/contrail/vnc_api_lib.ini; ln -s /etc/mycontrail/vnc.${POD_IP} /etc/contrail/vnc_api_lib.ini; /usr/bin/python /usr/bin/contrail-analytics-api -c /etc/mycontrail/analyticsapi.${POD_IP} -c /etc/mycontrail/contrail-keystone-auth.conf"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "queryengine":
			queryEngineContainer := &statefulSet.Spec.Template.Spec.Containers[idx]
			queryEngineContainer.Command = []string{"bash", "-c",
				"/usr/bin/contrail-query-engine --conf_file /etc/mycontrail/queryengine.${POD_IP}"}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			queryEngineContainer.VolumeMounts = volumeMountList
			queryEngineContainer.Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "collector":
			command := []string{"bash", "-c",
				"/usr/bin/contrail-collector --conf_file /etc/mycontrail/collector.${POD_IP}"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
				corev1.VolumeMount{
					Name:      request.Name + "-secret-certificates",
					MountPath: "/etc/certificates",
				},
				corev1.VolumeMount{
					Name:      csrSignerCaVolumeName,
					MountPath: cacertificates.CsrSignerCAMountPath,
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "redis":
			command := []string{"bash", "-c",
				"redis-server --lua-time-limit 15000 --dbfilename '' --bind 127.0.0.1 ${POD_IP} --port 6379"}
			if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
			volumeMountList = append(volumeMountList,
				corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/etc/mycontrail",
				},
			)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		case "nodemanagerconfig":
			if configNodeMgr {
				command := []string{"bash", "-c",
					"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanagerconfig.${POD_IP} > /etc/contrail/contrail-config-nodemgr.conf; /usr/bin/python /usr/bin/contrail-nodemgr --nodetype=contrail-config"}

				if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
					(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
				} else {
					(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
				}
				volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
				volumeMountList = append(volumeMountList,
					corev1.VolumeMount{
						Name:      request.Name + "-" + instanceType + "-volume",
						MountPath: "/etc/mycontrail",
					},
				)
				(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
			}
		case "nodemanageranalytics":
			if analyticsNodeMgr {
				command := []string{"bash", "-c",
					"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanageranalytics.${POD_IP} > /etc/contrail/contrail-analytics-nodemgr.conf;/usr/bin/python /usr/bin/contrail-nodemgr --nodetype=contrail-analytics"}

				if config.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
					(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
				} else {
					(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
				}
				volumeMountList := statefulSet.Spec.Template.Spec.Containers[idx].VolumeMounts
				volumeMountList = append(volumeMountList,
					corev1.VolumeMount{
						Name:      request.Name + "-" + instanceType + "-volume",
						MountPath: "/etc/mycontrail",
					},
				)
				(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
			}
		}
	}

	if !configNodeMgr {
		for idx, container := range statefulSet.Spec.Template.Spec.Containers {
			if container.Name == "nodemanagerconfig" {
				statefulSet.Spec.Template.Spec.Containers = utils.RemoveIndex(statefulSet.Spec.Template.Spec.Containers, idx)
			}
		}
	}

	if !analyticsNodeMgr {
		for idx, container := range statefulSet.Spec.Template.Spec.Containers {
			if container.Name == "nodemanageranalytics" {
				statefulSet.Spec.Template.Spec.Containers = utils.RemoveIndex(statefulSet.Spec.Template.Spec.Containers, idx)
			}
		}
	}

	var initVolumeList []corev1.Volume
	for storageName, storage := range configDefaultConfiguration.Storages {
		initHostPathType := corev1.HostPathType("DirectoryOrCreate")
		initHostPathSource := &corev1.HostPathVolumeSource{
			Path: storage.Path,
			Type: &initHostPathType,
		}
		initVolume := corev1.Volume{
			Name: request.Name + "-" + instanceType + "-" + storageName + "-init",
			VolumeSource: corev1.VolumeSource{
				HostPath: initHostPathSource,
			},
		}
		initVolumeList = append(initVolumeList, initVolume)
		statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, initVolume)
	}

	// Configure InitContainers
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = config.Spec.ServiceConfiguration.Containers[container.Name].Image
		if config.Spec.ServiceConfiguration.Containers[container.Name].Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = config.Spec.ServiceConfiguration.Containers[container.Name].Command
		}
		for _, volume := range initVolumeList {
			volumeMount := corev1.VolumeMount{
				Name:      volume.Name,
				MountPath: volume.VolumeSource.HostPath.Path,
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
		}
	}

	volumeBindingMode := storagev1.VolumeBindingMode("WaitForFirstConsumer")
	storageClass := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-storage",
			Namespace: config.Namespace,
		},
		Provisioner:       "kubernetes.io/no-provisioner",
		VolumeBindingMode: &volumeBindingMode,
	}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: storageClass.Name}, storageClass)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(config, storageClass, r.Scheme); err != nil {
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
	for k, v := range config.Spec.CommonConfiguration.NodeSelector {
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

	storageCount := 1
	for _, storage := range configDefaultConfiguration.Storages {
		localVolumeSource := corev1.LocalVolumeSource{
			Path: storage.Path,
		}
		diskSize, err := resource.ParseQuantity(storage.Size)
		if err != nil {
			return reconcile.Result{}, err
		}
		replicasInt := int(*config.Spec.CommonConfiguration.Replicas)
		for i := 0; i < replicasInt; i++ {
			pv := &corev1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name:      config.Name + "-pv-" + strconv.Itoa(i) + "-" + strconv.Itoa(storageCount),
					Namespace: config.Namespace,
				},
				Spec: corev1.PersistentVolumeSpec{
					Capacity:   corev1.ResourceList{corev1.ResourceStorage: diskSize},
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
				if err = controllerutil.SetControllerReference(config, pv, r.Scheme); err != nil {
					return reconcile.Result{}, err
				}
				if err = r.Client.Create(context.TODO(), pv); err != nil && !errors.IsAlreadyExists(err) {
					return reconcile.Result{}, err
				}
			}
		}
		storageCount++
	}

	configChanged := false
	if config.Status.ConfigChanged != nil {
		configChanged = *config.Status.ConfigChanged
	}

	if err = config.CreateSTS(statefulSet, &config.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	strategy := "deleteFirst"
	if err = config.UpdateSTS(statefulSet, &config.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client, strategy); err != nil {
		return reconcile.Result{}, err
	}

	podIPList, podIPMap, err := config.PodIPListAndIPMapFromInstance(request, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPMap) > 0 {
		if err = config.InstanceConfiguration(request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		hostNetwork := true
		if config.Spec.CommonConfiguration.HostNetwork != nil {
			hostNetwork = *config.Spec.CommonConfiguration.HostNetwork
		}
		if err = certificates.CreateAndSignCsr(r.Client, request, r.Scheme, config, r.Manager.GetConfig(), podIPList, hostNetwork); err != nil {
			return reconcile.Result{}, err
		}

		if err = config.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = config.WaitForPeerPods(request, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = config.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": instanceType, instanceType: request.Name})
		listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
		pvcList := &corev1.PersistentVolumeClaimList{}
		err = r.Client.List(context.TODO(), pvcList, listOps)
		if err != nil {
			return reconcile.Result{}, err
		}
		for _, pvc := range pvcList.Items {
			if err = controllerutil.SetControllerReference(config, &pvc, r.Scheme); err != nil {
				return reconcile.Result{}, err
			}
			if err = r.Client.Update(context.TODO(), &pvc); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	if currentConfigExists {
		newConfigMap := &corev1.ConfigMap{}
		_ = r.Client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-configmap", Namespace: request.Namespace}, newConfigMap)
		if !reflect.DeepEqual(currentConfigMap.Data, newConfigMap.Data) {
			configChanged = true
		} else {
			configChanged = false
		}
		config.Status.ConfigChanged = &configChanged
	}

	if config.Status.Active == nil {
		active := false
		config.Status.Active = &active
	}
	if err = config.SetInstanceActive(r.Client, config.Status.Active, statefulSet, request); err != nil {
		return reconcile.Result{}, err
	}
	if config.Status.ConfigChanged != nil {
		if *config.Status.ConfigChanged {
			return reconcile.Result{Requeue: true}, nil
		}
	}
	return reconcile.Result{}, nil
}
