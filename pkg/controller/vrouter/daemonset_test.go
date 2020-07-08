package vrouter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"
)

var multusContainer = v1.Container{
	Name:  "multusconfig",
	Image: "busybox",
	Command: []string{
		"sh",
		"-c",
		"mkdir -p /etc/kubernetes/cni/net.d && " +
			"cp -f /etc/mycontrail/10-contrail.conf /etc/kubernetes/cni/net.d/10-contrail.conf && " +
			"mkdir -p /var/run/multus/cni/net.d && " +
			"cp -f /etc/mycontrail/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf"},
	VolumeMounts: []v1.VolumeMount{
		v1.VolumeMount{
			Name:      "etc-kubernetes-cni",
			MountPath: "/etc/kubernetes/cni",
		},
		v1.VolumeMount{
			Name:      "multus-cni",
			MountPath: "/var/run/multus",
		},
	},
	ImagePullPolicy: "Always",
}

var testBinariesPath = "/test/cni/bin"

var expectedCniBinVolume = v1.Volume{
	Name: "cni-bin",
	VolumeSource: v1.VolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: testBinariesPath,
		},
	},
}

func TestGetDaemonsetK8s(t *testing.T) {
	testCNIDirs := vrouter.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "k8s",
	}
	ds := vrouter.GetDaemonset(testCNIDirs)
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.NotContains(t, ds.Spec.Template.Spec.InitContainers, multusContainer)
}

func TestGetDaemonsetOpenshift(t *testing.T) {
	testCNIDirs := vrouter.CniDirs{
		BinariesDirectory: testBinariesPath,
		DeploymentType:    "openshift",
	}
	ds := vrouter.GetDaemonset(testCNIDirs)
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Contains(t, ds.Spec.Template.Spec.InitContainers, multusContainer)
}
