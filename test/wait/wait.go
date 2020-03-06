package wait

import (
	"github.com/Juniper/contrail-operator/test/logger"
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
	Logger        logger.Logger
}

// ForReadyStatefulSet is used to wait until StatefulSet is ready
func (w Wait) ForReadyStatefulSet(name string) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
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
	w.dumpPodsOnError(err)
	return err
}

// ForReadyDeployment is used to wait until Deployment is ready
func (w Wait) ForReadyDeployment(name string) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
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
	w.dumpPodsOnError(err)
	return err
}

// ForStatefulSet is used to wait until StatefulSet is created
func (w Wait) ForStatefulSet(name string) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		_, err = w.KubeClient.AppsV1().StatefulSets(w.Namespace).Get(name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	})
	w.dumpPodsOnError(err)
	return err
}

// ForDeployment is used to wait until Deployment is created
func (w Wait) ForDeployment(name string) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		_, err = w.KubeClient.AppsV1().Deployments(w.Namespace).Get(name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	})
	w.dumpPodsOnError(err)
	return err
}

func (w Wait) dumpPodsOnError(err error) {
	if err != nil {
		w.Logger.DumpPods()
	}
}
