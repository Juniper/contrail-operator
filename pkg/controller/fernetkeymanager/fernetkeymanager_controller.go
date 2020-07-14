package fernetkeymanager

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_fernetkeymanager")

// Add creates a new FernetKeyManager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()))
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fernetkeymanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FernetKeyManager
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.FernetKeyManager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFernetKeyManager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFernetKeyManager{}

// ReconcileFernetKeyManager reconciles a FernetKeyManager object
type ReconcileFernetKeyManager struct {
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes) *ReconcileFernetKeyManager {
	return &ReconcileFernetKeyManager{client: client, scheme: scheme, kubernetes: kubernetes}
}

// Reconcile reads that state of the cluster for a FernetKeyManager object and makes changes based on the state read
// and what is in the FernetKeyManager.Spec
func (r *ReconcileFernetKeyManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling FernetKeyManager")

	fernetKeyManager := &contrailv1alpha1.FernetKeyManager{}
	err := r.client.Get(context.TODO(), request.NamespacedName, fernetKeyManager)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !fernetKeyManager.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	keySecretName := "fernet-keys-repository"
	if err = r.secret(keySecretName, "fernetKeyManager", fernetKeyManager).ensureKeysRepoSecretExists(); err != nil {
		return reconcile.Result{}, err
	}

	tokenExpiration := fernetKeyManager.Spec.TokenExpiration
	if tokenExpiration == 0 {
		tokenExpiration = 86400
	}
	allowExpiredWindow := fernetKeyManager.Spec.TokenAllowExpiredWindow
	if allowExpiredWindow == 0 {
		allowExpiredWindow = 172800
	}
	rotationInterval := fernetKeyManager.Spec.RotationInterval
	if rotationInterval == 0 {
		rotationInterval = tokenExpiration + allowExpiredWindow
	}

	maxActiveKeys := ((tokenExpiration + allowExpiredWindow + rotationInterval - 1) / rotationInterval) + 2
	keySecret, err := r.getSecret(keySecretName, fernetKeyManager.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := r.rotateKeys(keySecret, maxActiveKeys); err != nil {
		return reconcile.Result{}, err
	}

	fernetKeyManager.Status.SecretName = keySecretName
	if err := r.client.Status().Update(context.TODO(), fernetKeyManager); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: time.Second * time.Duration(rotationInterval),
	}, nil
}

func (r *ReconcileFernetKeyManager) getSecret(secretName, secretNamespace string) (*core.Secret, error) {
	secret := &core.Secret{}
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNamespace}
	if err := r.client.Get(context.Background(), namespacedName, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (r *ReconcileFernetKeyManager) rotateKeys(sc *core.Secret, maxActiveKeys int) error {
	keys := sc.Data
	existingKeysIndices := make([]int, 0, len(keys))
	for k := range keys {
		key, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		existingKeysIndices = append(existingKeysIndices, key)
	}

	activeKeysNumber := len(existingKeysIndices)
	if activeKeysNumber == 0 {
		return fmt.Errorf("key repository not initialized, secret is empty")
	}
	log.Info(fmt.Sprintf("Starting rotation with %d keys", activeKeysNumber))

	sort.Ints(existingKeysIndices)
	maxKeyIndex := existingKeysIndices[activeKeysNumber-1]
	log.Info(fmt.Sprintf("Current primary is: %d", maxKeyIndex))
	log.Info(fmt.Sprintf("Next primary key will be: %d", maxKeyIndex+1))

	stagedKeyIndex := strconv.Itoa(0)
	newPrimary := keys[stagedKeyIndex]
	keys[strconv.Itoa(maxKeyIndex+1)] = newPrimary
	delete(keys, stagedKeyIndex)

	newKey, err := generateKey()
	if err != nil {
		return err
	}
	keys[stagedKeyIndex] = newKey
	log.Info("Promoted key 0 to be primary key")

	if len(keys) > maxActiveKeys {
		minKeyIndex := existingKeysIndices[1]
		log.Info(fmt.Sprintf("Excess key to purge: %d", minKeyIndex))
		delete(keys, strconv.Itoa(minKeyIndex))
	}

	sc.Data = keys
	return r.client.Update(context.TODO(), sc)
}
