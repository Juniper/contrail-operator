package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_manager")

var resourcesList = []runtime.Object{
	&v1alpha1.Cassandra{},
	&v1alpha1.Zookeeper{},
	&v1alpha1.Webui{},
	&v1alpha1.ProvisionManager{},
	&v1alpha1.Config{},
	&v1alpha1.Control{},
	&v1alpha1.Rabbitmq{},
	&v1alpha1.Postgres{},
	&v1alpha1.Command{},
	&v1alpha1.Keystone{},
	&v1alpha1.Swift{},
	&v1alpha1.Memcached{},
	&v1alpha1.Vrouter{},
	&v1alpha1.Kubemanager{},
	&v1alpha1.Contrailmonitor{},
	&v1alpha1.ContrailCNI{},
	&corev1.ConfigMap{},
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Manager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	if err := apiextensionsv1beta1.AddToScheme(scheme.Scheme); err != nil {
		return err
	}
	var r reconcile.Reconciler
	reconcileManager := ReconcileManager{client: mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		manager:    mgr,
		cache:      mgr.GetCache(),
		kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme()),
	}
	r = &reconcileManager
	//r := newReconciler(mgr)
	c, err := createController(mgr, r)
	if err != nil {
		return err
	}
	reconcileManager.controller = c
	return addManagerWatch(c, mgr)
}

func createController(mgr manager.Manager, r reconcile.Reconciler) (controller.Controller, error) {
	c, err := controller.New("manager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return c, err
	}
	return c, nil
}

func addResourcesToWatch(c controller.Controller, obj runtime.Object) error {
	return c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Manager{},
	})
}

