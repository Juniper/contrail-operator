package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	cr "github.com/Juniper/contrail-operator/pkg/controller/manager/crs"
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
	return addManagerWatch(c)
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

func addManagerWatch(c controller.Controller) error {
	err := c.Watch(&source.Kind{Type: &v1alpha1.Manager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	for _, resource := range resourcesList {
		if err = addResourcesToWatch(c, resource); err != nil {
			return err
		}
	}
	return nil
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

	// Create CRDs
	/*
		cassandraCrdActive := false
		for _, crdStatus := range instance.Status.CrdStatus {
			if crdStatus.Name == "Cassandra" {
				if crdStatus.Active != nil {
					if *crdStatus.Active {
						cassandraCrdActive = true
					}
				}
			}
		}
		if !cassandraCrdActive && len(instance.Spec.Services.Cassandras) > 0 {
			cassandraResource := v1alpha1.Cassandra{}
			cassandraCrd := cassandraResource.GetCrd()
			err = r.createCrd(instance, cassandraCrd)
			if err != nil {
				return reconcile.Result{}, err
			}
			controllerRunning := false
			c := r.manager.GetCache()
			sharedIndexInformer, err := c.GetInformer(&v1alpha1.Cassandra{})
			if err == nil {
				store := sharedIndexInformer.GetStore()
				if store != nil {
					fmt.Println("STORE NOT NIL")
				} else {
					fmt.Println("STORE NIL")
				}
				if sharedIndexInformer.HasSynced() {
					fmt.Println("has synced")
				} else {
					fmt.Println("has not synced")
				}

				controller := sharedIndexInformer.GetController()
				if controller != nil {
					controllerRunning = true
				}
			}
			if !controllerRunning {
				err = cassandra.Add(r.manager)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			active := true
			crdStatus := &v1alpha1.CrdStatus{
				Name:   "Cassandra",
				Active: &active,
			}
			var crdStatusList []v1alpha1.CrdStatus
			if len(instance.Status.CrdStatus) > 0 {
				crdStatusList = instance.Status.CrdStatus
			}
			crdStatusList = append(crdStatusList, *crdStatus)
			instance.Status.CrdStatus = crdStatusList
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}

		}
	*/
	// Create CRs
	for _, cassandraService := range instance.Spec.Services.Cassandras {
		create := *cassandraService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetCassandraCr()
		cr.ObjectMeta = cassandraService.ObjectMeta
		cr.Labels = cassandraService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = cassandraService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Cassandra"
		for _, cassandraStatus := range instance.Status.Cassandras {
			if cassandraService.Name == *cassandraStatus.Name {
				if *cassandraService.Spec.CommonConfiguration.Create && *cassandraStatus.Created {
					err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
					if err == nil {
						create = false
						delete = false
						update = true
					}
				}
				if !*cassandraService.Spec.CommonConfiguration.Create && *cassandraStatus.Created {
					create = false
					delete = true
					update = false
				}
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					if err = controllerutil.SetControllerReference(instance, cr, r.scheme); err != nil {
						return reconcile.Result{}, err
					}
					if err = r.client.Create(context.TODO(), cr); err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			cassandraStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Cassandras != nil {
				for _, cassandraStatus := range instance.Status.Cassandras {
					if cassandraService.Name == *cassandraStatus.Name {
						status = cassandraStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				cassandraStatusList = append(cassandraStatusList, status)
				instance.Status.Cassandras = cassandraStatusList
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if cassandraService.Spec.CommonConfiguration.Replicas != nil {
				replicas = cassandraService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range cassandraService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			cassandraStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Cassandras != nil {
				for _, cassandraStatus := range instance.Status.Cassandras {
					if cassandraService.Name == *cassandraStatus.Name {
						status = cassandraStatus
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				cassandraStatusList = append(cassandraStatusList, status)
				instance.Status.Cassandras = cassandraStatusList
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			cassandraStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Cassandras != nil {
				for _, cassandraStatus := range instance.Status.Cassandras {
					if cassandraService.Name == *cassandraStatus.Name {
						status = cassandraStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				cassandraStatusList = append(cassandraStatusList, status)
				instance.Status.Cassandras = cassandraStatusList
			}
		}
	}

	for _, zookeeperService := range instance.Spec.Services.Zookeepers {
		create := *zookeeperService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetZookeeperCr()
		cr.ObjectMeta = zookeeperService.ObjectMeta
		cr.Labels = zookeeperService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = zookeeperService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Zookeeper"
		for _, zookeeperStatus := range instance.Status.Zookeepers {
			if zookeeperService.Name == *zookeeperStatus.Name {
				if *zookeeperService.Spec.CommonConfiguration.Create && *zookeeperStatus.Created {
					err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
					if err == nil {
						create = false
						delete = false
						update = true
					}
				}
				if !*zookeeperService.Spec.CommonConfiguration.Create && *zookeeperStatus.Created {
					create = false
					delete = true
					update = false
				}
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			zookeeperStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Zookeepers != nil {
				for _, zookeeperStatus := range instance.Status.Zookeepers {
					if zookeeperService.Name == *zookeeperStatus.Name {
						status = zookeeperStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				zookeeperStatusList = append(zookeeperStatusList, status)
				instance.Status.Zookeepers = zookeeperStatusList
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if zookeeperService.Spec.CommonConfiguration.Replicas != nil {
				replicas = zookeeperService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range zookeeperService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			zookeeperStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Zookeepers != nil {
				for _, zookeeperStatus := range instance.Status.Zookeepers {
					if zookeeperService.Name == *zookeeperStatus.Name {
						status = zookeeperStatus
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				zookeeperStatusList = append(zookeeperStatusList, status)
				instance.Status.Zookeepers = zookeeperStatusList
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			zookeeperStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Zookeepers != nil {
				for _, zookeeperStatus := range instance.Status.Zookeepers {
					if zookeeperService.Name == *zookeeperStatus.Name {
						status = zookeeperStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				zookeeperStatusList = append(zookeeperStatusList, status)
				instance.Status.Zookeepers = zookeeperStatusList
			}
		}
	}

	if instance.Spec.Services.Webui != nil {
		webuiService := instance.Spec.Services.Webui
		create := *webuiService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		webuiService.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName

		cr := cr.GetWebuiCr()
		cr.ObjectMeta = webuiService.ObjectMeta
		cr.Labels = webuiService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = webuiService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Webui"
		if instance.Status.Webui != nil {
			if *webuiService.Spec.CommonConfiguration.Create && *instance.Status.Webui.Created {
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
				if err == nil {
					create = false
					delete = false
					update = true
				}
			}
			if !*webuiService.Spec.CommonConfiguration.Create && *instance.Status.Webui.Created {
				create = false
				delete = true
				update = false
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Webui != nil {
				status = instance.Status.Webui
				status.Created = &create
				status.Active = &cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = &cr.Status.Active
				instance.Status.Webui = status
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if webuiService.Spec.CommonConfiguration.Replicas != nil {
				replicas = webuiService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range webuiService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			secretParamChanged := false
			if cr.Spec.ServiceConfiguration.KeystoneSecretName == "" {
				cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
				secretParamChanged = true
			}

			if imageChanged || replicasChanged || secretParamChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Webui != nil {
				status = instance.Status.Webui
				status.Active = &cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = &cr.Status.Active
				instance.Status.Webui = status
			}

		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Webui != nil {
				status = instance.Status.Webui
				status.Created = &create
				status.Active = &cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = &cr.Status.Active
				instance.Status.Webui = status
			}
		}
	}

	if instance.Spec.Services.ProvisionManager != nil {
		provisionManagerService := instance.Spec.Services.ProvisionManager
		create := *provisionManagerService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		provisionManagerService.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName

		cr := cr.GetProvisionManagerCr()
		cr.ObjectMeta = provisionManagerService.ObjectMeta
		cr.Labels = provisionManagerService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = provisionManagerService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "ProvisionManager"
		if instance.Status.ProvisionManager != nil {
			if *provisionManagerService.Spec.CommonConfiguration.Create && *instance.Status.ProvisionManager.Created {
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
				if err == nil {
					create = false
					delete = false
					update = true
				}
			}
			if !*provisionManagerService.Spec.CommonConfiguration.Create && *instance.Status.ProvisionManager.Created {
				create = false
				delete = true
				update = false
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			if instance.Status.ProvisionManager != nil {
				status = instance.Status.ProvisionManager
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.ProvisionManager = status
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if provisionManagerService.Spec.CommonConfiguration.Replicas != nil {
				replicas = provisionManagerService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}

			imageChanged := false
			for _, container := range provisionManagerService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			secretParamChanged := false
			if cr.Spec.ServiceConfiguration.KeystoneSecretName == "" {
				cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
				secretParamChanged = true

			}
			if imageChanged || replicasChanged || secretParamChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.ProvisionManager != nil {
				status = instance.Status.ProvisionManager
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.ProvisionManager = status
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.ProvisionManager != nil {
				status = instance.Status.ProvisionManager
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.ProvisionManager = status
			}
		}
	}

	if instance.Spec.Services.Config != nil {
		configService := instance.Spec.Services.Config
		create := *configService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		configService.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName

		cr := cr.GetConfigCr()
		cr.ObjectMeta = configService.ObjectMeta
		cr.Labels = configService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = configService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Config"

		if instance.Status.Config != nil {
			if *configService.Spec.CommonConfiguration.Create && *instance.Status.Config.Created {
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
				if err == nil {
					create = false
					delete = false
					update = true
				}
			}
			if !*configService.Spec.CommonConfiguration.Create && *instance.Status.Config.Created {
				create = false
				delete = true
				update = false
			}
		}

		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Config != nil {
				status = instance.Status.Config
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Config = status
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if configService.Spec.CommonConfiguration.Replicas != nil {
				replicas = configService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range configService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			secretParamChanged := false
			if cr.Spec.ServiceConfiguration.KeystoneSecretName == "" {
				cr.Spec.ServiceConfiguration.KeystoneSecretName = instance.Spec.KeystoneSecretName
				secretParamChanged = true
			}

			if imageChanged || replicasChanged || secretParamChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Config != nil {
				status = instance.Status.Config
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Config = status
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Config != nil {
				status = instance.Status.Config
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Config = status
			}
		}
	}

	for _, kubemanagerService := range instance.Spec.Services.Kubemanagers {
		create := *kubemanagerService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetKubemanagerCr()
		cr.ObjectMeta = kubemanagerService.ObjectMeta
		cr.Labels = kubemanagerService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = kubemanagerService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Kubemanager"
		for _, kubemanagerStatus := range instance.Status.Kubemanagers {
			if kubemanagerService.Name == *kubemanagerStatus.Name {
				if *kubemanagerService.Spec.CommonConfiguration.Create && *kubemanagerStatus.Created {
					err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
					if err == nil {
						create = false
						delete = false
						update = true
					}
				}
				if !*kubemanagerService.Spec.CommonConfiguration.Create && *kubemanagerStatus.Created {
					create = false
					delete = true
					update = false
				}
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			kubemanagerStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Kubemanagers != nil {
				for _, kubemanagerStatus := range instance.Status.Kubemanagers {
					if kubemanagerService.Name == *kubemanagerStatus.Name {
						status = kubemanagerStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				kubemanagerStatusList = append(kubemanagerStatusList, status)
				instance.Status.Kubemanagers = kubemanagerStatusList
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if kubemanagerService.Spec.CommonConfiguration.Replicas != nil {
				replicas = kubemanagerService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range kubemanagerService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			kubemanagerStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Kubemanagers != nil {
				for _, kubemanagerStatus := range instance.Status.Kubemanagers {
					if kubemanagerService.Name == *kubemanagerStatus.Name {
						status = kubemanagerStatus
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				kubemanagerStatusList = append(kubemanagerStatusList, status)
				instance.Status.Kubemanagers = kubemanagerStatusList
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			kubemanagerStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Kubemanagers != nil {
				for _, kubemanagerStatus := range instance.Status.Kubemanagers {
					if kubemanagerService.Name == *kubemanagerStatus.Name {
						status = kubemanagerStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				kubemanagerStatusList = append(kubemanagerStatusList, status)
				instance.Status.Kubemanagers = kubemanagerStatusList
			}
		}
	}

	for _, controlService := range instance.Spec.Services.Controls {
		create := *controlService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetControlCr()
		cr.ObjectMeta = controlService.ObjectMeta
		cr.Labels = controlService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = controlService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Control"
		for _, controlStatus := range instance.Status.Controls {
			if controlService.Name == *controlStatus.Name {
				if *controlService.Spec.CommonConfiguration.Create && *controlStatus.Created {
					err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
					if err == nil {
						create = false
						delete = false
						update = true
					}
				}
				if !*controlService.Spec.CommonConfiguration.Create && *controlStatus.Created {
					create = false
					delete = true
					update = false
				}
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			controlStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Controls != nil {
				for _, controlStatus := range instance.Status.Controls {
					if controlService.Name == *controlStatus.Name {
						status = controlStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				controlStatusList = append(controlStatusList, status)
				instance.Status.Controls = controlStatusList
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if controlService.Spec.CommonConfiguration.Replicas != nil {
				replicas = controlService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range controlService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			controlStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Controls != nil {
				for _, controlStatus := range instance.Status.Controls {
					if controlService.Name == *controlStatus.Name {
						status = controlStatus
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				controlStatusList = append(controlStatusList, status)
				instance.Status.Controls = controlStatusList
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			controlStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Controls != nil {
				for _, controlStatus := range instance.Status.Controls {
					if controlService.Name == *controlStatus.Name {
						status = controlStatus
						status.Created = &create
						status.Active = cr.Status.Active
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				controlStatusList = append(controlStatusList, status)
				instance.Status.Controls = controlStatusList
			}

		}
	}

	if instance.Spec.Services.Rabbitmq != nil {
		rabbitmqService := instance.Spec.Services.Rabbitmq
		create := *rabbitmqService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetRabbitmqCr()
		cr.ObjectMeta = rabbitmqService.ObjectMeta
		cr.Labels = rabbitmqService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = rabbitmqService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Rabbitmq"
		if instance.Status.Rabbitmq != nil {
			if *rabbitmqService.Spec.CommonConfiguration.Create && *instance.Status.Rabbitmq.Created {
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
				if err == nil {
					create = false
					delete = false
					update = true
				}
			}
			if !*rabbitmqService.Spec.CommonConfiguration.Create && *instance.Status.Rabbitmq.Created {
				create = false
				delete = true
				update = false
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Rabbitmq != nil {
				status = instance.Status.Rabbitmq
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Rabbitmq = status
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if rabbitmqService.Spec.CommonConfiguration.Replicas != nil {
				replicas = rabbitmqService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range rabbitmqService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Rabbitmq != nil {
				status = instance.Status.Rabbitmq
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Rabbitmq = status
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Delete(context.TODO(), cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			status := &v1alpha1.ServiceStatus{}
			if instance.Status.Rabbitmq != nil {
				status = instance.Status.Rabbitmq
				status.Created = &create
				status.Active = cr.Status.Active
			} else {
				status.Name = &cr.Name
				status.Created = &create
				status.Active = cr.Status.Active
				instance.Status.Rabbitmq = status
			}
		}
	}

	for _, vrouterService := range instance.Spec.Services.Vrouters {
		create := *vrouterService.Spec.CommonConfiguration.Create
		delete := false
		update := false

		cr := cr.GetVrouterCr()
		cr.ObjectMeta = vrouterService.ObjectMeta
		cr.Labels = vrouterService.ObjectMeta.Labels
		cr.Namespace = instance.Namespace
		cr.Spec.ServiceConfiguration = vrouterService.Spec.ServiceConfiguration
		cr.TypeMeta.APIVersion = "contrail.juniper.net/v1alpha1"
		cr.TypeMeta.Kind = "Vrouter"
		for _, vrouterStatus := range instance.Status.Vrouters {
			if vrouterService.Name == *vrouterStatus.Name {
				if *vrouterService.Spec.CommonConfiguration.Create && *vrouterStatus.Created {
					err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
					if err == nil {
						create = false
						delete = false
						update = true
					}
				}
				if !*vrouterService.Spec.CommonConfiguration.Create && *vrouterStatus.Created {
					create = false
					delete = true
					update = false
				}
			}
		}
		if create {
			err = r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				if errors.IsNotFound(err) {
					err = controllerutil.SetControllerReference(instance, cr, r.scheme)
					if err != nil {
						return reconcile.Result{}, err
					}
					err = r.client.Create(context.TODO(), cr)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}

			status := &v1alpha1.ServiceStatus{}

			vrouterStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Vrouters != nil {
				vrouterStatusList = instance.Status.Vrouters
			}
			if instance.Status.Vrouters != nil {
				vrouterFound := false
				for _, vrouterStatus := range instance.Status.Vrouters {
					if vrouterService.Name == *vrouterStatus.Name {
						vrouterFound = true
						status = vrouterStatus
						status.Created = &create
					}
				}
				if !vrouterFound {
					status.Name = &cr.Name
					status.Created = &create
					vrouterStatusList = append(vrouterStatusList, status)
					instance.Status.Vrouters = vrouterStatusList
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				vrouterStatusList = append(vrouterStatusList, status)
				instance.Status.Vrouters = vrouterStatusList
			}
		}
		if update {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicasChanged := false
			replicas := instance.Spec.CommonConfiguration.Replicas
			if vrouterService.Spec.CommonConfiguration.Replicas != nil {
				replicas = vrouterService.Spec.CommonConfiguration.Replicas
			}
			if cr.Spec.CommonConfiguration.Replicas != nil {
				if *replicas != *cr.Spec.CommonConfiguration.Replicas {
					cr.Spec.CommonConfiguration.Replicas = replicas
					replicasChanged = true
				}
			}
			imageChanged := false
			for _, container := range vrouterService.Spec.ServiceConfiguration.Containers {
				for idx, crContainer := range cr.Spec.ServiceConfiguration.Containers {
					if crContainer.Name == container.Name {
						if crContainer.Image != container.Image {
							cr.Spec.ServiceConfiguration.Containers[idx].Image = container.Image
							imageChanged = true
							break
						}
					}
				}
			}
			if imageChanged || replicasChanged {
				err = r.client.Update(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
		if delete {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
			if err == nil {
				err = r.client.Delete(context.TODO(), cr)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			status := &v1alpha1.ServiceStatus{}
			vrouterStatusList := []*v1alpha1.ServiceStatus{}
			if instance.Status.Vrouters != nil {
				vrouterStatusList = instance.Status.Vrouters
			}
			if instance.Status.Vrouters != nil {
				for _, vrouterStatus := range instance.Status.Vrouters {
					if vrouterService.Name == *vrouterStatus.Name {
						status = vrouterStatus
						status.Created = &create
					}
				}
			} else {
				status.Name = &cr.Name
				status.Created = &create
				vrouterStatusList = append(vrouterStatusList, status)
				instance.Status.Vrouters = vrouterStatusList
			}
		}
	}

	if err = r.processPostgres(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processCommand(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processKeystone(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processSwift(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.processMemcached(instance); err != nil {
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

func (r *ReconcileManager) processCommand(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Command == nil {
		return nil
	}
	command := &v1alpha1.Command{}
	command.ObjectMeta = manager.Spec.Services.Command.ObjectMeta
	command.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Command.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, command, func() error {
		command.Spec = manager.Spec.Services.Command.Spec
		if command.Spec.ServiceConfiguration.ClusterName == "" {
			command.Spec.ServiceConfiguration.ClusterName = manager.GetName()
		}
		return controllerutil.SetControllerReference(manager, command, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Active = &command.Status.Active
	manager.Status.Command = status
	return err
}

func (r *ReconcileManager) processKeystone(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Keystone == nil {
		return nil
	}

	keystone := &v1alpha1.Keystone{}
	keystone.ObjectMeta = manager.Spec.Services.Keystone.ObjectMeta
	keystone.ObjectMeta.Namespace = manager.Namespace
	manager.Spec.Services.Keystone.Spec.ServiceConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, keystone, func() error {
		keystone.Spec = manager.Spec.Services.Keystone.Spec
		return controllerutil.SetControllerReference(manager, keystone, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Active = &keystone.Status.Active
	manager.Status.Keystone = status
	return err
}

func (r *ReconcileManager) processPostgres(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Postgres == nil {
		return nil
	}
	psql := &v1alpha1.Postgres{}
	psql.ObjectMeta = manager.Spec.Services.Postgres.ObjectMeta
	psql.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, psql, func() error {
		psql.Spec = manager.Spec.Services.Postgres.Spec
		return controllerutil.SetControllerReference(manager, psql, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Active = &psql.Status.Active
	manager.Status.Postgres = status
	return err
}

func (r *ReconcileManager) processSwift(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Swift == nil {
		return nil
	}
	swift := &v1alpha1.Swift{}
	swift.ObjectMeta = manager.Spec.Services.Swift.ObjectMeta
	manager.Spec.Services.Swift.Spec.ServiceConfiguration.SwiftProxyConfiguration.KeystoneSecretName = manager.Spec.KeystoneSecretName
	swift.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swift, func() error {
		swift.Spec = manager.Spec.Services.Swift.Spec
		return controllerutil.SetControllerReference(manager, swift, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
	status.Active = &swift.Status.Active
	manager.Status.Swift = status
	return err
}

func (r *ReconcileManager) processMemcached(manager *v1alpha1.Manager) error {
	if manager.Spec.Services.Memcached == nil {
		return nil
	}
	memcached := &v1alpha1.Memcached{}
	memcached.ObjectMeta = manager.Spec.Services.Memcached.ObjectMeta
	memcached.ObjectMeta.Namespace = manager.Namespace
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, memcached, func() error {
		memcached.Spec = manager.Spec.Services.Memcached.Spec
		return controllerutil.SetControllerReference(manager, memcached, r.scheme)
	})
	status := &v1alpha1.ServiceStatus{}
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
