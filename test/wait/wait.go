package wait

import (
	"context"
	"time"

	"github.com/Juniper/contrail-operator/test/logger"

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
func (w Wait) ForReadyStatefulSet(name string, expectedReplicas int32) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		statefulSet, err := w.KubeClient.AppsV1().StatefulSets(w.Namespace).Get(context.Background(), name, meta.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			w.Logger.Logf("request to kube api returned error: %v", err)
			return false, nil
		}
		if statefulSet.Status.ReadyReplicas == expectedReplicas {
			return true, nil
		}
		return false, nil
	})
	w.dumpPodsOnError(err)
	return err
}

// ForReadyDeployment is used to wait until Deployment is ready
func (w Wait) ForReadyDeployment(name string, expectedReplicas int32) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		deployment, err := w.KubeClient.AppsV1().Deployments(w.Namespace).Get(context.Background(), name, meta.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			w.Logger.Logf("request to kube api returned error: %v", err)
			return false, nil
		}

		if deployment.Status.ReadyReplicas == expectedReplicas {
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
		_, err = w.KubeClient.AppsV1().StatefulSets(w.Namespace).Get(context.Background(), name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		w.Logger.Logf("request to kube api returned error: %v", err)
		return false, nil
	})
	w.dumpPodsOnError(err)
	return err
}

// ForDeployment is used to wait until Deployment is created
func (w Wait) ForDeployment(name string) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		_, err = w.KubeClient.AppsV1().Deployments(w.Namespace).Get(context.Background(), name, meta.GetOptions{})
		if err == nil {
			return true, nil
		}
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		w.Logger.Logf("request to kube api returned error: %v", err)
		return false, nil
	})
	w.dumpPodsOnError(err)
	return err
}

// Poll is a wrapper for retrying an function until it ends with success
func (w Wait) Poll(repeatable func() (done bool, err error)) error {
	return wait.Poll(w.RetryInterval, w.Timeout, repeatable)
}

// RetryRequest is used to retry requests until it doesn't return error
func (w Wait) RetryRequest(f func() error) error {
	err := wait.Poll(w.RetryInterval, w.Timeout, func() (done bool, err error) {
		if err = f(); err != nil {
			w.Logger.Logf("request failed: %v", err)
			return false, nil
		}
		return true, nil
	})

	return err
}

func (w Wait) dumpPodsOnError(err error) {
	if err != nil {
		w.Logger.DumpPods()
	}
}
