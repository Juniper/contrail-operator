package contrailtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func makeIntPointer(v int) *int {
	return &v
}

func TestCustomizedConfigConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "config1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "config1-config-configmap",
		Namespace: "default",
	}
	environment := SetupEnv()
	cl := *environment.client
	environment.configResource.Spec.ServiceConfiguration.AnalyticsConfigAuditTTL = makeIntPointer(111)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsDataTTL = makeIntPointer(222)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsFlowTTL = makeIntPointer(333)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsStatisticsTTL = makeIntPointer(444)

	require.NoError(t, environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl), "Error while configuring instance")

	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap))

	t.Run("custom queryengine config settings", func(t *testing.T) {
		queryEngine, err := ini.Load([]byte(environment.configConfigMap.Data["queryengine.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", queryEngine.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
	})

	t.Run("custom collector config settings", func(t *testing.T) {
		collector, err := ini.Load([]byte(environment.configConfigMap.Data["collector.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", collector.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
		assert.Equal(t, "111", collector.Section("DEFAULT").Key("analytics_config_audit_ttl").String(), "Invalid analytics_config_audit_ttl")
		assert.Equal(t, "444", collector.Section("DEFAULT").Key("analytics_statistics_ttl").String(), "Invalid analytics_statistics_ttl")
		assert.Equal(t, "333", collector.Section("DEFAULT").Key("analytics_flow_ttl").String(), "Invalid analytics_flow_ttl")
	})
}
