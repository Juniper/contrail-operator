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
	environment.configResource.Spec.ServiceConfiguration.AuthMode = "keyrock"

	require.NoError(t, environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl), "Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap))

	t.Run("custom queryengine config settings", func(t *testing.T) {
		queryEngineIni, err := ini.Load([]byte(environment.configConfigMap.Data["queryengine.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", queryEngineIni.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
	})

	t.Run("custom collector config settings", func(t *testing.T) {
		collectorIni, err := ini.Load([]byte(environment.configConfigMap.Data["collector.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", collectorIni.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
		assert.Equal(t, "111", collectorIni.Section("DEFAULT").Key("analytics_config_audit_ttl").String(), "Invalid analytics_config_audit_ttl")
		assert.Equal(t, "444", collectorIni.Section("DEFAULT").Key("analytics_statistics_ttl").String(), "Invalid analytics_statistics_ttl")
		assert.Equal(t, "333", collectorIni.Section("DEFAULT").Key("analytics_flow_ttl").String(), "Invalid analytics_flow_ttl")
	})
}
func TestCustomizedAuthConfig(t *testing.T) {
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
	environment.configResource.Spec.ServiceConfiguration.AuthMode = "keyrock"

	require.NoError(t, environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl), "Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap))

	t.Run("custom auth mode", func(t *testing.T) {

		apiIni, err := ini.Load([]byte(environment.configConfigMap.Data["api.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "keyrock", apiIni.Section("DEFAULTS").Key("auth").String())
		assert.Equal(t, "no-auth", apiIni.Section("DEFAULTS").Key("aaa_mode").String())

		vncApiIni, err := ini.Load([]byte(environment.configConfigMap.Data["vnc.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "keyrock", vncApiIni.Section("auth").Key("AUTHN_TYPE").String())
		assert.Equal(t, "", vncApiIni.Section("auth").Key("AUTHN_PROTOCOL").String())
		assert.Equal(t, "", vncApiIni.Section("auth").Key("AUTHN_SERVER").String())
		assert.Equal(t, "0", vncApiIni.Section("auth").Key("AUTHN_PORT").String())
		assert.Equal(t, "", vncApiIni.Section("auth").Key("AUTHN_DOMAIN").String())

		contrailKeystoneAuthIni, err := ini.Load([]byte(environment.configConfigMap.Data["contrail-keystone-auth.conf"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "", contrailKeystoneAuthIni.Section("KEYSTONE").Key("auth_host").String())
		assert.Equal(t, "0", contrailKeystoneAuthIni.Section("KEYSTONE").Key("auth_port").String())
		assert.Equal(t, "://:0/v3", contrailKeystoneAuthIni.Section("KEYSTONE").Key("auth_url").String())
		assert.Equal(t, "", contrailKeystoneAuthIni.Section("KEYSTONE").Key("user_domain_name").String())
		assert.Equal(t, "", contrailKeystoneAuthIni.Section("KEYSTONE").Key("project_domain_name").String())
		assert.Equal(t, "", contrailKeystoneAuthIni.Section("KEYSTONE").Key("region_name").String())

		serviceMonitorIni, err := ini.Load([]byte(environment.configConfigMap.Data["servicemonitor.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "no-auth", serviceMonitorIni.Section("SCHEDULER").Key("aaa_mode").String())

		analyticsApiIni, err := ini.Load([]byte(environment.configConfigMap.Data["analyticsapi.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "no-auth", analyticsApiIni.Section("DEFAULTS").Key("aaa_mode").String())
	})
}

func TestCustomizedPortConfig(t *testing.T) {
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
	environment.configResource.Spec.ServiceConfiguration.APIPort = makeIntPointer(1111)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsPort = makeIntPointer(2222)
	environment.configResource.Spec.ServiceConfiguration.CollectorPort = makeIntPointer(3333)
	environment.configResource.Spec.ServiceConfiguration.RedisPort = makeIntPointer(4444)
	environment.configResource.Spec.ServiceConfiguration.ApiIntrospectPort = makeIntPointer(5555)
	environment.configResource.Spec.ServiceConfiguration.SchemaIntrospectPort = makeIntPointer(6666)
	environment.configResource.Spec.ServiceConfiguration.DeviceManagerIntrospectPort = makeIntPointer(7777)
	environment.configResource.Spec.ServiceConfiguration.SvcMonitorIntrospectPort = makeIntPointer(8888)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsApiIntrospectPort = makeIntPointer(9999)
	environment.configResource.Spec.ServiceConfiguration.CollectorIntrospectPort = makeIntPointer(10101)

	require.NoError(t, environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl), "Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap))

	t.Run("custom api config", func(t *testing.T) {
		apiIni, err := ini.Load([]byte(environment.configConfigMap.Data["api.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "1111", apiIni.Section("DEFAULTS").Key("listen_port").String())
		assert.Equal(t, "5555", apiIni.Section("DEFAULTS").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", apiIni.Section("DEFAULTS").Key("collectors").String())
	})

	t.Run("custom vncapi config", func(t *testing.T) {
		vncApiIni, err := ini.Load([]byte(environment.configConfigMap.Data["vnc.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "1111", vncApiIni.Section("global").Key("WEB_PORT").String())
	})

	t.Run("custom devicemanager config", func(t *testing.T) {
		devicemanagerIni, err := ini.Load([]byte(environment.configConfigMap.Data["devicemanager.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "7777", devicemanagerIni.Section("DEFAULTS").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", devicemanagerIni.Section("DEFAULTS").Key("collectors").String())
	})

	t.Run("custom schematransformer config", func(t *testing.T) {
		schematransformerIni, err := ini.Load([]byte(environment.configConfigMap.Data["schematransformer.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "6666", schematransformerIni.Section("DEFAULTS").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", schematransformerIni.Section("DEFAULTS").Key("collectors").String())
	})

	t.Run("custom servicemonitor config", func(t *testing.T) {
		serviceMonitorIni, err := ini.Load([]byte(environment.configConfigMap.Data["servicemonitor.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "8888", serviceMonitorIni.Section("DEFAULTS").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", serviceMonitorIni.Section("DEFAULTS").Key("collectors").String())
		assert.Equal(t, "1.1.1.1:2222 1.1.1.2:2222 1.1.1.3:2222", serviceMonitorIni.Section("SCHEDULER").Key("analytics_server_list").String())
	})

	t.Run("custom servicemonitor config", func(t *testing.T) {
		analyticsApiIni, err := ini.Load([]byte(environment.configConfigMap.Data["analyticsapi.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "9999", analyticsApiIni.Section("DEFAULTS").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", analyticsApiIni.Section("DEFAULTS").Key("collectors").String())
		assert.Equal(t, "1.1.1.1:1111 1.1.1.2:1111 1.1.1.3:1111", analyticsApiIni.Section("DEFAULTS").Key("api_server").String())
		assert.Equal(t, "1.1.1.1:4444 1.1.1.2:4444 1.1.1.3:4444", analyticsApiIni.Section("REDIS").Key("redis_uve_list").String())
	})

	t.Run("custom queryengine config", func(t *testing.T) {
		queryEngineIni, err := ini.Load([]byte(environment.configConfigMap.Data["queryengine.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "1.1.1.1:3333 1.1.1.2:3333 1.1.1.3:3333", queryEngineIni.Section("DEFAULT").Key("collectors").String())
		assert.Equal(t, "1.1.1.1:4444 1.1.1.2:4444 1.1.1.3:4444", queryEngineIni.Section("REDIS").Key("server_list").String())
	})

	t.Run("custom collector config", func(t *testing.T) {
		collectorIni, err := ini.Load([]byte(environment.configConfigMap.Data["collector.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "10101", collectorIni.Section("DEFAULT").Key("http_server_port").String())
		assert.Equal(t, "1.1.1.1:1111 1.1.1.2:1111 1.1.1.3:1111", collectorIni.Section("API_SERVER").Key("api_server_list").String())
	})
}
