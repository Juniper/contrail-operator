package vrouter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"
)

func TestGetDaemonset(t *testing.T) {
	testBinariesPath := "/test/cni/bin"
	testConfigPath := "/config/test/cni"
	testCNIDirs := vrouter.CniDirs{
		BinariesDirectory:    testBinariesPath,
		ConfigFilesDirectory: testConfigPath,
	}
	ds := vrouter.GetDaemonset(testCNIDirs)
	expectedCniBinVolume := v1.Volume{
		Name: "cni-bin",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: testBinariesPath,
			},
		},
	}
	expectedCniConfigVolume := v1.Volume{
		Name: "cni-config-files",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: testConfigPath,
			},
		},
	}
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniBinVolume)
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, expectedCniConfigVolume)
}
