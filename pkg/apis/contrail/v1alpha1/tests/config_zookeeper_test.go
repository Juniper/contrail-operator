package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestZookeeperConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "zookeeper1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "zookeeper1-zookeeper-configmap",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client
	t.Run("default zookeeper config test", func(t *testing.T) {
		assert.NoError(t, environment.zookeeperResource.InstanceConfiguration(request,
			"zookeeper1-zookeeper-configmap", &environment.zookeeperPodList, cl),
			"Error while configuring instance")
		assert.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.zookeeperConfigMap),
			"Error while gathering zookeeper config map")

		if environment.zookeeperConfigMap.Data["zoo.cfg"] != zookeeperConfig {
			configDiff := diff.Diff(environment.zookeeperConfigMap.Data["zoo.cfg"], zookeeperConfig)
			t.Fatalf("get zookeeper config: \n%v\n", configDiff)
		}
	})

	adminEnableServer := false
	environment.zookeeperResource.Spec.ServiceConfiguration.AdminEnableServer = &adminEnableServer
	adminPort := 21833
	environment.zookeeperResource.Spec.ServiceConfiguration.AdminPort = &adminPort
	request = reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "zookeeper1",
			Namespace: "default",
		},
	}
	t.Run("custom zookeeper config test", func(t *testing.T) {
		require.NoError(t, environment.zookeeperResource.InstanceConfiguration(request,
			"zookeeper1-zookeeper-configmap", &environment.zookeeperPodList, cl),
			"Error while configuring instance")
		require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.zookeeperConfigMap),
			"Error while gathering zookeeper config map")
		zookeeperTest, err := ini.Load([]byte(environment.zookeeperConfigMap.Data["zoo.cfg"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "false", zookeeperTest.Section("DEFAULT").Key("admin.enableServer").String())
		assert.Equal(t, "21833", zookeeperTest.Section("DEFAULT").Key("admin.serverPort").String())
	})
}

var zookeeperConfig = `dataDir=/var/lib/zookeeper
tickTime=2000
initLimit=5
syncLimit=2
maxClientCnxns=60
maxSessionTimeout=120000
admin.enableServer=true
admin.serverPort=2182
standaloneEnabled=false
4lw.commands.whitelist=stat,ruok,conf,isro
reconfigEnabled=true
skipACL=yes
dynamicConfigFile=/var/lib/zookeeper/zoo.cfg.dynamic
`
