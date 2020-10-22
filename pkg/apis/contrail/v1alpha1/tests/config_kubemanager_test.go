package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestKubemanagerConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "kubemanager1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "kubemanager1-kubemanager-configmap",
		Namespace: "default",
	}
	secretNamespacedName := types.NamespacedName{
		Name:      "kubemanagersecret",
		Namespace: "default",
	}
	environment := SetupEnv()
	cl := *environment.client
	clientset := kubernetes.Clientset{}
	require.NoError(t, environment.kubemanagerResource.InstanceConfiguration(
		request, &environment.kubemanbagerPodList, cl, k8s.ClusterConfig{Client: clientset.CoreV1()}),
		"Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), secretNamespacedName, &environment.kubemanagerSecret),
		"Error while gathering kubemanager secret")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.kubemanagerConfigMap),
		"Error while gathering kubemanager configmap")

	if environment.kubemanagerConfigMap.Data["kubemanager.1.1.6.1"] != kubemanagerConfig {
		diff := diff.Diff(environment.kubemanagerConfigMap.Data["kubemanager.1.1.6.1"], kubemanagerConfig)
		t.Fatalf("get kubemanager config: \n%v\n", diff)
	}
}

var kubemanagerConfig = `[DEFAULTS]
host_ip=1.1.6.1
orchestrator=kubernetes
token=THISISATOKEN
log_file=/var/log/contrail/contrail-kube-manager.log
log_level=SYS_DEBUG
log_local=1
nested_mode=0
http_server_ip=0.0.0.0
[KUBERNETES]
kubernetes_api_server=10.96.0.1
kubernetes_api_port=8080
kubernetes_api_secure_port=6443
cluster_name=kubernetes
cluster_project={}
cluster_network={}
pod_subnets=10.32.0.0/12
ip_fabric_subnets=10.64.0.0/12
service_subnets=10.96.0.0/12
ip_fabric_forwarding=false
ip_fabric_snat=true
host_network_service=false
[VNC]
public_fip_pool={}
vnc_endpoint_ip=1.1.1.1,1.1.1.2,1.1.1.3
vnc_endpoint_port=8082
rabbit_server=1.1.4.1,1.1.4.2,1.1.4.3
rabbit_port=15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.6.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.6.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=True
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.6.1.pem
sandesh_certfile=/etc/certificates/server-1.1.6.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`
