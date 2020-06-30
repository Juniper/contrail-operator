package ha_upgrade_e2e

import (
	"context"
	e2e "github.com/Juniper/contrail-operator/test/e2e"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/test/logger"
	wait "github.com/Juniper/contrail-operator/test/wait"
)

var scmRevision = getEnv("BUILD_SCM_REVISION", "latest")
var scmBranch = getEnv("BUILD_SCM_BRANCH", "master")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var versionMap = map[string]string{
	"cassandra":              "3.11.4",
	"zookeeper":              "3.5.5",
	"cemVersion":             "2005.42",
	"python":                 "3.8.2-alpine",
	"redis":                  "4.0.2",
	"busybox":                "1.31",
	"postgres":               "12.2",
	"postgres-client":        "1.0",
	"openstack":              "train-2005",
	"rabbitmq":               "3.7",
	"contrail-statusmonitor": scmBranch + "." + scmRevision,
}

func TestHACoreContrailServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: e2e.CleanupTimeout, RetryInterval: e2e.CleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, e2e.RetryInterval, e2e.WaitTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true
		oneVal := int32(1)

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-keystone-adminpass-secret",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"password": "test123",
			},
		}

		psql := &contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "hatest-psql"},
			Spec: contrail.PostgresSpec{
				Containers: []*contrail.Container{
					{Name: "postgres", Image: "registry:5000/common-docker-third-party/contrail/postgres:" + versionMap["postgres"]},
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
				},
			},
		}

		memcached := &contrail.Memcached{
			ObjectMeta: meta.ObjectMeta{
				Namespace: namespace,
				Name:      "hatest-memcached",
			},
			Spec: contrail.MemcachedSpec{
				ServiceConfiguration: contrail.MemcachedConfiguration{
					Containers: []*contrail.Container{{Name: "memcached", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:" + versionMap["openstack"]}},
				},
			},
		}

		keystoneResource := &contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "hatest-keystone"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.CommonConfiguration{HostNetwork: &trueVal},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					MemcachedInstance:  "hatest-memcached",
					PostgresInstance:   "hatest-psql",
					KeystoneSecretName: "hatest-keystone-adminpass-secret",
					ListenPort:         5555,
					Containers: []*contrail.Container{
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
						{Name: "keystoneDbInit", Image: "registry:5000/common-docker-third-party/contrail/postgresql-client:" + versionMap["postgres-client"]},
						{Name: "keystoneInit", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:" + versionMap["openstack"]},
						{Name: "keystone", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:" + versionMap["openstack"]},
						{Name: "keystoneSsh", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone-ssh:" + versionMap["openstack"]},
						{Name: "keystoneFernet", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone-fernet:" + versionMap["openstack"]},
					},
				},
			},
		}

		cassandras := []*contrail.Cassandra{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-cassandra",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:" + versionMap["cassandra"]},
					},
				},
			},
		}}

		zookeepers := []*contrail.Zookeeper{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-zookeeper",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:" + versionMap["zookeeper"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
					},
				},
			},
		}}

		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-rabbitmq",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.RabbitmqSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:" + versionMap["rabbitmq"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
					},
				},
			},
		}

		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-config",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ConfigSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.ConfigConfiguration{
					CassandraInstance: "hatest-cassandra",
					ZookeeperInstance: "hatest-zookeeper",
					KeystoneInstance:  "hatest-keystone",
					AuthMode:          "keystone",
					Containers: []*contrail.Container{
						{Name: "api", Image: "registry:5000/contrail-nightly/contrail-controller-config-api:" + versionMap["cemVersion"]},
						{Name: "devicemanager", Image: "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + versionMap["cemVersion"]},
						{Name: "dnsmasq", Image: "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + versionMap["cemVersion"]},
						{Name: "schematransformer", Image: "registry:5000/contrail-nightly/contrail-controller-config-schema:" + versionMap["cemVersion"]},
						{Name: "servicemonitor", Image: "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + versionMap["cemVersion"]},
						{Name: "analyticsapi", Image: "registry:5000/contrail-nightly/contrail-analytics-api:" + versionMap["cemVersion"]},
						{Name: "collector", Image: "registry:5000/contrail-nightly/contrail-analytics-collector:" + versionMap["cemVersion"]},
						{Name: "queryengine", Image: "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + versionMap["cemVersion"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + versionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + versionMap["redis"]},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/busybox:" + versionMap["busybox"]},
						{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + versionMap["contrail-statusmonitor"]},
					},
				},
			},
		}

		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "hatest-webui",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:       &trueVal,
					Replicas:     &oneVal,
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					CassandraInstance: "hatest-cassandra",
					KeystoneInstance:  "hatest-keystone",
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
						{Name: "nodeinit", Image: "registry:5000/contrail-nightly/contrail-node-init:" + versionMap["cemVersion"]},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:" + versionMap["redis"]},
						{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + versionMap["cemVersion"]},
						{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + versionMap["cemVersion"]},
					},
				},
			},
		}

		cluster := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cluster1",
				Namespace: namespace,
			},
			Spec: contrail.ManagerSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Replicas:    &oneVal,
					HostNetwork: &trueVal,
				},
				KeystoneSecretName: "hatest-keystone-adminpass-secret",
				Services: contrail.Services{
					Cassandras: cassandras,
					Zookeepers: zookeepers,
					Config:     config,
					Webui:      webui,
					Keystone:   keystoneResource,
					Rabbitmq:   rabbitmq,
					Postgres:   psql,
					Memcached:  memcached,
				},
			},
		}

		t.Run("when manager resource with Config and dependencies is created", func(t *testing.T) {
			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: e2e.CleanupTimeout, RetryInterval: e2e.CleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: e2e.CleanupTimeout, RetryInterval: e2e.CleanupRetryInterval})
			assert.NoError(t, err)

			w := wait.Wait{
				Namespace:     namespace,
				Timeout:       e2e.WaitTimeout,
				RetryInterval: e2e.RetryInterval,
				KubeClient:    f.KubeClient,
				Logger:        log,
			}

			t.Run("then a ready Zookeeper StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-zookeeper-zookeeper-statefulset"))
			})

			t.Run("then a ready Cassandra StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-cassandra-cassandra-statefulset"))
			})

			t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-keystone-keystone-statefulset"))
			})

			t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset"))
			})

			t.Run("then a ready webui StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("hatest-webui-webui-statefulset"))
			})

			t.Run("and when config is scaled up from 1 to 3 node", func(t *testing.T) {
				configInstance := &contrail.Config{}
				err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "hatest-config", Namespace: namespace}, configInstance)
				assert.NoError(t, err)

				var replicas int32 = 3
				configInstance.Spec.CommonConfiguration.Replicas = &replicas
				err = f.Client.Update(context.TODO(), configInstance)
				assert.NoError(t, err)

				t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
					assert.NoError(t, w.ForReadyStatefulSet("hatest-config-config-statefulset"))
				})

			})

		})

	})
}
