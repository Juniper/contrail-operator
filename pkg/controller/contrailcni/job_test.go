package contrailcni_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/controller/contrailcni"
)

var trueVal = true

var requestName = "test-request"
var instanceType = "contrailcni"

var testBinariesPath = "/test/cni/bin"

var expectedCniBinVolume = core.Volume{
	Name: "cni-bin",
	VolumeSource: core.VolumeSource{
		HostPath: &core.HostPathVolumeSource{
			Path: testBinariesPath,
		},
	},
}

var k8sCommand = []string{"sh", "-c",
	"mkdir -p /host/etc_cni/net.d && " +
		"mkdir -p /var/lib/contrail/ports/vm && " +
		"cp -f /usr/bin/contrail-k8s-cni /host/opt_cni_bin && " +
		"chmod 0755 /host/opt_cni_bin/contrail-k8s-cni && " +
		"cp -f /etc/contrailconfigmaps/10-contrail.conf /host/etc_cni/net.d/10-contrail.conf && " +
		"tar -C /host/opt_cni_bin -xzf /opt/cni-v0.3.0.tgz"}

var openshiftCommand = []string{"sh", "-c",
	"mkdir -p /host/etc_cni/net.d && " +
		"mkdir -p /var/lib/contrail/ports/vm && " +
		"cp -f /usr/bin/contrail-k8s-cni /host/opt_cni_bin && " +
		"chmod 0755 /host/opt_cni_bin/contrail-k8s-cni && " +
		"cp -f /etc/contrailconfigmaps/10-contrail.conf /host/etc_cni/net.d/10-contrail.conf && " +
		"tar -C /host/opt_cni_bin -xzf /opt/cni-v0.3.0.tgz" +
		" && mkdir -p /etc/kubernetes/cni/net.d && " +
		"cp -f /etc/contrailconfigmaps/10-contrail.conf /etc/kubernetes/cni/net.d/10-contrail.conf && " +
		"mkdir -p /var/run/multus/cni/net.d && " +
		"cp -f /etc/contrailconfigmaps/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf"}

var replicas int32 = 6

func TestGetJobK8s(t *testing.T) {
	testCNIDirs := contrailcni.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "k8s",
	}
	job := contrailcni.GetJob(testCNIDirs, requestName, instanceType, &replicas)
	assert.Contains(t, job.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Equal(t, job.Spec.Template.Spec.Containers[0].Command, k8sCommand)
}

func TestGetJobOpenshift(t *testing.T) {
	testCNIDirs := contrailcni.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "openshift",
	}
	job := contrailcni.GetJob(testCNIDirs, requestName, instanceType, &replicas)
	assert.Contains(t, job.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Equal(t, job.Spec.Template.Spec.Containers[0].Command, openshiftCommand)
}
