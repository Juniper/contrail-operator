package contrailtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestRabbitmqConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		types.NamespacedName{
			Name:      "rabbitmq1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "rabbitmq1-rabbitmq-configmap",
		Namespace: "default",
	}
	configMapRunnerNamespacedName := types.NamespacedName{
		Name:      "rabbitmq1-rabbitmq-configmap-runner",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client
	t.Run("rabbitmq config test", func(t *testing.T) {
		require.NoError(t, environment.rabbitmqResource.InstanceConfiguration(request, &environment.rabbitmqPodList, cl),
			"Error while configuring instance")
		require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.rabbitmqConfigMap),
			"Error while gathering rabbitmq configmap")

		for _, k := range []string{
			"rabbitmq-1.1.4.1.conf", "rabbitmq-1.1.4.2.conf", "rabbitmq-1.1.4.3.conf", "rabbitmq.nodes",
			"RABBITMQ_ERLANG_COOKIE", "RABBITMQ_USE_LONGNAME", "RABBITMQ_CONFIG_FILE", "RABBITMQ_PID_FILE",
			"RABBITMQ_PID_FILE", "RABBITMQ_CONF_ENV_FILE", "plugins.conf",
			// "definitions.json" TODO: Handle random password_hash in test
		} {
			if configDiff := diff.Diff(environment.rabbitmqConfigMap.Data[k], rabbitmqConfig[k]); configDiff != "" {
				t.Fatalf("get rabbitmq config key = %v: \n%v\n", k, configDiff)
			}
		}
	})
	t.Run("rabbitmq configmap runner test", func(t *testing.T) {
		require.NoError(t, cl.Get(context.TODO(), configMapRunnerNamespacedName, &environment.rabbitmqConfigMap2),
			"Error while gathering configmap")

		if environment.rabbitmqConfigMap2.Data["run.sh"] != rabbitmqConfigRunner {
			configDiff := diff.Diff(environment.rabbitmqConfigMap2.Data["run.sh"], rabbitmqConfigRunner)
			t.Fatalf("get rabbitmq config: \n%v\n", configDiff)
		}
	})
}

var rabbitmqConfigRunner = `#!/bin/bash
echo $RABBITMQ_ERLANG_COOKIE > /var/lib/rabbitmq/.erlang.cookie
chmod 0600 /var/lib/rabbitmq/.erlang.cookie
export RABBITMQ_NODENAME=rabbit@${POD_IP}
if [[ $(grep $POD_IP /etc/rabbitmq/0) ]] ; then
  rabbitmq-server
else
  rabbitmqctl --node rabbit@${POD_IP} forget_cluster_node rabbit@${POD_IP}
  rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) ping
  while [[ $? -ne 0 ]]; do
	rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) ping
  done
  rabbitmq-server -detached
  rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) node_health_check
  while [[ $? -ne 0 ]]; do
	rabbitmqctl --node rabbit@$(cat /etc/rabbitmq/0) node_health_check
  done
  rabbitmqctl stop_app
  sleep 2
  rabbitmqctl join_cluster rabbit@$(cat /etc/rabbitmq/0)
  rabbitmqctl shutdown
  rabbitmq-server
fi
`

var rabbitmqConfig = map[string]string{
	"rabbitmq-1.1.4.1.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.1.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.1.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
cluster_partition_handling = autoheal
cluster_formation.peer_discovery_backend = classic_config
cluster_formation.classic_config.nodes.1 = rabbit@1.1.4.1
cluster_formation.classic_config.nodes.2 = rabbit@1.1.4.2
cluster_formation.classic_config.nodes.3 = rabbit@1.1.4.3
`,
	"rabbitmq-1.1.4.2.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.2.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.2.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
cluster_partition_handling = autoheal
cluster_formation.peer_discovery_backend = classic_config
cluster_formation.classic_config.nodes.1 = rabbit@1.1.4.1
cluster_formation.classic_config.nodes.2 = rabbit@1.1.4.2
cluster_formation.classic_config.nodes.3 = rabbit@1.1.4.3
`,
	"rabbitmq-1.1.4.3.conf": `listeners.tcp.default = 5673
listeners.ssl.default = 15673
loopback_users = none
management.tcp.port = 15671
management.load_definitions = /etc/rabbitmq/definitions.json
ssl_options.cacertfile = /etc/ssl/certs/kubernetes/ca-bundle.crt
ssl_options.keyfile = /etc/certificates/server-key-1.1.4.3.pem
ssl_options.certfile = /etc/certificates/server-1.1.4.3.crt
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = true
cluster_partition_handling = autoheal
cluster_formation.peer_discovery_backend = classic_config
cluster_formation.classic_config.nodes.1 = rabbit@1.1.4.1
cluster_formation.classic_config.nodes.2 = rabbit@1.1.4.2
cluster_formation.classic_config.nodes.3 = rabbit@1.1.4.3
`,
	"rabbitmq.nodes":         fmt.Sprintf("1.1.4.1\n1.1.4.2\n1.1.4.3\n"),
	"0":                      "1.1.4.1",
	"1":                      "1.1.4.2",
	"2":                      "1.1.4.3",
	"RABBITMQ_ERLANG_COOKIE": "47EFF3BB-4786-46E0-A5BB-58455B3C2CB4",
	"RABBITMQ_USE_LONGNAME":  "true",
	"RABBITMQ_CONFIG_FILE":   "/etc/rabbitmq/rabbitmq-${POD_IP}.conf",
	"RABBITMQ_PID_FILE":      "/var/run/rabbitmq.pid",
	"RABBITMQ_CONF_ENV_FILE": "/var/lib/rabbitmq/rabbitmq.env",
	"definitions.json":       rabbitmqDefinition,
	"plugins.conf":           "[rabbitmq_management,rabbitmq_management_agent,rabbitmq_peer_discovery_k8s].",
}

var rabbitmqDefinition = `{
  "users": [
    {
      "name": "user",
      "password_hash": "0TeE2b17AedPdBgaZSwAtAfgwXDSqqoDki44i3AcXqxU1DIB",
      "tags": "administrator"
    }
  ],
  "vhosts": [
    {
      "name": "vhost"
    }
  ],
  "permissions": [
    {
      "user": "user",
      "vhost": "vhost",
      "configure": ".*",
      "write": ".*",
      "read": ".*"
    }
  ],
}
`
