package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestCassandraConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "cassandra1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "cassandra1-cassandra-configmap",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	require.NoError(t, environment.cassandraResource.InstanceConfiguration(request, &environment.cassandraPodList, cl),
		"Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.cassandraConfigMap),
		"Error while gathering configmap")

	if environment.cassandraConfigMap.Data["1.1.2.1.yaml"] != cassandraConfig {
		configDiff := diff.Diff(environment.cassandraConfigMap.Data["1.1.2.1.yaml"], cassandraConfig)
		t.Fatalf("get cassandra config: \n%v\n", configDiff)
	}
}

var cassandraConfig = `cluster_name: ContrailConfigDB
num_tokens: 32
hinted_handoff_enabled: true
max_hint_window_in_ms: 10800000 # 3 hours
hinted_handoff_throttle_in_kb: 1024
max_hints_delivery_threads: 2
hints_directory: /var/lib/cassandra/hints
hints_flush_period_in_ms: 10000
max_hints_file_size_in_mb: 128
batchlog_replay_throttle_in_kb: 1024
authenticator: AllowAllAuthenticator
authorizer: AllowAllAuthorizer
role_manager: CassandraRoleManager
roles_validity_in_ms: 2000
permissions_validity_in_ms: 2000
credentials_validity_in_ms: 2000
partitioner: org.apache.cassandra.dht.Murmur3Partitioner
data_file_directories:
- /var/lib/cassandra/data
commitlog_directory: /var/lib/cassandra/commitlog
disk_failure_policy: stop
commit_failure_policy: stop
key_cache_size_in_mb:
key_cache_save_period: 14400
row_cache_size_in_mb: 0
row_cache_save_period: 0
counter_cache_size_in_mb:
counter_cache_save_period: 7200
saved_caches_directory: /var/lib/cassandra/saved_caches
commitlog_sync: periodic
commitlog_sync_period_in_ms: 10000
commitlog_segment_size_in_mb: 32
seed_provider:
- class_name: org.apache.cassandra.locator.SimpleSeedProvider
  parameters:
  - seeds: 1.1.2.1,1.1.2.2
concurrent_reads: 32
concurrent_writes: 32
concurrent_counter_writes: 32
concurrent_materialized_view_writes: 32
disk_optimization_strategy: ssd
memtable_allocation_type: heap_buffers
index_summary_capacity_in_mb:
index_summary_resize_interval_in_minutes: 60
trickle_fsync: false
trickle_fsync_interval_in_kb: 10240
storage_port: 7000
ssl_storage_port: 7001
listen_address: 1.1.2.1
broadcast_address: 1.1.2.1
start_native_transport: true
native_transport_port: 9042
start_rpc: true
rpc_address: 1.1.2.1
rpc_port: 9160
broadcast_rpc_address: 1.1.2.1
rpc_keepalive: true
rpc_server_type: sync
thrift_framed_transport_size_in_mb: 15
incremental_backups: false
snapshot_before_compaction: false
auto_snapshot: true
tombstone_warn_threshold: 1000
tombstone_failure_threshold: 100000
column_index_size_in_kb: 64
batch_size_warn_threshold_in_kb: 5
batch_size_fail_threshold_in_kb: 50
compaction_throughput_mb_per_sec: 16
compaction_large_partition_warning_threshold_mb: 100
sstable_preemptive_open_interval_in_mb: 50
read_request_timeout_in_ms: 5000
range_request_timeout_in_ms: 10000
write_request_timeout_in_ms: 2000
counter_write_request_timeout_in_ms: 5000
cas_contention_timeout_in_ms: 1000
truncate_request_timeout_in_ms: 60000
request_timeout_in_ms: 10000
cross_node_timeout: false
endpoint_snitch: SimpleSnitch
dynamic_snitch_update_interval_in_ms: 100
dynamic_snitch_reset_interval_in_ms: 600000
dynamic_snitch_badness_threshold: 0.1
request_scheduler: org.apache.cassandra.scheduler.NoScheduler
server_encryption_options:
  internode_encryption: all
  keystore: /etc/keystore/server-keystore.jks
  keystore_password: keystorePassword
  truststore: /etc/keystore/server-truststore.jks
  truststore_password: truststorePassword
  require_client_auth: true
  store_type: JKS
client_encryption_options:
  enabled: true
  optional: false
  keystore: /etc/keystore/server-keystore.jks
  keystore_password: keystorePassword
internode_compression: all
inter_dc_tcp_nodelay: false
tracetype_query_ttl: 86400
tracetype_repair_ttl: 604800
gc_warn_threshold_in_ms: 1000
enable_user_defined_functions: false
enable_scripted_user_defined_functions: false
windows_timer_interval: 1
transparent_data_encryption_options:
  enabled: false
  chunk_length_kb: 64
  cipher: AES/CBC/PKCS5Padding
  key_alias: testing:1
  key_provider:
  - class_name: org.apache.cassandra.security.JKSKeyProvider
    parameters:
    - keystore: conf/.keystore
      keystore_password: cassandra
      store_type: JCEKS
      key_password: cassandra
auto_bootstrap: true
`
