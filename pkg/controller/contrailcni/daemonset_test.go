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

var configMapMount = core.VolumeMount{
	Name:      "test-request-contrailcni-volume",
	MountPath: "/etc/contrailconfigmaps",
}

var multusContainer = core.Container{
	Name:  "multusconfig",
	Image: "busybox",
	Command: []string{
		"sh",
		"-c",
		"mkdir -p /etc/kubernetes/cni/net.d && " +
			"cp -f /etc/contrailconfigmaps/10-contrail.conf /etc/kubernetes/cni/net.d/10-contrail.conf && " +
			"mkdir -p /var/run/multus/cni/net.d && " +
			"cp -f /etc/contrailconfigmaps/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf"},
	VolumeMounts: []core.VolumeMount{
		core.VolumeMount{
			Name:      "etc-kubernetes-cni",
			MountPath: "/etc/kubernetes/cni",
		},
		core.VolumeMount{
			Name:      "multus-cni",
			MountPath: "/var/run/multus",
		},
	},
	ImagePullPolicy: "Always",
}

func TestGetDaemonsetK8s(t *testing.T) {
	testCNIDirs := contrailcni.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "k8s",
	}
	ds := contrailcni.GetDaemonset(testCNIDirs, requestName, instanceType)
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Contains(t, ds.Spec.Template.Spec.InitContainers[0].VolumeMounts, configMapMount)
	assert.NotContains(t, ds.Spec.Template.Spec.InitContainers, multusContainer)
}

func TestGetDaemonsetOpenshift(t *testing.T) {
	testCNIDirs := contrailcni.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "openshift",
	}
	ds := contrailcni.GetDaemonset(testCNIDirs, requestName, instanceType)
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Contains(t, ds.Spec.Template.Spec.InitContainers[0].VolumeMounts, configMapMount)
	assert.Contains(t, ds.Spec.Template.Spec.InitContainers, multusContainer)
}
