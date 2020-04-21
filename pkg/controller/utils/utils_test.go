package utils_test

import (
	"testing"
	tm "github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/stretchr/testify/assert"
)

func TestUntils(t *testing.T){
	t.Run("testing utils with WebuiGroupKind", func(t *testing.T){
        expected := "Webui.contrail.juniper.net"
		got := tm.WebuiGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind Vrouter.", func(t *testing.T){
        expected := "Vrouter.contrail.juniper.net"
		got := tm.VrouterGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ControlGroupKind", func(t *testing.T){
        expected := "Control.contrail.juniper.net"
		got := tm.ControlGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ConfigGroupKind", func(t *testing.T){
        expected := "Config.contrail.juniper.net"
		got := tm.ConfigGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind KubemanagerGroupKind", func(t *testing.T){
        expected := "Kubemanager.contrail.juniper.net"
		got := tm.KubemanagerGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind CassandraGroupKind", func(t *testing.T){
        expected := "Cassandra.contrail.juniper.net"
		got := tm.CassandraGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ZookeeperGroupKind", func(t *testing.T){
        expected := "Zookeeper.contrail.juniper.net"
		got := tm.ZookeeperGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind RabbitmqGroupKind", func(t *testing.T){
        expected := "Rabbitmq.contrail.juniper.net"
		got := tm.RabbitmqGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ManagerGroupKind", func(t *testing.T){
        expected := "Manager.contrail.juniper.net"
		got := tm.ManagerGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ReplicaSetGroupKind", func(t *testing.T){
        expected := "ReplicaSet.apps"
		got := tm.ReplicaSetGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind DeploymentGroupKind", func(t *testing.T){
        expected := "Deployment.apps"
		got := tm.DeploymentGroupKind()
		assert.Equal(t, expected, got.String())
	})
}
