package zookeeper

import (
	"atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
	"atom/atom/contrail/operator/pkg/controller/utils"
	"context"

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

var log = logf.Log.WithName("controller_zookeeper")
var err error

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ZookeeperList{}
			err := myclient.List(context.TODO(), listOps, list)
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
			list := &v1alpha1.ZookeeperList{}
			err := myclient.List(context.TODO(), listOps, list)
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
			list := &v1alpha1.ZookeeperList{}
			err := myclient.List(context.TODO(), listOps, list)
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
			list := &v1alpha1.ZookeeperList{}
			err := myclient.List(context.TODO(), listOps, list)
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

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileZookeeper{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller

	c, err := controller.New("zookeeper-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	if err = c.Watch(&source.Kind{Type: &v1alpha1.Zookeeper{}},
		&handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	serviceMap := map[string]string{"contrail_manager": "zookeeper"}
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
	predManagerSizeChange := utils.ManagerSizeChange(utils.ZookeeperGroupKind())
	if err = c.Watch(srcManager, managerHandler, predManagerSizeChange); err != nil {
		return err
	}

	srcSTS := &source.Kind{Type: &appsv1.StatefulSet{}}
	stsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Zookeeper{},
	}
	stsPred := utils.STSStatusChange(utils.ZookeeperGroupKind())
	if err = c.Watch(srcSTS, stsHandler, stsPred); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileZookeeper implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileZookeeper{}

// ReconcileZookeeper reconciles a Zookeeper object.
type ReconcileZookeeper struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile reconciles zookeeper.
func (r *ReconcileZookeeper) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Zookeeper")
	instanceType := "zookeeper"
	instance := &v1alpha1.Zookeeper{}
	if err = r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	currentZookeeperInstance := *instance
	managerInstance, err := instance.OwnedByManager(r.Client, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if managerInstance != nil {
		if managerInstance.Spec.Services.Zookeepers != nil {
			for _, zookeeperManagerInstance := range managerInstance.Spec.Services.Zookeepers {
				if zookeeperManagerInstance.Name == request.Name {
					instance.Spec.CommonConfiguration = utils.MergeCommonConfiguration(
						managerInstance.Spec.CommonConfiguration,
						zookeeperManagerInstance.Spec.CommonConfiguration)
					if err = r.Client.Update(context.TODO(), instance); err != nil {
						return reconcile.Result{}, nil
					}
				}
			}
		}
	}

	configMap, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	configMap2, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-1", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{configMap.Name: request.Name + "-" + instanceType + "-volume", configMap2.Name: request.Name + "-" + instanceType + "-volume-1"})

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {
		for containerName, image := range instance.Spec.ServiceConfiguration.Images {
			if containerName == container.Name {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = image
			}
			if containerName == "zookeeper" {
				command := []string{
					"bash", "-c", "grep ${POD_IP} /mydata/zoo.cfg.dynamic.100000000|" +
						"sed -e 's/.*server.\\(.*\\)=.*/\\1/' > /data/myid && " +
						"cp /conf-1/* /conf/ && " +
						"sed -i \"s/clientPortAddress=.*/clientPortAddress=${POD_IP}/g\" /conf/zoo.cfg && " +
						"zkServer.sh --config /conf start-foreground"}
				//command = []string{"sh", "-c", "while true; do echo hello; sleep 10;done"}
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
				volumeMountList := []corev1.VolumeMount{}
				volumeMount := corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume-1",
					MountPath: "/conf-1",
				}
				volumeMountList = append(volumeMountList, volumeMount)
				volumeMount = corev1.VolumeMount{
					Name:      request.Name + "-" + instanceType + "-volume",
					MountPath: "/mydata",
				}
				volumeMountList = append(volumeMountList, volumeMount)
				(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			}
		}
	}
	// Configure InitContainers.
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		for containerName, image := range instance.Spec.ServiceConfiguration.Images {
			if containerName == container.Name {
				(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = image
			}
		}
	}
	if err = instance.CreateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	strategy := "rolling"
	if currentZookeeperInstance.Spec.CommonConfiguration.Replicas != nil {
		if int(*currentZookeeperInstance.Spec.CommonConfiguration.Replicas) == 1 && int(*instance.Spec.CommonConfiguration.Replicas) > 1 {
			strategy = "deleteFirst"
		}
	}
	if err = instance.UpdateSTS(statefulSet, &instance.Spec.CommonConfiguration, instanceType, request, r.Scheme, r.Client, strategy); err != nil {
		return reconcile.Result{}, err
	}
	podIPList, podIPMap, err := instance.PodIPListAndIPMapFromInstance(instanceType, request, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPList.Items) > 0 {
		if err = instance.InstanceConfiguration(request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			//return reconcile.Result{}, err
			return reconcile.Result{Requeue: true}, nil
		}
		if err = instance.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
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
