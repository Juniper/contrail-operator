package wait

import (
	"context"
	"strings"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/test/logger"
)

// Contrail is used to wait until certain contrail resources reach some condition
type Contrail struct {
	Namespace     string
	RetryInterval time.Duration
	Timeout       time.Duration
	Client        test.FrameworkClient
	Logger        logger.Logger
}

// ForManagerCondition is used to wait until manager has expected condition met
func (c Contrail) ForManagerCondition(name string, expected contrail.ManagerConditionType) error {
	m := &contrail.Manager{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, m)
		if apierrors.IsNotFound(err) {
			return false, nil
		}

		for _, condition := range m.Status.Conditions {
			if condition.Type != expected {
				continue
			}

			if condition.Status == contrail.ConditionTrue {
				return true, nil
			}

		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

// ForSwiftActive is used to wait until Swift is active
func (c Contrail) ForSwiftActive(name string) error {
	s := &contrail.Swift{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, s)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if s.Status.Active {
			return true, nil
		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

// ForWebuiActive is used to wait until Webui is active
func (c Contrail) ForWebUIActive(name string) error {
	w := &contrail.Webui{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, w)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if w.Status.Active {
			return true, nil
		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

// ForPostgresActive is used to wait until Postgres is active
func (c Contrail) ForPostgresActive(name string) error {
	s := &contrail.Postgres{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, s)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if s.Status.Active {
			return true, nil
		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

// ForPodImageChange is used to wait until Image is updated
func (c Contrail) ForPodImageChange(kubeClient kubernetes.Interface, labelSelector string, newImage string, container string) error {
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		podList, err := kubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return false, nil
		}
		imageUpdatedCount := 0
		for _, pod := range podList.Items {
			for _, c := range pod.Spec.Containers {
				if c.Name == container && c.Image == newImage {
					imageUpdatedCount++
				}
			}
		}
		if imageUpdatedCount == len(podList.Items) {
			return true, nil
		}
		return false, nil
	})
	c.dumpPodsOnError(err)
	return err
}

// ForPodUidChange is used to wait until pod has a new Uid
func (c Contrail) ForPodUidChange(kubeClient kubernetes.Interface, podName string, oldUid types.UID) error {
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, getErr error) {
		pod, getErr := kubeClient.CoreV1().Pods("contrail").Get(podName, meta.GetOptions{})
		if getErr != nil {
			return false, nil
		}
		newUid := pod.UID

		return !strings.EqualFold(string(oldUid), string(newUid)), nil
	})
	c.dumpPodsOnError(err)
	return err
}

// ForManagerDeletion is used to wait until manager is deleted
func (c Contrail) ForManagerDeletion(name string) error {
	m := &contrail.Manager{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, m)
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

// ForCommandUpgradeStateChange is used to wait until Command Upgrade State is updated
func (c Contrail) ForCommandUpgradeState(name string, newState contrail.CommandUpgradeState) error {
	command := &contrail.Command{}
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, command)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if command.Status.UpgradeState == newState {
			return true, nil
		}
		return false, err
	})
	c.dumpPodsOnError(err)
	return err
}

func (c Contrail) dumpPodsOnError(err error) {
	if err != nil {
		c.Logger.DumpPods()
	}
}
