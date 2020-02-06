package wait

import (
	"context"
	"time"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Contrail is used to wait until certain contrail resources reach some condition
type Contrail struct {
	Namespace     string
	RetryInterval time.Duration
	Timeout       time.Duration
	Client        test.FrameworkClient
}

// ForManagerCondition is used to wait until manager has expected condition met
func (c Contrail) ForManagerCondition(name string, expected contrail.ManagerConditionType) error {
	m := &contrail.Manager{}
	return wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
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
}

// ForManagerDeletion is used to wait until manager is deleted
func (c Contrail) ForManagerDeletion(name string) error {
	m := &contrail.Manager{}
	return wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, m)
		if apierrors.IsNotFound(err) {
			return true, nil
		}

		return false, err
	})
}
