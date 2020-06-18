package wait

import (
	"context"
	"github.com/Juniper/contrail-operator/test/logger"
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

// ForPostgresActive is used to wait until Postgres is inactive
func (c Contrail) ForPostgresInactive(name string) error {
	err := wait.Poll(c.RetryInterval, c.Timeout, func() (done bool, err error) {
		s := &contrail.Postgres{}
		err = c.Client.Get(context.Background(), types.NamespacedName{
			Namespace: c.Namespace,
			Name:      name,
		}, s)
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		if s.Status.Active {
			return false, nil
		}
		return true, err
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

func (c Contrail) dumpPodsOnError(err error) {
	if err != nil {
		c.Logger.DumpPods()
	}
}
