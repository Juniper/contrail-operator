package config

import (
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"

	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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

	return &ReconcileConfig{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
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
	Client    client.Client
	Scheme    *runtime.Scheme
	podsReady *bool
}

// Reconcile reconciles Config.
func (r *ReconcileConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var err error
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Config")
	instanceType := "config"
	instance := &v1alpha1.Config{}
	cassandraInstance := &v1alpha1.Cassandra{}
	zookeeperInstance := &v1alpha1.Zookeeper{}
	rabbitmqInstance := &v1alpha1.Rabbitmq{}

	if err = r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	//currentInstance := *instance

	cassandraActive := cassandraInstance.IsActive(instance.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	zookeeperActive := zookeeperInstance.IsActive(instance.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, r.Client)
	rabbitmqActive := rabbitmqInstance.IsActive(instance.Labels["contrail_cluster"],
		request.Namespace, r.Client)
	cassandraUpgrading := cassandraInstance.IsUpgrading(instance.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	zookeeperUpgrading := zookeeperInstance.IsUpgrading(instance.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)
	rabbitmqUpgrading := rabbitmqInstance.IsUpgrading(instance.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.Client)

	if !cassandraActive || !rabbitmqActive || !zookeeperActive || cassandraUpgrading || zookeeperUpgrading || rabbitmqUpgrading {
		return reconcile.Result{}, nil
	}

	managerInstance, err := instance.OwnedByManager(r.Client, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if managerInstance != nil {
		if managerInstance.Spec.Services.Config != nil {
			configManagerInstance := managerInstance.Spec.Services.Config
			if configManagerInstance.Name == request.Name {
				instance.Spec.CommonConfiguration = utils.MergeCommonConfiguration(
					managerInstance.Spec.CommonConfiguration,
					configManagerInstance.Spec.CommonConfiguration)
				err = r.Client.Update(context.TODO(), instance)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	currentConfigMap, currentConfigExists := instance.CurrentConfigMapExists(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)

	configMap, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{configMap.Name: request.Name + "-" + instanceType + "-volume"})

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {
		if container.Name == "api" {
			command := []string{"bash", "-c",
				"/usr/bin/python /usr/bin/contrail-api --conf_file /etc/mycontrail/api.${POD_IP} --conf_file /etc/contrail/contrail-keystone-auth.conf --worker_id 0"}

			//command := []string{"python", "/app/contrail/configapi/config_api_debug.binary", "--conf_file /etc/mycontrail/api.${POD_IP}", "--worker_id", "0"}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}

			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "devicemanager" {
			command := []string{"bash", "-c",
				"/usr/bin/python /usr/bin/contrail-device-manager --conf_file /etc/mycontrail/devicemanager.${POD_IP} --conf_file /etc/contrail/contrail-keystone-auth.conf"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "servicemonitor" {
			command := []string{"bash", "-c",
				"/usr/bin/python /usr/bin/contrail-svc-monitor --conf_file /etc/mycontrail/servicemonitor.${POD_IP} --conf_file /etc/contrail/contrail-keystone-auth.conf"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "schematransformer" {
			command := []string{"bash", "-c",
				"/usr/bin/python /usr/bin/contrail-schema --conf_file /etc/mycontrail/schematransformer.${POD_IP}  --conf_file /etc/contrail/contrail-keystone-auth.conf"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "analyticsapi" {
			command := []string{"bash", "-c",
				"/usr/bin/python /usr/bin/contrail-analytics-api -c /etc/mycontrail/analyticsapi.${POD_IP} -c /etc/contrail/contrail-keystone-auth.conf"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "collector" {
			command := []string{"bash", "-c",
				"/usr/bin/contrail-collector --conf_file /etc/mycontrail/collector.${POD_IP}"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "redis" {
			command := []string{"bash", "-c",
				"redis-server --lua-time-limit 15000 --dbfilename '' --bind 127.0.0.1 ${POD_IP} --port 6379"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "nodemanagerconfig" {

			command := []string{"bash", "-c",
				"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanagerconfig.${POD_IP} > /etc/contrail/contrail-config-nodemgr.conf; /usr/bin/python /usr/bin/contrail-nodemgr --nodetype=contrail-config"}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.NodeManager != nil && !*instance.Spec.ServiceConfiguration.NodeManager {
				command = []string{"bash", "-c",
					"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanagerconfig.${POD_IP} > /etc/contrail/contrail-config-nodemgr.conf; while true; do sleep 10; done"}
			}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}

			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
		if container.Name == "nodemanageranalytics" {
			command := []string{"bash", "-c",
				"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanageranalytics.${POD_IP} > /etc/contrail/contrail-analytics-nodemgr.conf;/usr/bin/python /usr/bin/contrail-nodemgr --nodetype=contrail-analytics"}
			if instance.Spec.ServiceConfiguration.NodeManager != nil && !*instance.Spec.ServiceConfiguration.NodeManager {
				command = []string{"bash", "-c",
					"sed \"s/hostip=.*/hostip=${POD_IP}/g\" /etc/mycontrail/nodemanageranalytics.${POD_IP} > /etc/contrail/contrail-analytics-nodemgr.conf;while true; do sleep 10; done"}
			}
			//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
			if instance.Spec.ServiceConfiguration.Containers[container.Name].Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-volume",
				MountPath: "/etc/mycontrail",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			//(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Images[container.Name]
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		}
	}

	// Configure InitContainers
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instance.Spec.ServiceConfiguration.Containers[container.Name].Image
		if instance.Spec.ServiceConfiguration.Containers[container.Name].Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instance.Spec.ServiceConfiguration.Containers[container.Name].Command
		}
		//fmt.Println(instance.Spec.ServiceConfiguration.Containers[container.Name].Image)
		/*
			for containerName, image := range instance.Spec.ServiceConfiguration.Images {
				if containerName == container.Name {
					(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = image
				}
			}
		*/

	}

	configChanged := false
	if instance.Status.ConfigChanged != nil {
		configChanged = *instance.Status.ConfigChanged
	}

	if err = instance.CreateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	strategy := "deleteFirst"
	if err = instance.UpdateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client, strategy); err != nil {
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

		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.WaitForPeerPods(request, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
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
		instance.Status.ConfigChanged = &configChanged
	}

	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err = instance.SetInstanceActive(r.Client, instance.Status.Active, statefulSet, request); err != nil {
		return reconcile.Result{}, err
	}
	if instance.Status.ConfigChanged != nil {
		if *instance.Status.ConfigChanged {
			return reconcile.Result{Requeue: true}, nil
		}
	}
	return reconcile.Result{}, nil
}
