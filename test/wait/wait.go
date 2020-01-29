package wait

import (
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// Wait is used to wait until certain resource reach some condition
type Wait struct {
	Namespace     string
	RetryInterval time.Duration
	Timeout       time.Duration
	KubeClient    kubernetes.Interface
}

// ForReadyStatefulSet is used to wait until StatefulSet is ready
func (w Wait) ForReadyStatefulSet(name string) error {
	return wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		statefulSet, err := w.KubeClient.AppsV1().StatefulSets(w.Namespace).Get(name, meta.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		replicas := int32(1)
		if statefulSet.Spec.Replicas != nil {
			replicas = *statefulSet.Spec.Replicas
		}

		if statefulSet.Status.ReadyReplicas == replicas {
			return true, nil
		}
		return false, nil
	})
}

// ForReadyDeployment is used to wait until Deployment is ready
func (w Wait) ForReadyDeployment(name string) error {
	return wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		deployment, err := w.KubeClient.AppsV1().Deployments(w.Namespace).Get(name, meta.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		replicas := int32(1)
		if deployment.Spec.Replicas != nil {
			replicas = *deployment.Spec.Replicas
		}

		if deployment.Status.ReadyReplicas == replicas {
			return true, nil
		}
		return false, nil
	})
}

// ForStatefulSet is used to wait until StatefulSet is created
func (w Wait) ForStatefulSet(name string) error {
	return wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		_, err = w.KubeClient.AppsV1().StatefulSets(w.Namespace).Get(name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	})
}

// ForDeployment is used to wait until Deployment is created
func (w Wait) ForDeployment(name string) error {
	return wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		_, err = w.KubeClient.AppsV1().Deployments(w.Namespace).Get(name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	})
}