func addManagerWatch(c controller.Controller, mgr manager.Manager) error {
	err := c.Watch(&source.Kind{Type: &v1alpha1.Manager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	for _, resource := range resourcesList {
		if err = addResourcesToWatch(c, resource); err != nil {
			return err
		}
	}
	return c.Watch(&source.Kind{Type: &corev1.Node{}}, nodeChangeHandler(mgr.GetClient()))
}

// blank assignment to verify that ReconcileManager implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileManager{}

// ReconcileManager reconciles a Manager object.
type ReconcileManager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	client     client.Client
	scheme     *runtime.Scheme
	manager    manager.Manager
	controller controller.Controller
	cache      cache.Cache
	kubernetes *k8s.Kubernetes
}

// Reconcile reconciles the manager.
func (r *ReconcileManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Manager")
	instance := &v1alpha1.Manager{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	provisionConfigMap := &corev1.ConfigMap{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: "provision-config", Namespace: request.Namespace}, provisionConfigMap); err != nil {
		if errors.IsNotFound(err) {
			provisionConfigMap.Name = "provision-config"
			provisionConfigMap.Namespace = request.Namespace
			data := map[string]string{"apiserver.yaml": "", "confignodes.yaml": "", "controlnodes.yaml": "", "analyticsnodes.yaml": "", "vrouternodes.yaml": ""}
			provisionConfigMap.Data = data
			if err = controllerutil.SetControllerReference(instance, provisionConfigMap, r.scheme); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	if err := r.processCSRSignerCaConfigMap(instance); err != nil {
		return reconcile.Result{}, err
	}

	if instance.Spec.KeystoneSecretName == "" {
		instance.Spec.KeystoneSecretName = instance.Name + "-admin-password"
	}
	adminPasswordSecretName := instance.Spec.KeystoneSecretName
	if err = r.secret(adminPasswordSecretName, "manager", instance).ensureAdminPassSecretExist(); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.client.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}
	nodes, err := r.getNodes(labels.SelectorFromSet(instance.Spec.CommonConfiguration.NodeSelector))
	if err != nil {
		return reconcile.Result{}, err
	}

	replicas := r.getReplicas(nodes)
	if replicas == 0 {
		return reconcile.Result{}, nil
	}

	nodesHostAliases := r.getNodesHostAliases(nodes)
	if err := r.processCassandras(instance, replicas, nodesHostAliases); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processZookeepers(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processWebui(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processProvisionManager(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processConfig(instance, replicas, nodesHostAliases); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processKubemanagers(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processControls(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processRabbitMQ(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processVRouters(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.processContrailCNIs(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processPostgres(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processCommand(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processKeystone(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processSwift(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processMemcached(instance, replicas); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processContrailmonitor(instance); err != nil {
		return reconcile.Result{}, err
	}

	r.setConditions(instance)

	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileManager) setConditions(manager *v1alpha1.Manager) {
	readyStatus := v1alpha1.ConditionFalse
	if manager.IsClusterReady() {
		readyStatus = v1alpha1.ConditionTrue
	}

	manager.Status.Conditions = []v1alpha1.ManagerCondition{{
		Type:   v1alpha1.ManagerReady,
		Status: readyStatus,
	}}
}

func (r *ReconcileManager) getNodes(selector labels.Selector) ([]corev1.Node, error) {
	nodes := &corev1.NodeList{}
	listOpts := client.ListOptions{LabelSelector: selector}
	if err := r.client.List(context.Background(), nodes, &listOpts); err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

func (r *ReconcileManager) getReplicas(nodes []corev1.Node) int32 {
	nodesNumber := len(nodes)
	if nodesNumber%2 == 0 && nodesNumber != 0 {
		return int32(nodesNumber - 1)
	}
	return int32(nodesNumber)
}

func (r *ReconcileManager) getNodesHostAliases(nodes []corev1.Node) []corev1.HostAlias {
	hostAliases := []corev1.HostAlias{}
	for _, n := range nodes {
		ip := ""
		hostname := ""
		for _, a := range n.Status.Addresses {
			if a.Type == corev1.NodeInternalIP {
				ip = a.Address
			}
			if a.Type == corev1.NodeHostName {
				hostname = a.Address
			}
		}
		if ip == "" || hostname == "" {
			continue
		}
		hostAliases = append(hostAliases, corev1.HostAlias{
			IP:        ip,
			Hostnames: []string{hostname},
		})
	}

	return hostAliases
}

func (r *ReconcileManager) processZookeepers(manager *v1alpha1.Manager, replicas int32) error {
	for _, existingZookeeper := range manager.Status.Zookeepers {
		found := false
		for _, intendedZookeeper := range manager.Spec.Services.Zookeepers {
			if *existingZookeeper.Name == intendedZookeeper.Name {
				found = true
				break
			}
		}
		if !found {
			oldZookeeper := &v1alpha1.Zookeeper{}
			oldZookeeper.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingZookeeper.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldZookeeper)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var zookeeperServiceStatus []*v1alpha1.ServiceStatus
	for _, zookeeperService := range manager.Spec.Services.Zookeepers {
		zookeeper := &v1alpha1.Zookeeper{}
		zookeeper.ObjectMeta = zookeeperService.ObjectMeta
		zookeeper.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, zookeeper, func() error {
			zookeeper.Spec = zookeeperService.Spec
			zookeeper.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, zookeeper.Spec.CommonConfiguration)
			if zookeeper.Spec.CommonConfiguration.Replicas == nil {
				zookeeper.Spec.CommonConfiguration.Replicas = &replicas
			}
			return controllerutil.SetControllerReference(manager, zookeeper, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &zookeeper.Name
		status.Active = zookeeper.Status.Active
		zookeeperServiceStatus = append(zookeeperServiceStatus, status)
	}

	manager.Status.Zookeepers = zookeeperServiceStatus
	return nil
}

func (r *ReconcileManager) processCassandras(manager *v1alpha1.Manager, replicas int32, hostAliases []corev1.HostAlias) error {
	for _, existingCassandra := range manager.Status.Cassandras {
		found := false
		for _, intendedCassandra := range manager.Spec.Services.Cassandras {
			if *existingCassandra.Name == intendedCassandra.Name {
				found = true
				break
			}
		}
		if !found {
			oldCassandra := &v1alpha1.Cassandra{}
			oldCassandra.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingCassandra.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldCassandra)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var cassandraStatusList []*v1alpha1.ServiceStatus
	for _, cassandraService := range manager.Spec.Services.Cassandras {
		cassandra := &v1alpha1.Cassandra{}
		cassandra.ObjectMeta = cassandraService.ObjectMeta
		cassandra.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, cassandra, func() error {
			cassandra.Spec = cassandraService.Spec
			if cassandra.Spec.ServiceConfiguration.ClusterName == "" {
				cassandra.Spec.ServiceConfiguration.ClusterName = manager.GetName()
			}
			cassandra.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, cassandra.Spec.CommonConfiguration)
			if cassandra.Spec.CommonConfiguration.Replicas == nil {
				cassandra.Spec.CommonConfiguration.Replicas = &replicas
			}
			if len(cassandra.Spec.CommonConfiguration.HostAliases) == 0 {
				cassandra.Spec.CommonConfiguration.HostAliases = hostAliases
			}
			return controllerutil.SetControllerReference(manager, cassandra, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &cassandra.Name
		status.Active = cassandra.Status.Active
		cassandraStatusList = append(cassandraStatusList, status)
	}

	manager.Status.Cassandras = cassandraStatusList
	return nil
}

func (r *ReconcileManager) processWebui(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Webui == nil {
		if manager.Status.Webui != nil {
			oldWebUI := &v1alpha1.Webui{}
			oldWebUI.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Webui.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldWebUI)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Webui = nil
		}
		return nil
	}

	webui := &v1alpha1.Webui{}
	webui.ObjectMeta = manager.Spec.Services.Webui.ObjectMeta
	webui.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Webui.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, webui, func() error {
		webui.Spec = manager.Spec.Services.Webui.Spec
		webui.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, webui.Spec.CommonConfiguration)
		if webui.Spec.CommonConfiguration.Replicas == nil {
			webui.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, webui, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &webui.Name
	status.Active = &webui.Status.Active
	manager.Status.Webui = status
	return err
}

func (r *ReconcileManager) processProvisionManager(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.ProvisionManager == nil {
		if manager.Status.ProvisionManager != nil {
			oldPM := &v1alpha1.ProvisionManager{}
			oldPM.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.ProvisionManager.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldPM)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.ProvisionManager = nil
		}
		return nil
	}

	pm := &v1alpha1.ProvisionManager{}
	pm.ObjectMeta = manager.Spec.Services.ProvisionManager.ObjectMeta
	pm.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, pm, func() error {
		pm.Spec = manager.Spec.Services.ProvisionManager.Spec
		pm.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, pm.Spec.CommonConfiguration)
		if pm.Spec.CommonConfiguration.Replicas == nil {
			pm.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, pm, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &pm.Name
	status.Active = pm.Status.Active
	manager.Status.ProvisionManager = status
	return err
}

func (r *ReconcileManager) processConfig(manager *v1alpha1.Manager, replicas int32, hostAliases []corev1.HostAlias) error {
	if manager.Spec.Services.Config == nil {
		if manager.Status.Config != nil {
			oldConfig := &v1alpha1.Config{}
			oldConfig.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Config.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldConfig)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Config = nil
		}
		return nil
	}

	config := &v1alpha1.Config{}
	config.ObjectMeta = manager.Spec.Services.Config.ObjectMeta
	config.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Config.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, config, func() error {
		config.Spec = manager.Spec.Services.Config.Spec
		config.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, config.Spec.CommonConfiguration)
		if config.Spec.CommonConfiguration.Replicas == nil {
			config.Spec.CommonConfiguration.Replicas = &replicas
		}
		if len(config.Spec.CommonConfiguration.HostAliases) == 0 {
			config.Spec.CommonConfiguration.HostAliases = hostAliases
		}
		return controllerutil.SetControllerReference(manager, config, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &config.Name
	status.Active = config.Status.Active
	manager.Status.Config = status
	return err
}

func (r *ReconcileManager) processKubemanagers(manager *v1alpha1.Manager, replicas int32) error {
	for _, existingKubemanager := range manager.Status.Kubemanagers {
		found := false
		for _, intendedKubemanager := range manager.Spec.Services.Kubemanagers {
			if *existingKubemanager.Name == intendedKubemanager.Name {
				found = true
				break
			}
		}
		if !found {
			oldKubemanager := &v1alpha1.Kubemanager{}
			oldKubemanager.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingKubemanager.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldKubemanager)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var kubemanagerServiceStatus []*v1alpha1.ServiceStatus
	for _, kubemanagerService := range manager.Spec.Services.Kubemanagers {
		if !kubemanagerDependenciesReady(kubemanagerService.Spec.ServiceConfiguration.CassandraInstance, kubemanagerService.Spec.ServiceConfiguration.ZookeeperInstance, manager.ObjectMeta, r.client) {
			continue
		}
		kubemanager := &v1alpha1.Kubemanager{}
		kubemanager.ObjectMeta = kubemanagerService.ObjectMeta
		kubemanager.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, kubemanager, func() error {
			kubemanager.Spec.ServiceConfiguration.KubemanagerConfiguration = kubemanagerService.Spec.ServiceConfiguration.KubemanagerConfiguration
			if err := fillKubemanagerConfiguration(kubemanager, kubemanagerService.Spec.ServiceConfiguration.CassandraInstance, kubemanagerService.Spec.ServiceConfiguration.ZookeeperInstance, manager.ObjectMeta, r.client); err != nil {
				return err
			}
			kubemanager.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, kubemanagerService.Spec.CommonConfiguration)
			if kubemanager.Spec.CommonConfiguration.Replicas == nil {
				kubemanager.Spec.CommonConfiguration.Replicas = &replicas
			}
			return controllerutil.SetControllerReference(manager, kubemanager, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &kubemanager.Name
		status.Active = kubemanager.Status.Active
		kubemanagerServiceStatus = append(kubemanagerServiceStatus, status)
	}

	manager.Status.Kubemanagers = kubemanagerServiceStatus
	return nil
}

func (r *ReconcileManager) processControls(manager *v1alpha1.Manager, replicas int32) error {
	for _, existingControl := range manager.Status.Controls {
		found := false
		for _, intendedControl := range manager.Spec.Services.Controls {
			if *existingControl.Name == intendedControl.Name {
				found = true
				break
			}
		}
		if !found {
			oldControl := &v1alpha1.Control{}
			oldControl.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingControl.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldControl)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var controlServiceStatus []*v1alpha1.ServiceStatus
	for _, controlService := range manager.Spec.Services.Controls {
		control := &v1alpha1.Control{}
		control.ObjectMeta = controlService.ObjectMeta
		control.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, control, func() error {
			control.Spec = controlService.Spec
			control.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, control.Spec.CommonConfiguration)
			if control.Spec.CommonConfiguration.Replicas == nil {
				control.Spec.CommonConfiguration.Replicas = &replicas
			}
			return controllerutil.SetControllerReference(manager, control, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &control.Name
		status.Active = control.Status.Active
		controlServiceStatus = append(controlServiceStatus, status)
	}

	manager.Status.Controls = controlServiceStatus
	return nil
}

func (r *ReconcileManager) processRabbitMQ(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Rabbitmq == nil {
		if manager.Status.Rabbitmq != nil {
			oldRabbitMQ := &v1alpha1.Rabbitmq{}
			oldRabbitMQ.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Rabbitmq.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldRabbitMQ)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Rabbitmq = nil
		}
		return nil
	}
	rabbitMQ := &v1alpha1.Rabbitmq{}
	rabbitMQ.ObjectMeta = manager.Spec.Services.Rabbitmq.ObjectMeta
	rabbitMQ.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, rabbitMQ, func() error {
		rabbitMQ.Spec = manager.Spec.Services.Rabbitmq.Spec
		rabbitMQ.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, rabbitMQ.Spec.CommonConfiguration)
		if rabbitMQ.Spec.CommonConfiguration.Replicas == nil {
			rabbitMQ.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, rabbitMQ, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &rabbitMQ.Name
	status.Active = rabbitMQ.Status.Active
	manager.Status.Rabbitmq = status
	return err
}

func (r *ReconcileManager) processVRouters(manager *v1alpha1.Manager, replicas int32) error {
	for _, existingVRouter := range manager.Status.Vrouters {
		found := false
		for _, intendedVRouter := range manager.Spec.Services.Vrouters {
			if *existingVRouter.Name == intendedVRouter.Name {
				found = true
				break
			}
		}
		if !found {
			oldVRouter := &v1alpha1.Vrouter{}
			oldVRouter.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingVRouter.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldVRouter)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var vRouterServiceStatus []*v1alpha1.ServiceStatus
	for _, vRouterService := range manager.Spec.Services.Vrouters {
		if !vrouterDependenciesReady(vRouterService.Spec.ServiceConfiguration.ControlInstance, manager.ObjectMeta, r.client) {
			continue
		}
		vRouter := &v1alpha1.Vrouter{}
		vRouter.ObjectMeta = vRouterService.ObjectMeta
		vRouter.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, vRouter, func() error {
			vRouter.Spec.ServiceConfiguration.VrouterConfiguration = vRouterService.Spec.ServiceConfiguration.VrouterConfiguration
			vRouter.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, vRouterService.Spec.CommonConfiguration)
			if err := fillVrouterConfiguration(vRouter, vRouterService.Spec.ServiceConfiguration.ControlInstance, manager.ObjectMeta, r.client); err != nil {
				return err
			}
			if vRouter.Spec.CommonConfiguration.Replicas == nil {
				vRouter.Spec.CommonConfiguration.Replicas = &replicas
			}
			return controllerutil.SetControllerReference(manager, vRouter, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &vRouter.Name
		status.Active = vRouter.Status.Active
		vRouterServiceStatus = append(vRouterServiceStatus, status)
	}

	manager.Status.Vrouters = vRouterServiceStatus
	return nil
}

func (r *ReconcileManager) processContrailCNIs(manager *v1alpha1.Manager) error {
	for _, existingContrailCNI := range manager.Status.ContrailCNIs {
		found := false
		for _, intendedContrailCNI := range manager.Spec.Services.ContrailCNIs {
			if *existingContrailCNI.Name == intendedContrailCNI.Name {
				found = true
				break
			}
		}
		if !found {
			oldContrailCNI := &v1alpha1.ContrailCNI{}
			oldContrailCNI.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *existingContrailCNI.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldContrailCNI)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}

	var ContrailCNIServiceStatus []*v1alpha1.ServiceStatus
	for _, ContrailCNIService := range manager.Spec.Services.ContrailCNIs {
		ContrailCNI := &v1alpha1.ContrailCNI{}
		ContrailCNI.ObjectMeta = ContrailCNIService.ObjectMeta
		ContrailCNI.ObjectMeta.Namespace = manager.Namespace
		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, ContrailCNI, func() error {
			modifyContrailCNI(ContrailCNI, ContrailCNIService, manager)
			return controllerutil.SetControllerReference(manager, ContrailCNI, r.scheme)
		})
		if err != nil {
			return err
		}
		status := &v1alpha1.ServiceStatus{}
		status.Name = &ContrailCNI.Name
		status.Active = &ContrailCNI.Status.Active
		ContrailCNIServiceStatus = append(ContrailCNIServiceStatus, status)
	}

	manager.Status.ContrailCNIs = ContrailCNIServiceStatus
	return nil
}

func (r *ReconcileManager) processCommand(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Command == nil {
		if manager.Status.Command != nil {
			oldCommand := &v1alpha1.Command{}
			oldCommand.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Command.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldCommand)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Command = nil
		}
		return nil
	}

	command := &v1alpha1.Command{}
	command.ObjectMeta = manager.Spec.Services.Command.ObjectMeta
	command.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Command.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, command, func() error {
		command.Spec = manager.Spec.Services.Command.Spec
		command.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, command.Spec.CommonConfiguration)
		if command.Spec.ServiceConfiguration.ClusterName == "" {
			command.Spec.ServiceConfiguration.ClusterName = manager.GetName()
		}
		if command.Spec.CommonConfiguration.Replicas == nil {
			command.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, command, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &command.Name
	status.Active = &command.Status.Active
	manager.Status.Command = status
	return err
}

func (r *ReconcileManager) processKeystone(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Keystone == nil {
		if manager.Status.Keystone != nil {
			oldKeystone := &v1alpha1.Keystone{}
			oldKeystone.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Keystone.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldKeystone)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Keystone = nil
		}
		return nil
	}

	keystone := &v1alpha1.Keystone{}
	keystone.ObjectMeta = manager.Spec.Services.Keystone.ObjectMeta
	keystone.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Keystone.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, keystone, func() error {
		keystone.Spec = manager.Spec.Services.Keystone.Spec
		keystone.SetDefaultValues()
		keystone.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, keystone.Spec.CommonConfiguration)
		if keystone.Spec.CommonConfiguration.Replicas == nil {
			keystone.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, keystone, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &keystone.Name
	status.Active = &keystone.Status.Active
	manager.Status.Keystone = status
	return err
}

func (r *ReconcileManager) processPostgres(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Postgres == nil {
		if manager.Status.Postgres != nil {
			oldPostgres := &v1alpha1.Postgres{}
			oldPostgres.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Postgres.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldPostgres)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Postgres = nil
		}
		return nil
	}

	psql := &v1alpha1.Postgres{}
	psql.ObjectMeta = manager.Spec.Services.Postgres.ObjectMeta
	psql.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Postgres.Spec.ServiceConfiguration.RootPassSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, psql, func() error {
		psql.Spec = manager.Spec.Services.Postgres.Spec
		if psql.Spec.ServiceConfiguration.ListenPort == 0 {
			psql.Spec.ServiceConfiguration.ListenPort = 5432
		}
		psql.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, psql.Spec.CommonConfiguration)
		if psql.Spec.CommonConfiguration.Replicas == nil {
			psql.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, psql, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &psql.Name
	status.Active = &psql.Status.Active
	manager.Status.Postgres = status
	return err
}

func (r *ReconcileManager) processSwift(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Swift == nil {
		if manager.Status.Swift != nil {
			oldSwift := &v1alpha1.Swift{}
			oldSwift.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Swift.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldSwift)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Swift = nil
		}
		return nil
	}

	swift := &v1alpha1.Swift{}
	swift.ObjectMeta = manager.Spec.Services.Swift.ObjectMeta
	manager.Spec.Services.Swift.Spec.ServiceConfiguration.SwiftProxyConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	swift.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swift, func() error {
		swift.Spec = manager.Spec.Services.Swift.Spec
		swift.SetDefaultValues()
		swift.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, swift.Spec.CommonConfiguration)
		if swift.Spec.CommonConfiguration.Replicas == nil {
			swift.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, swift, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &swift.Name
	status.Active = &swift.Status.Active
	manager.Status.Swift = status
	return err
}

func (r *ReconcileManager) processMemcached(manager *v1alpha1.Manager, replicas int32) error {
	if manager.Spec.Services.Memcached == nil {
		if manager.Status.Memcached != nil {
			oldMemcached := &v1alpha1.Memcached{}
			oldMemcached.ObjectMeta = v1.ObjectMeta{
				Namespace: manager.Namespace,
				Name:      *manager.Status.Memcached.Name,
				Labels: map[string]string{
					"contrail_cluster": manager.Name,
				},
			}
			err := r.client.Delete(context.TODO(), oldMemcached)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			manager.Status.Memcached = nil
		}
		return nil
	}

	memcached := &v1alpha1.Memcached{}
	memcached.ObjectMeta = manager.Spec.Services.Memcached.ObjectMeta
	memcached.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, memcached, func() error {
		memcached.Spec = manager.Spec.Services.Memcached.Spec
		memcached.Spec.CommonConfiguration = utils.MergeCommonConfiguration(manager.Spec.CommonConfiguration, memcached.Spec.CommonConfiguration)
		if memcached.Spec.CommonConfiguration.Replicas == nil {
			memcached.Spec.CommonConfiguration.Replicas = &replicas
		}
		return controllerutil.SetControllerReference(manager, memcached, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Name = &memcached.Name
	status.Active = &memcached.Status.Active
	manager.Status.Memcached = status
	return err
}

func (r *ReconcileManager) processCSRSignerCaConfigMap(manager *v1alpha1.Manager) error {
	caCertificate := certificates.NewCACertificate(r.client, r.scheme, manager, "manager")
	if err := caCertificate.EnsureExists(); err != nil {
		return err
	}

	csrSignerCaConfigMap := &corev1.ConfigMap{}
	csrSignerCaConfigMap.ObjectMeta.Name = certificates.SignerCAConfigMapName
	csrSignerCaConfigMap.ObjectMeta.Namespace = manager.Namespace

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, csrSignerCaConfigMap, func() error {
		csrSignerCAValue, err := caCertificate.GetCaCert()
		if err != nil {
			return err
		}
		csrSignerCaConfigMap.Data = map[string]string{certificates.SignerCAFilename: string(csrSignerCAValue)}
		return controllerutil.SetControllerReference(manager, csrSignerCaConfigMap, r.scheme)
	})

	return err
}

func (r *ReconcileManager) processContrailmonitor(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Contrailmonitor == nil {
		return nil
	}
	cms := &v1alpha1.Contrailmonitor{}
	cms.ObjectMeta = manager.Spec.Services.Contrailmonitor.ObjectMeta
	cms.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, cms, func() error {
		cms.Spec = manager.Spec.Services.Contrailmonitor.Spec
		return controllerutil.SetControllerReference(manager, cms, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Active = &cms.Status.Active
	manager.Status.Contrailmonitor = status
	return err
}

func modifyContrailCNI(ContrailCNI, ContrailCNIService *v1alpha1.ContrailCNI, manager *v1alpha1.Manager) {
	ContrailCNI.Spec = ContrailCNIService.Spec
	if len(ContrailCNI.Spec.CommonConfiguration.NodeSelector) == 0 && len(manager.Spec.CommonConfiguration.NodeSelector) > 0 {
		(&ContrailCNI.Spec.CommonConfiguration).NodeSelector = manager.Spec.CommonConfiguration.NodeSelector
	}
	if ContrailCNI.Spec.CommonConfiguration.HostNetwork == nil && manager.Spec.CommonConfiguration.HostNetwork != nil {
		(&ContrailCNI.Spec.CommonConfiguration).HostNetwork = manager.Spec.CommonConfiguration.HostNetwork
	}
	if len(ContrailCNI.Spec.CommonConfiguration.ImagePullSecrets) == 0 && len(manager.Spec.CommonConfiguration.ImagePullSecrets) > 0 {
		(&ContrailCNI.Spec.CommonConfiguration).ImagePullSecrets = manager.Spec.CommonConfiguration.ImagePullSecrets
	}
	if len(ContrailCNI.Spec.CommonConfiguration.Tolerations) == 0 && len(manager.Spec.CommonConfiguration.Tolerations) > 0 {
		(&ContrailCNI.Spec.CommonConfiguration).Tolerations = manager.Spec.CommonConfiguration.Tolerations
	}
}

func kubemanagerDependenciesReady(cassandraName, zookeeperName string, managerMeta v1.ObjectMeta, client client.Client) bool {
	cassandraInstance := v1alpha1.Cassandra{}
	zookeeperInstance := v1alpha1.Zookeeper{}
	rabbitmqInstance := v1alpha1.Rabbitmq{}
	configInstance := v1alpha1.Config{}

	cassandraActive := cassandraInstance.IsActive(cassandraName, managerMeta.Namespace, client)
	zookeeperActive := zookeeperInstance.IsActive(zookeeperName, managerMeta.Namespace, client)
	rabbitmqActive := rabbitmqInstance.IsActive(managerMeta.Name,
		managerMeta.Namespace, client)
	configActive := configInstance.IsActive(managerMeta.Name,
		managerMeta.Namespace, client)

	return cassandraActive && zookeeperActive && rabbitmqActive && configActive
}

func fillKubemanagerConfiguration(kubemanager *v1alpha1.Kubemanager, cassandraName, zookeeperName string, managerMeta v1.ObjectMeta, client client.Client) error {
	cassandraConfig, err := v1alpha1.NewCassandraClusterConfiguration(cassandraName, managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&kubemanager.Spec.ServiceConfiguration).CassandraNodesConfiguration = &cassandraConfig
	zookeeperConfig, err := v1alpha1.NewZookeeperClusterConfiguration(zookeeperName, managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&kubemanager.Spec.ServiceConfiguration).ZookeeperNodesConfiguration = &zookeeperConfig
	rabbitmqConfig, err := v1alpha1.NewRabbitmqClusterConfiguration(managerMeta.Name, managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&kubemanager.Spec.ServiceConfiguration).RabbbitmqNodesConfiguration = &rabbitmqConfig
	configConfig, err := v1alpha1.NewConfigClusterConfiguration(managerMeta.Name, managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&kubemanager.Spec.ServiceConfiguration).ConfigNodesConfiguration = &configConfig
	return nil
}

func vrouterDependenciesReady(controlName string, managerMeta v1.ObjectMeta, client client.Client) bool {
	controlInstance := v1alpha1.Control{}
	configInstance := v1alpha1.Config{}
	configInstanceActive := configInstance.IsActive(managerMeta.Name, managerMeta.Namespace, client)
	controlInstanceActive := controlInstance.IsActive(controlName, managerMeta.Namespace, client)

	return configInstanceActive && controlInstanceActive
}

func fillVrouterConfiguration(vrouter *v1alpha1.Vrouter, controlName string, managerMeta v1.ObjectMeta, client client.Client) error {
	controlConfig, err := v1alpha1.NewControlClusterConfiguration(controlName, "", managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&vrouter.Spec.ServiceConfiguration).ControlNodesConfiguration = &controlConfig
	configConfig, err := v1alpha1.NewConfigClusterConfiguration(managerMeta.Name, managerMeta.Namespace, client)
	if err != nil {
		return err
	}
	(&vrouter.Spec.ServiceConfiguration).ConfigNodesConfiguration = &configConfig
	return nil
}
