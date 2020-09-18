package zookeeper

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/label"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
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
	policy "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_zookeeper")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ZookeeperList{}
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
			list := &v1alpha1.ZookeeperList{}
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
			list := &v1alpha1.ZookeeperList{}
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
			list := &v1alpha1.ZookeeperList{}
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
	if err := r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	currentZookeeperInstance := *instance

	_, err := instance.CreateConfigMap("correct-zookeeper-conf", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{"correct-zookeeper-conf": "correct-zookeeper-conf"})

	zookeeperDefaultConfigurationInterface := instance.ConfigurationParameters()
	zookeeperDefaultConfiguration := zookeeperDefaultConfigurationInterface.(v1alpha1.ZookeeperConfiguration)

	storageResource := corev1.ResourceStorage
	diskSize, err := resource.ParseQuantity(zookeeperDefaultConfiguration.Storage.Size)
	if err != nil {
		return reconcile.Result{}, err
	}
	storageClassName := "local-storage"
	statefulSet.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pvc",
			Namespace: request.Namespace,
			Labels:    label.New(instanceType, request.Name),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: label.New(instanceType, request.Name),
			},
			StorageClassName: &storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{storageResource: diskSize},
			},
		},
	}}
	for idx, container := range statefulSet.Spec.Template.Spec.Containers {

		if container.Name == "zookeeper" {
			command := []string{
				"bash", "-c", "zkServer.sh --config /var/lib/zookeeper start-foreground"}
			//"zkServer.sh --config /var/lib/zookeeper start-foreground"
			instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
			if instanceContainer.Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instanceContainer.Command
			}
			volumeMountList := []corev1.VolumeMount{}
			volumeMount := corev1.VolumeMount{
				Name:      "pvc",
				MountPath: "/var/lib/zookeeper",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instanceContainer.Image

		}

	}
	initHostPathType := corev1.HostPathType("DirectoryOrCreate")
	initHostPathSource := &corev1.HostPathVolumeSource{
		Path: zookeeperDefaultConfiguration.Storage.Path,
		Type: &initHostPathType,
	}
	initVolume := corev1.Volume{
		Name: request.Name + "-" + instanceType + "-init",
		VolumeSource: corev1.VolumeSource{
			HostPath: initHostPathSource,
		},
	}
	statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, initVolume)
	statefulSet.Spec.PodManagementPolicy = appsv1.OrderedReadyPodManagement
	// Configure InitContainers.
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
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instanceContainer.Image
		if instanceContainer.Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instanceContainer.Command
		}
		if container.Name == "init" {
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-init",
				MountPath: zookeeperDefaultConfiguration.Storage.Path,
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
		}
		if container.Name == "conf-init" {
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-init",
				MountPath: "/mnt/zookeeper",
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      "correct-zookeeper-conf",
				MountPath: "/zookeeper-conf",
			}
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts = append((&statefulSet.Spec.Template.Spec.InitContainers[idx]).VolumeMounts, volumeMount)
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = []string{
				"sh", "-c", "if [ ! -f /mnt/zookeeper/zoo.cfg ]; then cp /zookeeper-conf/* /mnt/zookeeper/ && cp /mnt/zookeeper/zoo.cfg.$POD_IP /mnt/zookeeper/zoo.cfg && cp /mnt/zookeeper/myid.$POD_IP /mnt/zookeeper/myid; fi;"}
		}
	}

	volumeBindingMode := storagev1.VolumeBindingMode("WaitForFirstConsumer")
	storageClass := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "local-storage",
		},
		Provisioner:       "kubernetes.io/no-provisioner",
		VolumeBindingMode: &volumeBindingMode,
	}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: storageClass.Name}, storageClass)
	if err != nil && errors.IsNotFound(err) {
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
		Path: zookeeperDefaultConfiguration.Storage.Path,
	}

	replicasInt := int(*instance.Spec.CommonConfiguration.Replicas)
	for i := 0; i < replicasInt; i++ {
		pv := &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   instance.Name + "-pv-" + strconv.Itoa(i),
				Labels: label.New(instanceType, request.Name),
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
			if err = r.Client.Create(context.TODO(), pv); err != nil && !errors.IsAlreadyExists(err) {
				return reconcile.Result{}, err
			}
		}
	}

	if err = instance.CreateSTS(statefulSet, instanceType, request, r.Client); err != nil {
		return reconcile.Result{}, err
	}
	strategy := "rolling"
	if currentZookeeperInstance.Spec.CommonConfiguration.Replicas != nil {
		if int(*currentZookeeperInstance.Spec.CommonConfiguration.Replicas) == 1 && int(*instance.Spec.CommonConfiguration.Replicas) > 1 {
			strategy = "deleteFirst"
		}
	}
	if err = instance.UpdateSTS(statefulSet, instanceType, request, r.Client, strategy); err != nil {
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

		pods := make([]corev1.Pod, len(podIPList.Items))
		copy(pods, podIPList.Items)
		sort.SliceStable(pods, func(i, j int) bool { return pods[i].Name < pods[j].Name })

		var found *corev1.Pod
		for _, pod := range pods {
			ip, ok := instance.Status.Nodes[pod.Name]
			if !ok || ip != pod.Status.PodIP {
				for _, c := range pod.Status.Conditions {
					if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
						found = &pod
						break
					}
				}
			}
		}

		if found != nil && len(pods) > 1 {
			myidString := found.Name[len(found.Name)-1:]
			myidInt, err := strconv.Atoi(myidString)
			if err != nil {
				return reconcile.Result{}, err
			}
			serverDef := fmt.Sprintf("server.%d=%s:%s;%s:2181",
				myidInt+1, found.Status.PodIP,
				strconv.Itoa(*zookeeperDefaultConfiguration.ElectionPort)+":"+strconv.Itoa(*zookeeperDefaultConfiguration.ServerPort), found.Status.PodIP)
			runScript := fmt.Sprintf("zkCli.sh -server %s reconfig -add \"%s\"", found.Status.PodIP, serverDef)
			reqLogger.Info(serverDef)
			command := []string{"bash", "-c", runScript, serverDef}
			_, _, err = v1alpha1.ExecToPodThroughAPI(command, "zookeeper", found.Name, found.Namespace, nil)
			if err != nil {
				return reconcile.Result{}, err
			}
			instance.Status.Reconfigs = instance.Status.Reconfigs + 1
		}

		for _, pod := range podIPList.Items {
			for _, c := range pod.Status.Conditions {
				if c.Type == corev1.PodReady && c.Status == corev1.ConditionFalse {
					delete(podIPMap, pod.Name)
				}
			}
		}

		if err = instance.ManageNodeStatus(podIPMap, r.Client); err != nil {
			return reconcile.Result{}, err
		}

		labelSelector := labels.SelectorFromSet(label.New(instanceType, request.Name))
		listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
		pvcList := &corev1.PersistentVolumeClaimList{}
		err = r.Client.List(context.TODO(), pvcList, listOps)
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

	if err = r.ensurePodDisruptionBudgetExists(instance); err != nil {
		return reconcile.Result{}, err
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
func (r *ReconcileZookeeper) ensurePodDisruptionBudgetExists(zookeeper *v1alpha1.Zookeeper) error {
	pdb := &policy.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      zookeeper.Name + "-zookeeper",
			Namespace: zookeeper.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.Client, pdb, func() error {
		oneVal := intstr.FromInt(1)
		pdb.ObjectMeta.Labels = label.New("zookeeper", zookeeper.Name)
		pdb.Spec.MaxUnavailable = &oneVal
		pdb.Spec.Selector = metav1.SetAsLabelSelector(label.New("zookeeper", zookeeper.Name))
		return controllerutil.SetControllerReference(zookeeper, pdb, r.Scheme)
	})

	return err
}
