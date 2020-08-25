package swiftproxy

import (
	"context"
	"fmt"
	"strings"
	"time"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
)

func (r *ReconcileSwiftProxy) ensureSwiftRegistered(sp *contrail.SwiftProxy, adminSecret, swiftSecret *core.Secret, k *contrail.Keystone) (reconcile.Result, error) {

	jobNamespacedName := types.NamespacedName{Name: sp.Name + "-swift-register-job", Namespace: sp.Namespace}
	job := &batch.Job{}
	err := r.client.Get(context.TODO(), jobNamespacedName, job)
	if err == nil {
		if job.Status.CompletionTime.IsZero() {
			log.Info(fmt.Sprintf("job %v in progress", jobNamespacedName))
			return reconcile.Result{}, nil
		}
		if err := r.client.Delete(context.TODO(), job, client.PropagationPolicy(meta.DeletePropagationForeground)); err != nil {
			return reconcile.Result{}, err
		}
		// We have to wait for some time until job gets deleted because r.client.Delete does not delete job synchronously.
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 5,
		}, nil
	}
	if !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	if err := r.ensureRegisterJobConfig(jobConfigMapName(sp), sp, adminSecret, swiftSecret, k); err != nil {
		return reconcile.Result{}, err
	}

	job = newBootstrapJob(
		jobNamespacedName, jobConfigMapName(sp),
		sp.Spec.ServiceConfiguration.Containers, sp.Spec.CommonConfiguration.Tolerations,
	)
	if err = controllerutil.SetControllerReference(sp, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, r.client.Create(context.TODO(), job)
}

func jobConfigMapName(sp *contrail.SwiftProxy) string {
	return sp.Name + "-swiftproxy-register-job"
}

func (r *ReconcileSwiftProxy) ensureRegisterJobConfig(
	jobConfigName string,
	sp *contrail.SwiftProxy,
	adminSecret,
	swiftSecret *core.Secret,
	k *contrail.Keystone,
) error {
	keystoneData := &keystoneEndpoint{
		address:         k.Status.Endpoint,
		port:            k.Spec.ServiceConfiguration.ListenPort,
		region:          k.Spec.ServiceConfiguration.Region,
		authProtocol:    k.Spec.ServiceConfiguration.AuthProtocol,
		userDomainID:    k.Spec.ServiceConfiguration.UserDomainID,
		projectDomainID: k.Spec.ServiceConfiguration.ProjectDomainID,
	}

	cm := r.configMap(jobConfigName, sp, keystoneData, adminSecret, swiftSecret)
	publicIP := "0.0.0.0"
	if sp.Status.LoadBalancerIP != "" {
		publicIP = sp.Status.LoadBalancerIP
	}
	clusterIP := "0.0.0.0"
	if sp.Status.ClusterIP != "" {
		clusterIP = sp.Status.ClusterIP
	}
	if err := cm.ensureServiceExists(clusterIP, publicIP); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileSwiftProxy) isSwiftRegistered(sp *contrail.SwiftProxy, k *contrail.Keystone, swiftSecret *core.Secret) (bool, error) {
	keystoneClient, err := keystone.NewClient(r.client, r.scheme, r.mgrConfig, k)
	if err != nil {
		return false, err
	}
	token, err := keystoneClient.PostAuthTokens(string(swiftSecret.Data["user"]), string(swiftSecret.Data["password"]), "service")
	if keystone.IsUnauthorized(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to get keystone token: %v", err)
	}

	url := token.EndpointURL("swift", "internal")
	if !strings.Contains(url, sp.Status.ClusterIP) {
		return false, nil
	}
	url = token.EndpointURL("swift", "public")
	if sp.Status.LoadBalancerIP != "" && !strings.Contains(url, sp.Status.LoadBalancerIP) {
		return false, nil
	}

	return true, nil
}

func newBootstrapJob(
	nameSpacedName types.NamespacedName,
	jobConfigName string,
	containers []*contrail.Container,
	tolerations []core.Toleration,
) *batch.Job {
	return &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      nameSpacedName.Name,
			Namespace: nameSpacedName.Namespace,
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					Volumes: []core.Volume{
						{
							Name: "register",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: jobConfigName,
									},
								},
							},
						},
						{
							Name: "csr-signer-ca",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: certificates.SignerCAConfigMapName,
									},
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:            "register",
							Image:           getImage(containers, "init"),
							ImagePullPolicy: core.PullAlways,
							Command:         getCommand(containers, "init"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "register", MountPath: "/var/lib/ansible/register", ReadOnly: true},
								core.VolumeMount{Name: "csr-signer-ca", MountPath: certificates.SignerCAMountPath, ReadOnly: true},
							},
							Args: []string{"/var/lib/ansible/register/register.yaml", "-e", "@/var/lib/ansible/register/config.yaml"},
						},
					},
					Tolerations: tolerations,
				},
			},
		},
	}
}
